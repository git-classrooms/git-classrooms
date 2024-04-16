// Package database contains reference types for representing data with gorm
package database

import (
	"time"

	"gorm.io/gorm"
)

// User is the representation of the user in database
type User struct {
	ID              int               `gorm:"primary_key;autoIncrement:false" json:"id"`
	GitlabEmail     string            `gorm:"unique;not null" json:"gitlabEmail"`
	GitLabAvatar    UserAvatar        `json:"gitlabAvatar"`
	Name            string            `gorm:"not null" json:"name"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
	OwnedClassrooms []*Classroom      `gorm:"foreignKey:OwnerID" json:"-"`
	Classrooms      []*UserClassrooms `gorm:"foreignKey:UserID" json:"-"`
}

type UserAvatar struct {
	UserID            int            `gorm:"primary_key" json:"userId"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	AvatarURL         *string        `json:"avatarURL"`
	FallbackAvatarURL *string        `json:"fallbackAvatarURL"`
}
