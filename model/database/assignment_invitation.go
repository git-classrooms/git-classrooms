package database

import (
	"github.com/google/uuid"
	"time"
)

type AssignmentInvitationStatus uint8

const (
	AssignmentInvitationPending AssignmentInvitationStatus = iota
	AssignmentInvitationAccepted
	AssignmentInvitationDeclined
	AssignmentInvitationRevoked
)

type AssignmentInvitation struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Status       AssignmentInvitationStatus `gorm:"not null"`
	AssignmentID uuid.UUID                  `gorm:"not null"`
	Assignment   Classroom
	UserID       int `gorm:"not null"`
	User         User
	Enabled      bool      `gorm:"not null"`
	ExpiryDate   time.Time `gorm:"not null"`
}
