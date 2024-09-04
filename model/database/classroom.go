// Package database contains reference types for representing data with gorm
package database

import (
	"time"

	"github.com/google/uuid"
)

// Classroom is a struct that represents a classroom in the database
type Classroom struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	OwnerID     int    `gorm:"not null" json:"ownerId"`
	Owner       User   `gorm:";" json:"owner"`

	CreateTeams bool `gorm:"not null" json:"createTeams"`
	MaxTeamSize int  `gorm:"not null;default:1" json:"maxTeamSize"`
	MaxTeams    int  `gorm:"not null;default:0" json:"maxTeams"`

	GroupID            int    `gorm:"<-:create;not null" json:"groupId"`
	GroupAccessTokenID int    `gorm:"not null" json:"-"`
	GroupAccessToken   string `gorm:"not null" json:"-"`

	Member                  []*UserClassrooms      `gorm:"foreignKey:ClassroomID;constraint:OnDelete:CASCADE;" json:"-"`
	Teams                   []*Team                `gorm:"foreignKey:ClassroomID;constraint:OnDelete:CASCADE;" json:"-"`
	Assignments             []*Assignment          `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Invitations             []*ClassroomInvitation `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ManualGradingRubrics    []*ManualGradingRubric `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	StudentsViewAllProjects bool                   `gorm:"not null" json:"studentsViewAllProjects"`

	Archived           bool `gorm:"not null;default:false" json:"archived"`
	PotentiallyDeleted bool `gorm:"not null;default:false" json:"potentiallyDeleted"`
} //@Name Classroom
