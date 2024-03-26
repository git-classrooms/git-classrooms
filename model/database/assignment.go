// Package database contains reference types for representing data with gorm
package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Assignment is a struct that represents an assignment in the database
type Assignment struct {
	ID                uuid.UUID             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt        `gorm:"index" json:"-"`
	ClassroomID       uuid.UUID             `gorm:"not null" json:"classroomId"`
	Classroom         Classroom             `json:"-"`
	TemplateProjectID int                   `gorm:"<-:create;not null" json:"templateProjectId"`
	Name              string                `gorm:"not null" json:"name"`
	Description       string                `json:"description"`
	DueDate           *time.Time            `json:"dueDate"`
	Projects          []*AssignmentProjects `json:"-"`
}

// AssignmentProjects is a struct that represents an assignment-projects in the database
type AssignmentProjects struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	AssignmentID       uuid.UUID      `gorm:"<-:create;not null" json:"assignmentId"`
	Assignment         Assignment     `json:"-"`
	UserID             int            `gorm:"<-:create;not null" json:"userId"`
	User               User           `json:"user"`
	AssignmentAccepted bool           `gorm:"not null" json:"assignmentAccepted"`
	ProjectID          int            `json:"projectId"`
}
