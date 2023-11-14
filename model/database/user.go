package database

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                     int
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
	Username               string
	Name                   string
	WebUrl                 string
	Email                  string
	Classrooms             []Classroom `gorm:"many2many:user_classrooms;"`
	Assignments            []Assignment
	AssignmentRepositories []AssignmentRepository
}
