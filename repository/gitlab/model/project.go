package model

type Visibility int

const (
	Private Visibility = iota
	Internal
	Public
)

type Project struct {
	Name        string     `json:"name"`
	ID          int        `json:"id"`
	Visibility  Visibility `json:"visibility"`
	WebUrl      string     `json:"webUrl"`
	Description string     `json:"description"`
	Owner       *User      `json:"owner"`
	Member      []User     `json:"member"`
}
