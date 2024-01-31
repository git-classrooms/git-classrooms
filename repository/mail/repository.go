package mail

import "net/url"

type MailData struct {
	name       string
	email      string
	inviteLink *url.URL
}

type Repository interface {
	Send(to, subject string, mailData MailData) error
}
