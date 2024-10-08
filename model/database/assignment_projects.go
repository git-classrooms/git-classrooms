package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type status string

const (
	Pending  status = "pending"
	Creating status = "creating"
	Accepted status = "accepted"
	Failed   status = "failed"
)

// AssignmentProjects is a struct that represents an assignment-projects in the database
type AssignmentProjects struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	TeamID uuid.UUID `gorm:"<-:create;type:uuid;not null" json:"teamId"`
	Team   Team      `json:"team"`

	AssignmentID uuid.UUID  `gorm:"<-:create;not null" json:"-"`
	Assignment   Assignment `json:"assignment"`

	ProjectStatus status `gorm:"not null;default:pending" json:"projectStatus"`
	ProjectID     int    `json:"projectId"`

	GradingJUnitTestResult *JUnitTestResult       `gorm:"type:jsonb;" json:"gradingJUnitTestResult" validate:"optional"`
	GradingManualResults   []*ManualGradingResult `gorm:"foreignKey:AssignmentProjectID;constraint:OnDelete:CASCADE;" json:"gradingManualResults"`
} //@Name AssignmentProjects

type JUnitTestResult struct {
	model.TestReport
}

func (a JUnitTestResult) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JUnitTestResult) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
