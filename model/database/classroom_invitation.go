package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type ClassroomInvitationStatus uint8 //@Name ClassroomInvitationStatus

const (
	ClassroomInvitationPending ClassroomInvitationStatus = iota
	ClassroomInvitationAccepted
	ClassroomInvitationRejected
	ClassroomInvitationRevoked
	ClassroomInvitationFailed
)

var _ (fmt.Stringer) = (*ClassroomInvitationStatus)(nil)

func (s ClassroomInvitationStatus) String() string {
	switch s {
	case ClassroomInvitationPending:
		return "pending"
	case ClassroomInvitationAccepted:
		return "accepted"
	case ClassroomInvitationRejected:
		return "rejected"
	case ClassroomInvitationRevoked:
		return "revoked"
	case ClassroomInvitationFailed:
		return "failed"
	default:
		return "INVALID INVITATION STATUS"
	}
}

type ClassroomInvitation struct {
	ID        uuid.UUID                 `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time                 `json:"createdAt"`
	UpdatedAt time.Time                 `json:"updatedAt"`
	Status    ClassroomInvitationStatus `gorm:"not null" json:"status"`

	ClassroomID uuid.UUID `gorm:"not null" json:"-"`
	Classroom   Classroom `json:"classroom"`

	Email      string    `gorm:"not null" json:"email"`
	ExpiryDate time.Time `gorm:"not null" json:"expiryDate"`
} //@Name ClassroomInvitation

var _ (slog.LogValuer) = (*ClassroomInvitation)(nil)

// LogValue implements slog.LogValuer.
func (c *ClassroomInvitation) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", c.ID.String()),
		slog.String("status", c.Status.String()),
		slog.String("classroomID", c.ClassroomID.String()),
	)
}
