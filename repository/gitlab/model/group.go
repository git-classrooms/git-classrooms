package model

type Group struct {
	Name        string
	ID          int
	Description string
	WebURL      string
	Visibility  Visibility
	Member      []User
	Projects    []Project
}
