package database

import (
	"github.com/google/uuid"
	"time"
)

type ClassroomInvitationStatus uint8

const (
	ClassroomInvitationPending ClassroomInvitationStatus = iota
	ClassroomInvitationAccepted
	ClassroomInvitationDeclined
	ClassroomInvitationRevoked
)

type ClassroomInvitation struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      ClassroomInvitationStatus `gorm:"not null"`
	ClassroomID uuid.UUID                 `gorm:"not null"`
	Classroom   Classroom
	Email       string    `gorm:"not null"`
	Enabled     bool      `gorm:"not null"`
	ExpiryDate  time.Time `gorm:"not null"`
}
