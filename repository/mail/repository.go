package mail

import "net/url"

type MailData struct {
	Name       string
	Email      string
	InviteLink *url.URL
}

type Repository interface {
	Send(to, subject string, mailData MailData) error
}
