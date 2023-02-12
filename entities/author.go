package entities

import "strconv"

// Author will be the master file.
type Author struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//define a string method:
func (a Author) String() string {
	return strconv.Itoa(a.ID) + " " + a.Name + " " + a.Email + " " + a.Password
}
