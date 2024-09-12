package database

import (
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
