package database

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type ManualGradingRubric struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`

	ClassroomID uuid.UUID `gorm:"not null" json:"-"`
	Classroom   Classroom `gorm:";" json:"-"`

	MaxScore int `gorm:"not null" json:"maxScore"`

	Assignments []*Assignment          `gorm:"many2many:assignment_manual_grading_rubrics;constraint:OnDelete:CASCADE;" json:"-"`
	Results     []*ManualGradingResult `gorm:"foreignKey:RubricID;constraint:OnDelete:CASCADE;" json:"-"`
} //@Name ManualGradingRubric

var _ (slog.LogValuer) = (*ManualGradingRubric)(nil)

// LogValue implements slog.LogValuer.
func (r *ManualGradingRubric) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", r.ID.String()),
		slog.String("name", r.Name),
		slog.String("classroomID", r.ClassroomID.String()),
		slog.Int("maxScore", r.MaxScore),
	)
}
