// Package database contains reference types for representing data with gorm
package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Classroom is a struct that represents a classroom in the database
type Classroom struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	OwnerID     int    `gorm:"not null" json:"ownerId"`
	Owner       User   `json:"owner"`

	CreateTeams bool `gorm:"not null" json:"createTeams"`
	MaxTeamSize int  `gorm:"not null;default:1" json:"maxTeamSize"`
	MaxTeams    int  `gorm:"not null;default:0" json:"maxTeams"`

	GroupID            int    `gorm:"<-:create;not null" json:"groupId"`
	GroupAccessTokenID int    `gorm:"not null" json:"-"`
	GroupAccessToken   string `gorm:"not null" json:"-"`

	Member                  []*UserClassrooms      `gorm:"foreignKey:ClassroomID" json:"-"`
	Teams                   []*Team                `gorm:"foreignKey:ClassroomID" json:"-"`
	Assignments             []*Assignment          `json:"-"`
	Invitations             []*ClassroomInvitation `json:"-"`
	StudentsViewAllProjects bool                   `gorm:"not null" json:"studentsViewAllProjects"`
} //@Name Classroom
