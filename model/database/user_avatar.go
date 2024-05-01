package database

import (
	"time"

	"gorm.io/gorm"
)

type UserAvatar struct {
	UserID            int            `gorm:"primary_key" json:"-"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	AvatarURL         *string        `json:"avatarURL"`
	FallbackAvatarURL *string        `json:"fallbackAvatarURL"`
} //@Name UserAvatar
