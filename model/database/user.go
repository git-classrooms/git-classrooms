// Package database contains reference types for representing data with gorm
package database

import (
	"gorm.io/gorm"
	"time"
)

// User is the representation of the user in database
type User struct {
	ID                     int                   `gorm:"primary_key;autoIncrement:false" json:"id"`
	GitlabEmail            string                `gorm:"unique;not null" json:"gitlab_email"`
	Name                   string                `gorm:"not null" json:"name"`
	CreatedAt              time.Time             `json:"-"`
	UpdatedAt              time.Time             `json:"-"`
	DeletedAt              gorm.DeletedAt        `gorm:"index" json:"-"`
	Classrooms             []*UserClassrooms     `gorm:"foreignKey:UserID" json:"-"`
	AssignmentRepositories []*AssignmentProjects `json:"-"`
}
