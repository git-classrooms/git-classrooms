package database

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type AssignmentJunitTest struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Name         string     `gorm:"not null;uniqueIndex:idx_unique_assignment_assignmentjunittestName" json:"name"`
	AssignmentID uuid.UUID  `gorm:"not null;uniqueIndex:idx_unique_assignment_assignmentjunittestName" json:"-"`
	Assignment   Assignment `json:"-"`

	Score int `gorm:"not null" json:"score"`
} //@Name AssignmentJunitTest

var _ (slog.LogValuer) = (*AssignmentJunitTest)(nil)

// LogValue implements slog.LogValuer.
func (a *AssignmentJunitTest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", a.ID.String()),
		slog.String("name", a.Name),
		slog.String("assignmentID", a.AssignmentID.String()),
		slog.Int("score", a.Score),
	)
}
