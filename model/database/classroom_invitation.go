package database

import (
	"github.com/google/uuid"
	"time"
)

type InvitationStatus uint8

const (
	InvitationPending Role = iota
	InvitationAccepted
	InvitationDeclined
	InvitationRevoked
)

type ClassroomInvitation struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      InvitationStatus `gorm:"not null"`
	Email       string           `gorm:"not null"`
	Enabled     bool             `gorm:"not null"`
	ExpiryDate  time.Time        `gorm:"not null"`
	ClassroomID uuid.UUID        `gorm:"not null"`
	Classroom   Classroom
}
