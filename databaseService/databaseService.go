package databaseService

import (
	"encoding/gob"
	"fmt"
	"os"

	"BD_Lab1/entities"
)

/*
Define a databaseService class that takes in string fileName. Then it creates if not exist, or opens if exist files
fileName.fl & fileName.ind. Then it creates a new instance of the databaseService class and returns it.
*/
type DatabaseService struct {
	fileName         string
	mainFile         *os.File
	indexFile        *os.File
	indexTableStruct *entities.IndexTable
}

func NewDatabaseService(fileName string) (*DatabaseService, error) {
	// Create the main file || open it if it exists.
	mainFile, err := os.OpenFile(fileName+".fl", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	// Create the index file || open it if it exists.
	indexFile, err := os.OpenFile(fileName+".ind", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	// check if index file is not empty
	indexFileStat, err := indexFile.Stat()
	if err != nil {
		return nil, err
	}

	indexTable := entities.IndexTable{
		PreviousMaxUID: 1,
		OverallOffset:  0,
		Rows:           []entities.IndexTableRow{},
	}
	// if it is empty, initialize it with empty IndexTable using GOB:
	if indexFileStat.Size() == 0 {
		// initialize index table

		// write it to the index file
		err = gob.NewEncoder(indexFile).Encode(indexTable)
		if err != nil {
			return nil, err
		}
	}
	return &DatabaseService{
		fileName:         fileName,
		mainFile:         mainFile,
		indexFile:        indexFile,
		indexTableStruct: &indexTable,
	}, nil
}

func (db *DatabaseService) PrintIndexTable() {
	fmt.Println("Max user ID:", db.indexTableStruct.PreviousMaxUID)
	fmt.Println("Overall offset:", db.indexTableStruct.OverallOffset)
	fmt.Println("Rows:")
	fmt.Println("UID\tFileStartOffset\tFileEndOffset")
	for _, row := range db.indexTableStruct.Rows {
		fmt.Println(row.UID, "\t", row.FileStartOffset, "\t\t", row.FileEndOffset)
	}
}

//TODO: refactor this function to avoid returning error
func (db *DatabaseService) UpdateIndexTable(row entities.IndexTableRow, bytesWritten int) error {
	db.indexTableStruct.Rows = append(db.indexTableStruct.Rows, row)
	db.indexTableStruct.PreviousMaxUID += 1
	db.indexTableStruct.OverallOffset += bytesWritten
	return nil
}

func (db *DatabaseService) Close() error {
	err := db.mainFile.Close()
	if err != nil {
		return err
	}

	err = db.indexFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DatabaseService) SaveIndexTable() error {
	// write it to the index file
	err := gob.NewEncoder(db.indexFile).Encode(db.indexTableStruct)
	if err != nil {
		return err
	}
	return nil
}
func (db *DatabaseService) WriteToMainFile(data []byte) error {
	// Write to the main file.
	size, err := db.mainFile.Write(data)
	if err != nil {
		return err
	}

	//update the index table:
	err = db.UpdateIndexTable(entities.IndexTableRow{
		UID:             db.indexTableStruct.PreviousMaxUID + 1,
		FileStartOffset: db.indexTableStruct.OverallOffset,
		FileEndOffset:   db.indexTableStruct.OverallOffset + len(data),
	}, size)
	if err != nil {
		return err
	}
	return nil
}

// it will take in the ID & return the bytes array:
func (db *DatabaseService) GetEntryById(id uint64) ([]byte, error) {
	for _, row := range db.indexTableStruct.Rows {
		if row.UID == id {
			fmt.Println(row.FileStartOffset)
			fmt.Println(row.FileEndOffset)
			// read the bytes from the main file:
			bytes := make([]byte, row.FileEndOffset-row.FileStartOffset)
			_, err := db.mainFile.ReadAt(bytes, int64(row.FileStartOffset))
			if err != nil {
				return nil, err
			}
			return bytes, nil
		}
	}
	return nil, fmt.Errorf("no such id")
}

func (db *DatabaseService) GetAllEntriesIds() ([]uint64, error) {
	var ids []uint64
	for _, row := range db.indexTableStruct.Rows {
		ids = append(ids, row.UID)
	}
	return ids, nil
}
func (db *DatabaseService) CleanUp() error {
	// clean the main file:
	err := db.mainFile.Truncate(0)
	if err != nil {
		return err
	}
	// clean the index file:
	db.indexTableStruct.PreviousMaxUID = 0
	db.indexTableStruct.OverallOffset = 0
	db.indexTableStruct.Rows = []entities.IndexTableRow{}
	err = db.SaveIndexTable()
	if err != nil {
		fmt.Println("error saving index table: ", err)
		return err
	}
	return nil
}

func (db *DatabaseService) GetUID() int {
	return int(db.indexTableStruct.PreviousMaxUID)
}

func (db *DatabaseService) PrintIndexTableInfo() {
	fmt.Println("Printing index table info...")
	fmt.Println("PreviousMaxUID: ", db.indexTableStruct.PreviousMaxUID)
	fmt.Println("OverallOffset: ", db.indexTableStruct.OverallOffset)
	fmt.Println("Rows: ", db.indexTableStruct.Rows)
}
func (db *DatabaseService) FindIndexTableRowById(id uint64) (entities.IndexTableRow, error) {
	for _, row := range db.indexTableStruct.Rows {
		if row.UID == id {
			return row, nil
		}
	}
	return entities.IndexTableRow{}, fmt.Errorf("no such id")
}

func (db *DatabaseService) DeleteEntryById(id uint64) error {
	// 1. find the row with the given id in the index table
	_, err := db.GetEntryById(id)
	if err != nil {
		return err
	}

	// find the start & end offsets of the row with the given id
	row, err := db.FindIndexTableRowById(id)
	if err != nil {
		return err
	}
	startOffset := row.FileStartOffset
	endOffset := row.FileEndOffset
	entrySize := endOffset - startOffset

	// 2. remove the row from the main file (we have the start & end offsets)
	// 2.1 read the bytes after the deleted row
	bytes := make([]byte, db.indexTableStruct.OverallOffset-entrySize)
	_, err = db.mainFile.ReadAt(bytes, int64(endOffset))
	if err != nil {
		return err
	}

	// 2.2 overwrite the deleted row with the bytes read in the previous step
	_, err = db.mainFile.WriteAt(bytes, int64(startOffset))
	if err != nil {
		return err
	}

	// 2.3 truncate the main file
	err = db.mainFile.Truncate(int64(db.indexTableStruct.OverallOffset - entrySize))
	if err != nil {
		return err
	}

	// 2.4 update the index table
	db.indexTableStruct.OverallOffset -= entrySize

	// 3. unshift the rows after the deleted row
	for i := range db.indexTableStruct.Rows {
		if db.indexTableStruct.Rows[i].UID > id {
			db.indexTableStruct.Rows[i].FileStartOffset -= entrySize
			db.indexTableStruct.Rows[i].FileEndOffset -= entrySize
		}
	}

	// 4. remove the row from the index table
	for i := range db.indexTableStruct.Rows {
		if db.indexTableStruct.Rows[i].UID == id {
			db.indexTableStruct.Rows = append(db.indexTableStruct.Rows[:i], db.indexTableStruct.Rows[i+1:]...)
			break
		}
	}

	return nil
}
