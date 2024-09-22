// Package mail provides functionality for sending email notifications such as classroom invitations
// and assignment notifications. It uses templates to generate dynamic emails.
package mail

import "time"

// ClassroomInvitationData holds the information required for sending a classroom invitation email.
type ClassroomInvitationData struct {
	ClassroomName      string
	ClassroomOwnerName string
	RecipientEmail     string
	InvitationPath     string
	ExpireDate         time.Time
}

// AssignmentNotificationData holds the information required for sending an assignment notification email.
type AssignmentNotificationData struct {
	ClassroomName      string
	ClassroomOwnerName string
	RecipientName      string
	AssignmentName     string
	JoinPath           string
}

// Repository is an interface that defines the contract for sending email notifications.
type Repository interface {
	SendClassroomInvitation(to string, subject string, data ClassroomInvitationData) error
	SendAssignmentNotification(to string, subject string, data AssignmentNotificationData) error
}
