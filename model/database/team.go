package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name    string `gorm:"not null" json:"name"`
	GroupID int    `gorm:"<-:create;not null" json:"groupId"`

	ClassroomID uuid.UUID `gorm:"<-:create;type:uuid;not null" json:"-"`
	Classroom   Classroom `json:"classroom"`

	Member []*UserClassrooms `gorm:"foreignKey:TeamID" json:"member"`

	AssignmentProjects []*AssignmentProjects `json:"-"`
}
