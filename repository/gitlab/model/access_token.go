package model

import "time"

type GroupAccessToken struct {
	ID          int
	UserID      int
	Name        string
	Scopes      []string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Token       string
	AccessLevel AccessLevelValue
}
