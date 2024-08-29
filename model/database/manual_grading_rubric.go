package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManualGradingRubric struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`

	AssignmentID uuid.UUID  `gorm:"not null" json:"-"`
	Assignment   Assignment `json:"-"`

	MaxScore int `gorm:"not null" json:"maxScore"`

	Results []*ManualGradingResult `gorm:"foreignKey:RubricID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
} //@Name ManualGradingRubric

func (m *ManualGradingRubric) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("rubric_id = ?", m.ID).Delete(&ManualGradingResult{})
	return
}
