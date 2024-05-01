package database

import (
	"time"

	"github.com/google/uuid"
)

type ClassroomInvitationStatus uint8 //@Name ClassroomInvitationStatus

const (
	ClassroomInvitationPending ClassroomInvitationStatus = iota
	ClassroomInvitationAccepted
	ClassroomInvitationDeclined
	ClassroomInvitationRevoked
)

type ClassroomInvitation struct {
	ID          uuid.UUID                 `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
	Status      ClassroomInvitationStatus `gorm:"not null" json:"status"`
	ClassroomID uuid.UUID                 `gorm:"not null" json:"-"`
	Classroom   Classroom                 `json:"-"`
	Email       string                    `gorm:"not null" json:"email"`
	ExpiryDate  time.Time                 `gorm:"not null" json:"expiryDate"`
} //@Name ClassroomInvitation
