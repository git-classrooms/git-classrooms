// Package database contains reference types for representing data with gorm
package database

import (
	"time"
)

// User is the representation of the user in database
type User struct {
	ID             int    `gorm:"primary_key;autoIncrement:false" json:"id"`
	GitlabUsername string `gorm:"unique;not null" json:"gitlabUsername"`
	GitlabEmail    string `gorm:"unique;not null" json:"gitlabEmail"`

	AvatarURL         *string `json:"avatarURL"`
	FallbackAvatarURL *string `json:"fallbackAvatarURL"`

	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	OwnedClassrooms []*Classroom      `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE;" json:"-"`
	Classrooms      []*UserClassrooms `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
} //@Name User
