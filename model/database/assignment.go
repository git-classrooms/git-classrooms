// Package database contains reference types for representing data with gorm
package database

import (
	"time"

	"github.com/google/uuid"
)

// Assignment is a struct that represents an assignment in the database
type Assignment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ClassroomID uuid.UUID `gorm:"not null;uniqueIndex:idx_unique_classroom_assignmentName" json:"classroomId"`
	Classroom   Classroom `json:"-"`

	TemplateProjectID int        `gorm:"<-:create;not null" json:"templateProjectId"`
	Name              string     `gorm:"not null;uniqueIndex:idx_unique_classroom_assignmentName" json:"name"`
	Description       string     `json:"description"`
	DueDate           *time.Time `json:"dueDate" validate:"optional"`
	Closed            bool       `gorm:"default:false" json:"closed"`

	Projects                      []*AssignmentProjects  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	GradingJUnitAutoGradingActive bool                   `json:"gradingJUnitAutoGradingActive"`
	JUnitTests                    []*AssignmentJunitTest `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	GradingManualRubrics []*ManualGradingRubric `gorm:"many2many:assignment_manual_grading_rubrics;constraint:OnDelete:CASCADE;" json:"-"`
} //@Name Assignment
