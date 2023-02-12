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
	postsServ, err := postsService.NewPostService()
	if err != nil {
		fmt.Println("error while creating post service: ", err.Error())
		return
	}

	authorsServ, err := authorsService.NewAuthorService(postsServ)
	if err != nil {
		fmt.Println("error while creating author service: ", err.Error())
		return
	}

	// 1. create 5 authors:
	authorsServ.CreateAuthor(entities.Author{Name: "Yurii"})
	authorsServ.CreateAuthor(entities.Author{Name: "Vlad"})
	authorsServ.CreateAuthor(entities.Author{Name: "Ivan"})
	authorsServ.CreateAuthor(entities.Author{Name: "Oleg"})
	authorsServ.CreateAuthor(entities.Author{Name: "Vitalii"})

	// 2. for first author, add 1 post, for second author - 2 posts, for third author - 3 posts
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 1, Title: "Post 1", Content: "Content 1"})

	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 2, Title: "Post 2", Content: "Content 2"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 2, Title: "Post 3", Content: "Content 3"})

	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 3, Title: "Post 4", Content: "Content 4"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 3, Title: "Post 5", Content: "Content 5"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 3, Title: "Post 6", Content: "Content 6"})

	// 3. print system info for authors & posts:
	authorsServ.PrintSystemInfo()
	postsServ.PrintSystemInfo()

	// 4. remove master file with 2 posts (the second author)
	authorsServ.DeleteAuthor(2)

	// 5. remove one post from the third author
	authorsServ.RemovePostFromAuthor(4)

	// 6. print system info for authors & posts:
	authorsServ.PrintSystemInfo()
	postsServ.PrintSystemInfo()

	// 7. Add one more master & add post by him:
	authorsServ.CreateAuthor(entities.Author{Name: "Steve"})
	authorsServ.AddPostByAuthor(entities.Post{AuthorId: 6, Title: "Post 7", Content: "Content 7"})

	// 8. print system info for authors & posts:
	authorsServ.PrintSystemInfo()
	postsServ.PrintSystemInfo()

	// 9. Update one master & post by him:
	authorsServ.UpdateAuthor(entities.Author{ID: 6, Name: "Steve Jobs"})
	authorsServ.UpdatePostFromAuthor(entities.Post{ID: 7, AuthorId: 6, Title: "Post 7", Content: "Content 7"})

	// 10. print system info for authors & posts:
	authorsServ.PrintSystemInfo()
	postsServ.PrintSystemInfo()
}
