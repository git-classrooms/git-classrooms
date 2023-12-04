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

// ClassRoomDTO is the data transfer object representing a user
type ClassRoomDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	OwnerID     int       `json:"ownerId"`
	Description string    `json:"description"`
	GroupID     int       `json:"groupId"`
}

// UserClassrooms is a struct that represents the relationship between a user and a classroom
type UserClassrooms struct {
	UserID      int `gorm:"primaryKey;autoIncrement:false"`
	User        User
	ClassroomID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Classroom   Classroom
	Role        Role
}

// UserClassroomDTO is the data transfer object representing a user
type UserClassroomDTO struct {
	UserID      int    `json:"userId"`
	ClassroomID string `json:"classroomId"`
	Role        Role   `json:"role"`
}
