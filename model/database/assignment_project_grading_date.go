package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type AssignmentProjectGradingDate struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	Branchname string `json:"-"`

	AssignmentProject   AssignmentProjects `json:"-"`
	AssignmentProjectID uuid.UUID          `gorm:"type:uuid;uniqueIndex:idx_project_date;not null" json:"assignmentProjectId"`

	AssignmentDate   AssignmentDate `json:"-"`
	AssignmentDateID uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_project_date;not null" json:"assignmentDateId"`

	GradingJUnitTestResult *JUnitTestResult       `gorm:"type:jsonb;" json:"gradingJUnitTestResult" validate:"optional"`
	GradingManualResults   []*ManualGradingResult `gorm:"foreignKey:AssignmentProjectGradingDateID;constraint:OnDelete:CASCADE;" json:"gradingManualResults"`
} //@Name AssignmentProjectGradingDate

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
