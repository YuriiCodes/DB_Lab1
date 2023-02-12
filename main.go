package main

import (
	"fmt"

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

	fmt.Println("update-m <author_id> <name> <email> <password> - update author by id")
	fmt.Println("update-s <post_id> <author_id> <views> <image> <title> <content> - update post by id\n")

	fmt.Println("insert-m <name> <email> <password> - insert author")
	fmt.Println("insert-s <author_id> <views> <image> <title> <content> - insert post\n")

	fmt.Println("calc-m - calculate number of authors")
	fmt.Println("calc-s - calculate number of posts\n")

	fmt.Println("ut-m - utility function for printing all the fields (including the system ones) of the author table")
	fmt.Println("ut-s - utility function for printing all the fields (including the system ones) of the post table\n")

	fmt.Println("exit - exit the program")

}
func main() {
	greet()

	postsSer, _ := postsService.NewPostService()
	post := entities.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Title",
		Content:  "Content",
	}

	postUpdated := entities.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "TitleUpdated",
		Content:  "ContentUpdated",
	}

	err := postsSer.CreatePost(post)
	if err != nil {
		fmt.Println("error while creating post: ", err.Error())
	}

	post1, err := postsSer.GetPostById(1)
	if err != nil {
		fmt.Println("error while getting post by id: ", err.Error())
	}
	fmt.Println("post1: ", post1)

	err = postsSer.UpdatePost(postUpdated)
	if err != nil {
		fmt.Println("error while updating post: ", err.Error())
	}

	post1Updated, err := postsSer.GetPostById(1)
	if err != nil {
		fmt.Println("error while getting post by id: ", err.Error())
	}
	fmt.Println("post1Updated: ", post1Updated)

	//postsSer.DeletePost(1)
	err = postsSer.DeleteAllPostsByAuthorId(1)
	if err != nil {
		fmt.Println("error while deleting all posts by author id: ", err.Error())
	}
}
