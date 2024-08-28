package database

import (
	"time"
)

type UserAvatar struct {
	UserID            int       `gorm:"primary_key" json:"-"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	AvatarURL         *string   `json:"avatarURL"`
	FallbackAvatarURL *string   `json:"fallbackAvatarURL"`
} //@Name UserAvatar
