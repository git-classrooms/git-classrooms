package database

import (
	"time"

	"github.com/google/uuid"
)

type AssignmentDate struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	AssignmentID uuid.UUID  `gorm:"not null;" json:"assignmentID"`
	Assignment   Assignment `json:"-"`

	Description string `gorm:"not null" json:"description"`

	DueDate time.Time `gorm:"not null" json:"dueDate"`
	Closed  bool      `gorm:"default:false" json:"closed"`

	GradingManualRubrics []*ManualGradingRubric `gorm:"many2many:assignment_dates_manual_grading_rubrics;constraint:OnDelete:CASCADE;" json:"-"`

	AssignmentProjectGradingDate []*AssignmentProjectGradingDate `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	JUnitTests []*AssignmentJunitTest `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
