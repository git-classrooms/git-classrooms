package model

type Visibility int

const (
	Private Visibility = iota
	Internal
	Public
)

type Project struct {
	Name          string
	ID            int
	Visibility    Visibility
	WebUrl        string
	Description   string
	Owner         *User
	DefaultBranch string
	Members       []User
}
