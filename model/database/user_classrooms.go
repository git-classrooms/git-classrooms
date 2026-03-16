package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Role uint8 //@Name Role

const (
	Owner Role = iota
	Moderator
	Student
)

var _ fmt.Stringer = (*Role)(nil)

func (r Role) String() string {
	switch r {
	case Owner:
		return "owner"
	case Moderator:
		return "moderator"
	case Student:
		return "student"
	default:
		return "INVALID ROLE"
	}
}

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

var _ (slog.LogValuer) = (*UserClassrooms)(nil)

// LogValue implements slog.LogValuer.
func (u *UserClassrooms) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("userID", u.UserID),
		slog.Any("teamID", u.TeamID),
		slog.String("role", u.Role.String()),
	)
}
