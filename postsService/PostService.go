package postsService

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"BD_Lab1/entities"
)

//define PostService struct:
type PostService struct {
	posts      []entities.Post
	indexTable entities.IndexTable

	postsFile      *os.File
	indexTableFile *os.File
}

//define a constructor:
func NewPostService() (*PostService, error) {
	posts := make([]entities.Post, 0)
	indexTable := entities.IndexTable{}

	mainFile, err := os.OpenFile("posts.fl.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	main, err := mainFile.Stat()
	if err != nil {
		return nil, err
	}
	if main.Size() == 0 {
		fmt.Println("posts.fl is empty")
	} else {
		fmt.Println("posts.fl is not empty")
		// read from file to posts:
		err = json.NewDecoder(mainFile).Decode(&posts)
		if err != nil {
			return nil, errors.New("error while reading posts.fl: " + err.Error())
		}
	}
	fmt.Println("posts.fl is read")

	indexFile, err := os.OpenFile("posts.ind.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	indexStat, err := indexFile.Stat()
	if err != nil {
		return nil, err
	}

	// Create the index table if it is empty:
	if indexStat.Size() == 0 {
		fmt.Println("posts.ind is empty. Creating new index table...")
		// create new index table:
		indexTable = entities.IndexTable{
			Uid:  0,
			Rows: make([]entities.IndexTableRow, 0),
		}
		// write index table to file:
		err = json.NewEncoder(indexFile).Encode(indexTable)
		if err != nil {
			return nil, errors.New("error while writing index table to posts.ind: " + err.Error())
		}
	} else {
		fmt.Println("posts.ind is not empty")
		// read from file to posts:
		err = json.NewDecoder(indexFile).Decode(&indexTable)
		if err != nil {
			return nil, errors.New("error while reading posts.fl: " + err.Error())
		}
	}
	fmt.Println(posts)
	return &PostService{
		posts:      posts,
		indexTable: indexTable,

		postsFile:      mainFile,
		indexTableFile: indexFile,
	}, nil
}

func (p *PostService) GetPosts() []entities.Post {
	return p.posts
}
func (p *PostService) CreatePost(post entities.Post) error {
	// update ID to be unique:
	post.ID = p.indexTable.Uid + 1
	p.posts = append(p.posts, post)

	// write to file:
	err := p.postsFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = p.postsFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of posts.fl: ", err.Error())
		return err
	}
	err = json.NewEncoder(p.postsFile).Encode(p.posts)
	if err != nil {
		fmt.Println("error while writing posts to posts.fl: ", err.Error())
		return err
	}

	// update index table:
	p.indexTable.Uid += 1
	p.indexTable.Rows = append(p.indexTable.Rows, entities.IndexTableRow{
		UID:        p.indexTable.Uid,
		NumInArray: len(p.posts) - 1,
	})

	// write to file: (we have to update the existing json in file, not append the whole new one):
	err = p.indexTableFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = p.indexTableFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of posts.ind: ", err.Error())
		return err
	}
	err = json.NewEncoder(p.indexTableFile).Encode(p.indexTable)
	if err != nil {
		fmt.Println("error while writing index table to posts.ind: ", err.Error())
		return err
	}
	return nil
}

func (p *PostService) GetPostById(id int) (entities.Post, error) {
	// search in index table:
	for _, row := range p.indexTable.Rows {
		if row.UID == id {
			return p.posts[row.NumInArray], nil
		}
	}
	return entities.Post{
		ID:       -1,
		Title:    "",
		Content:  "",
		AuthorId: -1,
	}, errors.New("no post with such id")
}

func (p *PostService) UpdatePost(post entities.Post) error {
	// search in index table:
	for _, row := range p.indexTable.Rows {
		if row.UID == post.ID {
			p.posts[row.NumInArray] = post
			// write to file:
			err := p.postsFile.Truncate(0)
			if err != nil {
				return err
			}
			_, err = p.postsFile.Seek(0, 0)
			if err != nil {
				fmt.Println("error while seeking to the beginning of posts.fl: ", err.Error())
				return err
			}
			err = json.NewEncoder(p.postsFile).Encode(p.posts)
			if err != nil {
				fmt.Println("error while writing posts to posts.fl: ", err.Error())
				return err
			}
			return nil
		}
	}
	return errors.New("no post with such id")
}

func (p *PostService) DeletePost(id int) error {
	// search for post in index table and update indexes of all posts after it (decrease by 1):
	var index int
	var found bool
	for _, row := range p.indexTable.Rows {
		if row.UID == id {
			index = row.NumInArray
			found = true
			break
		}
	}
	if !found {
		return errors.New("no post with such id")
	}
	// iterate over each row after index and decrease NumInArray by 1:
	for i := index + 1; i < len(p.indexTable.Rows); i++ {
		p.indexTable.Rows[i].NumInArray -= 1
	}
	// update index table:
	p.indexTable.Rows = append(p.indexTable.Rows[:index], p.indexTable.Rows[index+1:]...)
	//save index table to file:
	err := p.indexTableFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = p.indexTableFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of posts.ind: ", err.Error())
		return err
	}
	err = json.NewEncoder(p.indexTableFile).Encode(p.indexTable)
	if err != nil {
		fmt.Println("error while writing index table to posts.ind: ", err.Error())
		return err
	}

	// update posts:
	p.posts = append(p.posts[:index], p.posts[index+1:]...)
	// save posts to file:
	err = p.postsFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = p.postsFile.Seek(0, 0)
	if err != nil {
		fmt.Println("error while seeking to the beginning of posts.fl: ", err.Error())
		return err
	}
	err = json.NewEncoder(p.postsFile).Encode(p.posts)
	if err != nil {
		fmt.Println("error while writing posts to posts.fl: ", err.Error())
		return err
	}
	return nil
}

func (p *PostService) GetAllPosts() []entities.Post {
	return p.posts
}
func (p *PostService) DeleteAllPostsByAuthorId(id int) error {
	// iterate over posts using range function and collect all posts with author id = id
	var IDs []int
	for _, post := range p.posts {
		if post.AuthorId == id {
			IDs = append(IDs, post.ID)
		}
	}

	// then delete them using DeletePost method
	for _, id := range IDs {
		err := p.DeletePost(id)
		if err != nil {
			return err
		}
	}
	return nil
}