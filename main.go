package main

import (
	"fmt"

	"BD_Lab1/authorsService"
	"BD_Lab1/entities"
	"BD_Lab1/postsService"
)

func greet() {
	fmt.Println("Welcome to Lab#1 of BD course!")
	fmt.Println("This is a one-to-many relationship example between users and posts, using only .fl files and index tables.")
	fmt.Println("Made by: Yurii Pidlisnyi, K-26")
	fmt.Println("Here are available commands 👇")

	fmt.Println("get-m <author_id> - get author by id")
	fmt.Println("get-s <post_id> - get post by id\n")

	fmt.Println("del-m <author_id> - delete author by id")
	fmt.Println("del-s <post_id> - delete post by id\n")

	fmt.Println("update-m <author_id> <name>  - update author by id")
	fmt.Println("update-s <post_id> <title> <content> - update post by id\n")

	fmt.Println("insert-m <name>  - insert author")
	fmt.Println("insert-s <author_id> <title> <content> - insert post\n")

	fmt.Println("calc-m - calculate number of authors")
	fmt.Println("calc-s - calculate number of posts\n")

	fmt.Println("ut-m - utility function for printing all the fields (including the system ones) of the author table")
	fmt.Println("ut-s - utility function for printing all the fields (including the system ones) of the post table\n")

	fmt.Println("exit - exit the program")

}
func main() {
	greet()

	// initialize the post service:
	postsServ, err := postsService.NewPostService()
	if err != nil {
		fmt.Println("error while creating post service: ", err.Error())
		return
	}

	// initialize the author service:
	authorsServ, err := authorsService.NewAuthorService(postsServ)
	if err != nil {
		fmt.Println("error while creating author service: ", err.Error())
		return
	}

	// 1. create 5 authors:
	author1, err := authorsServ.CreateAuthor(entities.Author{Name: "Yurii"})
	if err != nil {
		fmt.Println("error while creating author: ", err.Error())
		return
	}
	author2, err := authorsServ.CreateAuthor(entities.Author{Name: "Vlad"})
	if err != nil {
		fmt.Println("error while creating author: ", err.Error())
		return
	}

	author3, err := authorsServ.CreateAuthor(entities.Author{Name: "Ivan"})
	if err != nil {
		fmt.Println("error while creating author: ", err.Error())
		return
	}
	_, err = authorsServ.CreateAuthor(entities.Author{Name: "Oleg"})
	if err != nil {
		fmt.Println("error while creating author: ", err.Error())
		return
	}
	_, err = authorsServ.CreateAuthor(entities.Author{Name: "Vitalii"})
	if err != nil {
		fmt.Println("error while creating author: ", err.Error())
		return
	}

	// 2. for first author, add 1 post, for second author - 2 posts, for third author - 3 posts
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: author1.ID, Title: "Post 1", Content: "Content 1"})

	authorsServ.AddPostByAuthor(entities.Post{AuthorId: author2.ID, Title: "Post 2", Content: "Content 2"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: author2.ID, Title: "Post 3", Content: "Content 3"})

	authorsServ.AddPostByAuthor(entities.Post{AuthorId: author3.ID, Title: "Post 4", Content: "Content 4"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: author3.ID, Title: "Post 5", Content: "Content 5"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: author3.ID, Title: "Post 6", Content: "Content 6"})

	// 3. print system info for authors & posts:
	err = authorsServ.PrintSystemInfo()
	if err != nil {
		fmt.Println("error while printing system info: ", err.Error())
		return
	}
	postsServ.PrintSystemInfo()

	// 4. remove master file with 2 posts (the second author)
	_, err = authorsServ.DeleteAuthor(author2.ID)
	if err != nil {
		fmt.Println("error while deleting author: ", err.Error())
		return
	}

	// 5. remove one post from the third author
	posts, err := authorsServ.GetPostsByAuthorId(author3.ID)
	if err != nil {
		fmt.Println("error while getting posts by author id: ", err.Error())
		return
	}
	authorsServ.RemovePostFromAuthor(posts[0].ID)

	// 6. print system info for authors & posts:
	err = authorsServ.PrintSystemInfo()
	if err != nil {
		fmt.Println("error while printing system info: ", err.Error())
		return
	}
	postsServ.PrintSystemInfo()

	// 7. Add one more master & add post by him:
	author6, err := authorsServ.CreateAuthor(entities.Author{Name: "Steve"})
	if err != nil {
		fmt.Println("error while creating author: ", err.Error())
		return
	}
	postByAuthor6, err := authorsServ.AddPostByAuthor(entities.Post{AuthorId: author6.ID, Title: "Post 7", Content: "Content 7"})
	if err != nil {
		fmt.Println("error while adding post by author: ", err.Error())
		return
	}

	// 8. print system info for authors & posts:
	err = authorsServ.PrintSystemInfo()
	if err != nil {
		fmt.Println("error while printing system info: ", err.Error())
		return
	}
	postsServ.PrintSystemInfo()

	// 9. Update one master & post by him:
	author6Updated, err := authorsServ.UpdateAuthor(entities.Author{ID: author6.ID, Name: "Steve Jobs"})
	if err != nil {
		fmt.Println("error while updating author: ", err.Error())
		return
	}
	authorsServ.UpdatePostFromAuthor(entities.Post{ID: postByAuthor6.ID, AuthorId: author6Updated.ID, Title: "Post 7 UPDATED", Content: "Content 7 UPDATED"})

	// 10. print system info for authors & posts:
	err = authorsServ.PrintSystemInfo()
	if err != nil {
		fmt.Println("error while printing system info: ", err.Error())
		return
	}
	postsServ.PrintSystemInfo()
}
