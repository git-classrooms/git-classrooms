package model

import "time"

type GroupAccessToken struct {
	ID          int
	UserID      int
	Name        string
	Scopes      []string
	ExpiresAt   time.Time
	Token       string
	AccessLevel AccessLevelValue
}
