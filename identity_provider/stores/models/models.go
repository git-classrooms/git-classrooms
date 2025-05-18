package models

import (
	"time"

	"github.com/akatranlp/sentinel/account"
	"github.com/google/uuid"
)

type Session struct {
	Token  string    `gorm:"primary_key"`
	Data   []byte    `gorm:"type:bytea;not null"`
	Expiry time.Time `gorm:"not null;index"`
}

type TokenSession struct {
	SessionID  string    `gorm:"primary_key"`
	RefreshJTI string    `gorm:"not null;index"`
	Expiry     time.Time `gorm:"not null;index"`
}

type AuthUser struct {
	ID            uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name          string        `gorm:"not null"`
	Username      string        `gorm:"not null"`
	Picture       string        `gorm:"not null"`
	Email         string        `gorm:"not null"`
	EmailVerified bool          `gorm:"not null;default:false"`
	Accounts      []AuthAccount `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (a AuthUser) MapToUser() account.User {
	return account.User{
		UserID:        account.UserID(a.ID.String()),
		Name:          a.Name,
		Username:      a.Username,
		Picture:       a.Picture,
		Email:         a.Email,
		EmailVerified: a.EmailVerified,
	}
}

type AuthAccount struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Provider          string    `gorm:"not null;uniqueIndex:idx_unique_provider_provider_id;uniqueIndex:idx_unique_provider_user_id"`
	ProviderID        string    `gorm:"not null;uniqueIndex:idx_unique_provider_provider_id"`
	Email             string    `gorm:"not null"`
	EmailVerified     bool      `gorm:"not null;default:false"`
	AccessToken       string    `gorm:"not null"`
	Expiry            time.Time `gorm:"not null"`
	RefreshToken      string    `gorm:"not null"`
	RefreshExpiry     time.Time `gorm:"not null"`
	TokenType         string    `gorm:"not null"`
	IDToken           string    `gorm:"not null"`
	Name              string    `gorm:"not null"`
	PreferredUsername string    `gorm:"not null"`
	Nickname          string    `gorm:"not null"`
	Picture           string    `gorm:"not null"`
	Profile           string    `gorm:"not null"`
	UserID            uuid.UUID `gorm:"<-:create;not null;uniqueIndex:idx_unique_provider_user_id"`
	User              AuthUser  `gorm:"not null"`
}

func (a AuthAccount) MapToAccount() account.Account {
	return account.Account{
		AccountID: account.AccountID{
			Provider:   a.Provider,
			ProviderID: a.ProviderID,
		},
		Email:             a.Email,
		EmailVerified:     a.EmailVerified,
		AccessToken:       a.AccessToken,
		Expiry:            a.Expiry,
		RefreshToken:      a.RefreshToken,
		RefreshExpiry:     a.RefreshExpiry,
		TokenType:         a.TokenType,
		IDToken:           a.IDToken,
		Name:              a.Name,
		PreferredUsername: a.PreferredUsername,
		Nickname:          a.Nickname,
		Picture:           a.Picture,
		Profile:           a.Profile,
		UserID:            account.UserID(a.UserID.String()),
	}
}
