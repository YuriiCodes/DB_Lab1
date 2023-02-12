package main

import (
	"fmt"

	authorsService2 "BD_Lab1/authorsService"
	"BD_Lab1/entities"
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
	authorsService, err := authorsService2.NewAuthorsService("authors")
	if err != nil {
		panic(err)
	}

	author1 := entities.Author{
		Name:     "Yurii",
		Email:    "Yurii@gmail.com",
		Password: "Mello",
	}

	author2 := entities.Author{
		Name:     "Yurii2",
		Email:    "Yurii2@gmail.com",
		Password: "Hello",
	}

	author2Updated := entities.Author{
		ID:    2,
		Name:  "Yurii2",
		Email: "Yurii2IUPDATED@gmail.com",
	}

	authorsService.CreateAuthor(&author1)
	authorsService.CreateAuthor(&author2)

	//authorsService.PrintAllInfo()

	authorsLength, err := authorsService.CalculateAuthors()
	if err != nil {
		panic(err)
	}
	fmt.Println("Authors length: ", authorsLength)

	authorsService.PrintAllAuthors()
	authorsService.UpdateAuthor(&author2Updated)

	authors, err := authorsService.GetAllAuthors()
	fmt.Println(authors)

	authorsService.Close()

}
