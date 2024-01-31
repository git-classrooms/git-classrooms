package mail

import (
	"bytes"
	"crypto/tls"
	mailConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/mail"
	"html/template"

	gomail "gopkg.in/gomail.v2"
)

type MailRepository struct {
	template *template.Template
	dialer   *gomail.Dialer
}

func NewMailRepository(config mailConfig.Config) (*MailRepository, error) {
	t, err := template.ParseFiles(config.GetTemplateFilePath())
	if err != nil {
		return nil, err
	}
	dailer := gomail.NewDialer(config.GetHost(), config.GetPort(), config.GetUser(), config.GetPassword())
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &MailRepository{template: t, dialer: dailer}, nil
}

func (m *MailRepository) Send(to, subject string, mailData MailData) error {
	var tpl bytes.Buffer
	if err := m.template.Execute(&tpl, mailData); err != nil {
		return err
	}

	result := tpl.String()

	mail := gomail.NewMessage()
	mail.SetHeader("From", m.dialer.Username)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", result)

	return m.dialer.DialAndSend(mail)
}
