package database

import (
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	ID             int `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	ClassroomID    int
	TemplateRepoID string `gorm:"<-:create"`
	Repositories   []AssignmentRepository
}

type AssignmentRepository struct {
	ID           int `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	AssignmentID int
	UserID       int
}
