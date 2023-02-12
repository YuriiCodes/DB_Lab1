package entities

// Post will be the slave file.
type Post struct {
	ID       int64  `json:"id"`
	AuthorId int64  `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}
