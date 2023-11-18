package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	ClassroomID       uuid.UUID
	Classroom         Classroom
	TemplateProjectID int `gorm:"<-:create"`
	Projects          []AssignmentProjects
}

type AssignmentDTO struct {
	ID          uuid.UUID `json:"id"`
	ClassroomID uuid.UUID `json:"classroom_id"`
}

type AssignmentProjects struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	AssignmentID uuid.UUID      `gorm:"<-:create"`
	Assignment   Assignment
	UserID       uuid.UUID `gorm:"<-:create"`
	User         User
	ProjectID    int `gorm:"<-:create"`
}

type AssignmentRepositoryDTO struct {
}
