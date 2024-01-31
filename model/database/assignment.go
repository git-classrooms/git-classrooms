// Package database contains reference types for representing data with gorm
package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Assignment is a struct that represents an assignment in the database
type Assignment struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	ClassroomID       uuid.UUID      `gorm:"not null"`
	Classroom         Classroom
	TemplateProjectID int `gorm:"<-:create;not null"`
	Projects          []AssignmentProjects
}

// AssignmentDTO is the data transfer object representing a user
type AssignmentDTO struct {
	ID                uuid.UUID `json:"id"`
	ClassroomID       uuid.UUID `json:"classroomId"`
	TemplateProjectID int       `json:"templateProjectId"`
}

// AssignmentProjects is a struct that represents an assignment-projects in the database
type AssignmentProjects struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	AssignmentID uuid.UUID      `gorm:"<-:create;not null"`
	Assignment   Assignment
	UserID       uuid.UUID `gorm:"<-:create;not null"`
	User         User
	ProjectID    int `gorm:"<-:create;not null"`
}

// AssignmentRepositoryDTO is the data transfer object representing a user
type AssignmentRepositoryDTO struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	AssignmentID uuid.UUID `json:"assignmentId"`
	ProjectID    int       `json:"projectId"`
}
