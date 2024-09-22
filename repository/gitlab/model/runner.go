package model

import "time"

type Runner struct {
	ID             int        `json:"id"`
	Description    string     `json:"description"`
	Active         bool       `json:"active"`
	Paused         bool       `json:"paused"`
	IsShared       bool       `json:"is_shared"`
	IPAddress      string     `json:"ip_address"`
	RunnerType     string     `json:"runner_type"`
	Name           string     `json:"name"`
	Online         bool       `json:"online"`
	Status         string     `json:"status"`
	Token          string     `json:"token"`
	TokenExpiresAt *time.Time `json:"token_expires_at"`
} // @Name Runner
