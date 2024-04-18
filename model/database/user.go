// Package database contains reference types for representing data with gorm
package database

import (
	"time"

	"gorm.io/gorm"
)

// User is the representation of the user in database
type User struct {
	ID              int               `gorm:"primary_key;autoIncrement:false" json:"id"`
	GitlabEmail     string            `gorm:"unique;not null" json:"gitlabEmail"`
	Name            string            `gorm:"not null" json:"name"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
	OwnedClassrooms []*Classroom      `gorm:"foreignKey:OwnerID" json:"-"`
	Classrooms      []*UserClassrooms `gorm:"foreignKey:UserID" json:"-"`
}
