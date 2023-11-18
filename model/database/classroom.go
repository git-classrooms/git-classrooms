// Package database contains reference types for representing data with gorm
package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role uint8

const (
	Owner Role = iota
	Moderator
	Student
)

// Classroom is a struct that represents a classroom in the database
type Classroom struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	OwnerID     int
	Owner       User
	Description string
	GroupID     int              `gorm:"<-:create"`
	Member      []UserClassrooms `gorm:"foreignKey:ClassroomID"`
	Assignments []Assignment
}

type UserClassrooms struct {
	UserID      int `gorm:"primaryKey;autoIncrement:false"`
	User        User
	ClassroomID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Classroom   Classroom
	Role        Role
}
