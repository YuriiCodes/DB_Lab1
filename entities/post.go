package entities

// Post will be the slave file.
type Post struct {
	ID       int    `json:"id"`
	AuthorId int    `json:"authorId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}
