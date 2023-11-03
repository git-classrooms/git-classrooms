package model

type Classroom struct {
	Name        string
	ID          int
	Description string
	WebUrl      string
	Visibility  Visibility
	Member      []ClassroomMember
	Projects    []Project
}
