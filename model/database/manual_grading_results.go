package database

import "github.com/google/uuid"

type ManualGradingResult struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`

	RubricID uuid.UUID           `gorm:"type:uuid;not null" json:"-"`
	Rubric   ManualGradingRubric `json:"rubric"`

	ProjectGradingDateID uuid.UUID                    `gorm:"<-:create;type:uuid;not null" json:"-"`
	ProjectGradingDate   AssignmentProjectGradingDate `json:"-"`

	Score    int     `gorm:"not null" json:"score"`
	Feedback *string `json:"feedback" validate:"optional"`
} //@Name ManualGradingResult
