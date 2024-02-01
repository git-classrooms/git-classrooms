package mail

import "time"

type ClassroomInvitationData struct {
	ClassroomName      string
	ClassroomOwnerName string
	RecipientEmail     string
	InvitationPath     string
	ExpireDate         time.Time
}

type Repository interface {
	SendClassroomInvitation(to string, subject string, data ClassroomInvitationData) error
}
