// Package database contains reference types for representing data with gorm
package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role uint8

const (
	Owner Role = iota
	Moderator
	Student
)

// Classroom is a struct that represents a classroom in the database
type Classroom struct {
	ID                 uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt          time.Time              `json:"-"`
	UpdatedAt          time.Time              `json:"-"`
	DeletedAt          gorm.DeletedAt         `gorm:"index" json:"-"`
	Name               string                 `gorm:"not null" json:"name"`
	OwnerID            int                    `gorm:"not null" json:"ownerId"`
	Owner              User                   `json:"owner"`
	Description        string                 `json:"description"`
	GroupID            int                    `gorm:"<-:create;not null" json:"groupId"`
	GroupAccessTokenID int                    `gorm:"not null" json:"-"`
	GroupAccessToken   string                 `gorm:"not null" json:"-"`
	Member             []*UserClassrooms      `gorm:"foreignKey:ClassroomID" json:"-"`
	Assignments        []*Assignment          `json:"-"`
	Invitations        []*ClassroomInvitation `json:"-"`
}

// UserClassrooms is a struct that represents the relationship between a user and a classroom
type UserClassrooms struct {
	UserID      int       `gorm:"primaryKey;autoIncrement:false;not null" json:"-"`
	User        User      `json:"-"`
	ClassroomID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"-"`
	Classroom   Classroom `json:"classroom"`
	Role        Role      `gorm:"not null" json:"role"`
}
