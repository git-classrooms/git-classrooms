package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AssignmentProjects is a struct that represents an assignment-projects in the database
type AssignmentProjects struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	TeamID uuid.UUID `gorm:"<-:create;type:uuid;not null" json:"teamId"`
	Team   Team      `json:"team"`

	AssignmentID uuid.UUID  `gorm:"<-:create;not null" json:"-"`
	Assignment   Assignment `json:"assignment"`

	AssignmentAccepted bool `gorm:"not null" json:"assignmentAccepted"`
	ProjectID          int  `json:"projectId"`
} //@Name AssignmentProjects
