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
	UserID        int        `gorm:"primaryKey;autoIncrement:false;not null" json:"-"`
	User          User       `json:"user"`
	ClassroomID   uuid.UUID  `gorm:"type:uuid;primaryKey;not null" json:"-"`
	Classroom     Classroom  `json:"classroom"`
	TeamID        *uuid.UUID `gorm:"type:uuid;index" json:"-"`
	Team          *Team      `json:"team" validate:"optional"`
	Role          Role       `gorm:"not null" json:"role"`
	LeftClassroom bool       `gorm:"not null;default:false" json:"leftClassroom"` //TODO: Ich sehe noch keinen Vorteil/Nutzen darin, den die connection einfach zu entfernen wenn user via gitlab leaved oder gekickt wird
	LeftTeam      bool       `gorm:"not null;default:false" json:"leftTeam"`
} //@Name UserClassrooms
