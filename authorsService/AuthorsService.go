package authorsService

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"BD_Lab1/entities"
	"BD_Lab1/postsService"
)

type AuthorsService struct {
	authors    []entities.Author
	indexTable entities.IndexTable

	authorsFile    *os.File
	indexTableFile *os.File

	postService *postsService.PostService
}

// define a constructor
func NewAuthorService(postServ *postsService.PostService) (*AuthorsService, error) {
	authors := make([]entities.Author, 0)
	indexTable := entities.IndexTable{}

	mainFile, err := os.OpenFile("authors.fl.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	main, err := mainFile.Stat()
	if err != nil {
		return nil, err
	}
	if main.Size() == 0 {
		fmt.Println("authors.fl is empty")
	} else {
		fmt.Println("authors.fl is not empty")
		// read from file to posts:
		err = json.NewDecoder(mainFile).Decode(&authors)
		if err != nil {
			return nil, errors.New("error while reading posts.fl: " + err.Error())
		}
	}
	fmt.Println("authors.fl is read")

	indexFile, err := os.OpenFile("authors.ind.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	indexStat, err := indexFile.Stat()
	if err != nil {
		return nil, err
	}

	// Create the index table if it is empty:
	if indexStat.Size() == 0 {
		fmt.Println("authors.ind is empty. Creating new index table...")
		// create new index table:
		indexTable = entities.IndexTable{
			Uid:  0,
			Rows: make([]entities.IndexTableRow, 0),
		}
		// write index table to file:
		err = json.NewEncoder(indexFile).Encode(indexTable)
		if err != nil {
			return nil, errors.New("error while writing index table to authors.ind: " + err.Error())
		}
	} else {
		fmt.Println("authots.ind is not empty")
		// read from file to posts:
		err = json.NewDecoder(indexFile).Decode(&indexTable)
		if err != nil {
			return nil, errors.New("error while reading authors.fl: " + err.Error())
		}
	}
	fmt.Println(authors)
	return &AuthorsService{
		authors:    authors,
		indexTable: indexTable,

		authorsFile:    mainFile,
		indexTableFile: indexFile,
		postService:    postServ,
	}, nil
}

func (a *AuthorsService) CreateAuthor(author entities.Author) (entities.Author, error) {
	// update ID to be unique:
	author.ID = a.indexTable.Uid + 1
	a.authors = append(a.authors, author)

	// write to file:
	err := a.authorsFile.Truncate(0)
	if err != nil {
		return entities.Author{ID: -1}, err
	}
	_, err = a.authorsFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of authors.fl: ", err.Error())
		return entities.Author{ID: -1}, err
	}
	err = json.NewEncoder(a.authorsFile).Encode(a.authors)
	if err != nil {
		fmt.Println("error while writing posts to authors.fl: ", err.Error())
		return entities.Author{ID: -1}, err
	}

	// update index table:
	a.indexTable.Uid += 1
	a.indexTable.Rows = append(a.indexTable.Rows, entities.IndexTableRow{
		UID:        a.indexTable.Uid,
		NumInArray: len(a.authors) - 1,
	})

	// write to file: (we have to update the existing json in file, not append the whole new one):
	err = a.indexTableFile.Truncate(0)
	if err != nil {
		return entities.Author{ID: -1}, err
	}
	_, err = a.indexTableFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of authors.ind: ", err.Error())
		return entities.Author{ID: -1}, err
	}
	err = json.NewEncoder(a.indexTableFile).Encode(a.indexTable)
	if err != nil {
		fmt.Println("error while writing index table to authors.ind: ", err.Error())
		return entities.Author{ID: -1}, err
	}
	return author, nil
}

func (a *AuthorsService) GetAuthorById(id int) (entities.Author, error) {
	// search in index table:
	for _, row := range a.indexTable.Rows {
		if row.UID == id {
			return a.authors[row.NumInArray], nil
		}
	}
	return entities.Author{
		ID:   -1,
		Name: "",
	}, errors.New("no post with such id")
}

func (a *AuthorsService) UpdateAuthor(author entities.Author) (entities.Author, error) {
	// search in index table:
	for _, row := range a.indexTable.Rows {
		if row.UID == author.ID {
			a.authors[row.NumInArray] = author
			// write to file:
			err := a.authorsFile.Truncate(0)
			if err != nil {
				return entities.Author{ID: -1}, err
			}
			_, err = a.authorsFile.Seek(0, 0)
			if err != nil {
				fmt.Println("error while seeking to the beginning of authors.fl: ", err.Error())
				return entities.Author{ID: -1}, err
			}
			err = json.NewEncoder(a.authorsFile).Encode(a.authors)
			if err != nil {
				fmt.Println("error while writing posts to authors.fl: ", err.Error())
				return entities.Author{ID: -1}, err
			}
			return author, nil
		}
	}
	return entities.Author{ID: -1}, errors.New("no author with such id")
}

func (a *AuthorsService) DeleteAuthor(id int) (entities.Author, error) {
	// to delete the author, we first have to delete all his posts:
	err := a.postService.DeleteAllPostsByAuthorId(id)
	if err != nil {
		return entities.Author{ID: -1}, err
	}

	// search for post in index table and update indexes of all posts after it (decrease by 1):
	var index int
	var found bool
	for _, row := range a.indexTable.Rows {
		if row.UID == id {
			index = row.NumInArray
			found = true
			break
		}
	}
	if !found {
		return entities.Author{ID: -1}, errors.New("no author with such id")
	}
	// iterate over each row after index and decrease NumInArray by 1:
	for i := index + 1; i < len(a.indexTable.Rows); i++ {
		a.indexTable.Rows[i].NumInArray -= 1
	}
	// update index table:
	a.indexTable.Rows = append(a.indexTable.Rows[:index], a.indexTable.Rows[index+1:]...)
	//save index table to file:
	err = a.indexTableFile.Truncate(0)
	if err != nil {
		return entities.Author{ID: -1}, err
	}
	_, err = a.indexTableFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of authors.ind: ", err.Error())
		return entities.Author{ID: -1}, err
	}
	err = json.NewEncoder(a.indexTableFile).Encode(a.indexTable)
	if err != nil {
		fmt.Println("error while writing index table to authors.ind: ", err.Error())
		return entities.Author{ID: -1}, err
	}

	// update posts:
	a.authors = append(a.authors[:index], a.authors[index+1:]...)
	// save posts to file:
	err = a.authorsFile.Truncate(0)
	if err != nil {
		return entities.Author{ID: -1}, err
	}
	_, err = a.authorsFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of author.fl: ", err.Error())
		return entities.Author{ID: -1}, err
	}
	err = json.NewEncoder(a.authorsFile).Encode(a.authors)
	if err != nil {
		fmt.Println("error while writing posts to posts.fl: ", err.Error())
		return entities.Author{ID: -1}, err
	}
	return entities.Author{ID: id}, nil
}

func (a *AuthorsService) GetAllAuthors() ([]entities.Author, error) {
	if len(a.authors) == 0 {
		return nil, errors.New("no authors")
	}
	return a.authors, nil
}

func (a *AuthorsService) GetNumberOfAuthors() int {
	return len(a.authors)
}

func (a *AuthorsService) GetPostsByAuthorId(id int) ([]entities.Post, error) {
	return a.postService.GetPostsByAuthorId(id)
}

func (a *AuthorsService) PrintSystemInfo() error {
	a.postService.PrintSystemInfo()
	// print all the authors
	authors, err := a.GetAllAuthors()
	if err != nil {
		return err
	}
	fmt.Println("Authors:")
	for _, author := range authors {
		// print author info and posts by him:
		fmt.Println("Author: ", author)
		posts, err := a.GetPostsByAuthorId(author.ID)
		if err != nil {
			fmt.Println("error while getting posts by author: ", err.Error())
		}
		fmt.Println("Posts by author: ", posts)
	}
	return nil
}

func (a *AuthorsService) AddPostByAuthor(post entities.Post) (entities.Post, error) {
	// search for author and throw error if author not found:
	_, err := a.GetAuthorById(post.AuthorId)
	if err != nil {
		return entities.Post{ID: -1}, err
	}
	// add post to post service:
	createdPost, err := a.postService.CreatePost(post)
	if err != nil {
		return entities.Post{ID: -1}, err
	}
	return createdPost, nil
}

func (a *AuthorsService) RemovePostFromAuthor(postId int) (entities.Post, error) {
	removedPost, err := a.postService.DeletePost(postId)
	if err != nil {
		return entities.Post{ID: -1}, err
	}
	return removedPost, nil
}

func (a *AuthorsService) UpdatePostFromAuthor(post entities.Post) (entities.Post, error) {
	updatedPost, err := a.postService.UpdatePost(post)
	if err != nil {
		return entities.Post{ID: -1}, err
	}
	return updatedPost, nil
}

func (a *AuthorsService) Close() error {
	err := a.postService.Close()
	if err != nil {
		return err
	}
	err = a.authorsFile.Close()

	if err != nil {
		return err
	}
	err = a.indexTableFile.Close()
	if err != nil {
		return err
	}
	return nil

}
