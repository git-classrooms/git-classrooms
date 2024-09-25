package model

type Group struct {
	Name        string
	ID          int
	Description string
	WebUrl      string
	Visibility  Visibility
	Member      []User
	Projects    []Project
}
