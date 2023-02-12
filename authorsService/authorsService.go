package authorsService

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"BD_Lab1/databaseService"
	"BD_Lab1/entities"
)

type AuthorsService struct {
	dbService *databaseService.DatabaseService
}

func NewAuthorsService(name string) (*AuthorsService, error) {
	dbService, err := databaseService.NewDatabaseService(name)
	if err != nil {
		return nil, err
	}
	return &AuthorsService{
		dbService: dbService,
	}, nil
}

func (a *AuthorsService) Close() error {
	err := a.dbService.Close()
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthorsService) CreateAuthor(author *entities.Author) error {
	// check if author already exists
	UID := a.dbService.GetUID()
	author.ID = UID

	// convert author to bytes via GOB
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(author)
	if err != nil {
		return err
	}
	err = a.dbService.WriteToMainFile(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}
func (a *AuthorsService) GetAuthorById(id uint64) (*entities.Author, error) {
	// get entry from main file
	entry, err := a.dbService.GetEntryById(id)
	if err != nil {
		return nil, err
	}

	// decode entry to author
	var author entities.Author
	dec := gob.NewDecoder(bytes.NewReader(entry))
	err = dec.Decode(&author)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (a *AuthorsService) GetAllAuthors() ([]entities.Author, error) {
	IDs, err := a.dbService.GetAllEntriesIds()
	if err != nil {
		fmt.Println("error while getting all entries ids")
		return nil, err
	}
	fmt.Println(IDs)

	// iterate over all IDs and get the authors
	var authors []entities.Author
	for _, id := range IDs {
		author, err := a.GetAuthorById(id)
		if err != nil {
			fmt.Println("error while getting author by id")
			return nil, err
		}
		authors = append(authors, *author)
	}
	return authors, nil
}

func (a *AuthorsService) PrintAllInfo() error {
	authors, err := a.GetAllAuthors()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("\nAuthors:")
	fmt.Println(authors)

	fmt.Println("\nIndex Table:")
	a.dbService.PrintIndexTable()
	return nil
}

func (a *AuthorsService) PrintAllAuthors() {
	authors, err := a.GetAllAuthors()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(authors)
}

func (a *AuthorsService) CalculateAuthors() (int, error) {
	entries, err := a.dbService.GetAllEntriesIds()
	if err != nil {
		return -1, err
	}
	return len(entries), nil
}

func (a *AuthorsService) DeleteAuthorById(id uint64) error {
	//find all the authors except the one with the given id
	authors, err := a.GetAllAuthors()
	if err != nil {
		fmt.Println("error while getting all authors")
		return err
	}
	var newAuthors []entities.Author
	for _, author := range authors {
		if uint64(author.ID) != id {
			newAuthors = append(newAuthors, author)
		}
	}
	// clean up
	err = a.dbService.CleanUp()
	if err != nil {
		return err
		fmt.Println("error while cleaning up")
	}
	// write all the authors except the one with the given id
	for _, author := range newAuthors {
		err = a.CreateAuthor(&author)
		if err != nil {
			fmt.Println("error while creating author")
			return err
		}
	}
	return nil
}
func (a *AuthorsService) UpdateAuthor(author *entities.Author) error {
	//find all the authors except the one with the given id
	authors, err := a.GetAllAuthors()
	if err != nil {
		fmt.Println("error while getting all authors")
		return err
	}
	var newAuthors []entities.Author
	for _, author := range authors {
		if uint64(author.ID) != uint64(author.ID) {
			newAuthors = append(newAuthors, author)
		}
	}
	// clean up
	err = a.dbService.CleanUp()
	if err != nil {
		return err
		fmt.Println("error while cleaning up")
	}
	// write all the authors except the one with the given id
	for _, author := range newAuthors {
		err = a.CreateAuthor(&author)
		if err != nil {
			fmt.Println("error while creating author")
			return err
		}
	}
	// write the updated author
	err = a.CreateAuthor(author)
	if err != nil {
		fmt.Println("error while creating author")
		return err
	}
	return nil
}
