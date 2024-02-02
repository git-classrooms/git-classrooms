package mail

import "time"

type ClassroomInvitationData struct {
	ClassroomName      string
	ClassroomOwnerName string
	RecipientEmail     string
	InvitationPath     string
	ExpireDate         time.Time
}

type AssignmentNotificationData struct {
	ClassroomName      string
	ClassroomOwnerName string
	RecipientName      string
	AssignmentName     string
	JoinPath           string
}

type Repository interface {
	SendClassroomInvitation(to string, subject string, data ClassroomInvitationData) error
	SendAssignmentNotification(to string, subject string, data AssignmentNotificationData) error
}
