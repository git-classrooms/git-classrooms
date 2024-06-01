package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ManualGradingRubric struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`

	AssignmentID uuid.UUID  `gorm:"not null" json:"assignmentId"`
	Assignment   Assignment `json:"-"`

	MaxScore int `gorm:"not null" json:"maxScore"`

	Results []*ManualGradingResult `gorm:"foreignKey:RubricID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
} //@Name ManualGradingRubric
