package entities

import "strconv"

// Author will be the master file.
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//define a string method:
func (a Author) String() string {
	return strconv.Itoa(a.ID) + " " + a.Name
}
