package model

import "time"

type AccessLevelValue int

const (
	NoPermissions            AccessLevelValue = 0
	MinimalAccessPermissions AccessLevelValue = 5
	GuestPermissions         AccessLevelValue = 10
	ReporterPermissions      AccessLevelValue = 20
	DeveloperPermissions     AccessLevelValue = 30
	MaintainerPermissions    AccessLevelValue = 40
	OwnerPermissions         AccessLevelValue = 50
	AdminPermissions         AccessLevelValue = 60
)

type PendingInvite struct {
	ID            int
	InviteEmail   string
	CreatedAt     *time.Time
	AccessLevel   AccessLevelValue
	ExpiresAt     *time.Time
	UserName      string
	CreatedByName string
}
