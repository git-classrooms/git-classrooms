package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Role uint8 //@Name Role

const (
	Owner Role = iota
	Moderator
	Student
)

// String implements [fmt.Stringer].
func (r *Role) String() string {
	switch *r {
	case Owner:
		return "Owner"
	case Moderator:
		return "Moderator"
	case Student:
		return "Student"
	default:
		return "Invalid Role"
	}
}

var _ fmt.Stringer = (*Role)(nil)

// UserClassrooms is a struct that represents the relationship between a user and a classroom
type UserClassrooms struct {
	UserID int  `gorm:"primaryKey;autoIncrement:false;not null" json:"-"`
	User   User `gorm:";" json:"user"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"-"`

	ClassroomID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"-"`
	Classroom   Classroom `gorm:";" json:"classroom"`

	TeamID *uuid.UUID `gorm:"type:uuid;index" json:"-"`
	Team   *Team      `json:"team" validate:"optional"`
	Role   Role       `gorm:"not null" json:"role"`
} //@Name UserClassrooms
