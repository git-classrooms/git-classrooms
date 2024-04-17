package database

import (
	"gorm.io/gorm"
	"time"
)

type UserAvatar struct {
	UserID            int            `gorm:"primary_key" json:"userId"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	AvatarURL         *string        `json:"avatarURL"`
	FallbackAvatarURL *string        `json:"fallbackAvatarURL"`
}
