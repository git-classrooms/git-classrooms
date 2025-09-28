package database

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name    string `gorm:"not null" json:"name"`
	GroupID int    `gorm:"<-:create;not null" json:"groupId"`

	ClassroomID uuid.UUID `gorm:"<-:create;type:uuid;not null" json:"-"`
	Classroom   Classroom `json:"-"`

	Member             []*UserClassrooms     `gorm:"foreignKey:TeamID;constraint:OnDelete:SET NULL;" json:"-"`
	AssignmentProjects []*AssignmentProjects `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
} //@Name Team

var _ (slog.LogValuer) = (*Team)(nil)

// LogValue implements slog.LogValuer.
func (t *Team) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", t.ID.String()),
		slog.String("name", t.Name),
		slog.String("classroomID", t.ClassroomID.String()),
	)
}
