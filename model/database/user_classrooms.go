package database

import "github.com/google/uuid"

type Role uint8 //@Name Role

const (
	Owner Role = iota
	Moderator
	Student
)

// UserClassrooms is a struct that represents the relationship between a user and a classroom
type UserClassrooms struct {
	UserID int  `gorm:"primaryKey;autoIncrement:false;not null" json:"-"`
	User   User `gorm:";" json:"user"`

	ClassroomID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"-"`
	Classroom   Classroom `gorm:";" json:"classroom"`

	TeamID *uuid.UUID `gorm:"type:uuid;index" json:"-"`
	Team   *Team      `json:"team" validate:"optional"`
	Role   Role       `gorm:"not null" json:"role"`
} //@Name UserClassrooms
