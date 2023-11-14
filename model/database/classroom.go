package database

import (
	"gorm.io/gorm"
	"time"
)

type Classroom struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Owner       User
	Description string
	WebUrl      string
	Member      []User `gorm:"many2many:user_classrooms;"`
	Assignments []Assignment
}
