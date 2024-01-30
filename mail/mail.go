package mail

import (
	"bytes"
	"crypto/tls"
	mailConfig "de.hs-flensburg.gitlab/gitlab-classroom/config/mail"
	"html/template"

	gomail "gopkg.in/gomail.v2"
)

type Mail struct {
	to      string
	subject string
}

func New(to, subject string) Mail {
	return Mail{to, subject}
}

func (m *Mail) Send(config mailConfig.Config, templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	result := tpl.String()

	mail := gomail.NewMessage()
	mail.SetHeader("From", config.GetUser())
	mail.SetHeader("To", m.to)
	mail.SetHeader("Subject", m.subject)
	mail.SetBody("text/html", result)

	// This needs to be updated when
	// https://gitlab.hs-flensburg.de/fb3-masterprojekt-gitlab-classroom/gitlab-classroom/-/merge_requests/5
	// is merged
	dailer := gomail.NewDialer(config.GetHost(), config.GetPort(), config.GetUser(), config.GetPassword())
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dailer.DialAndSend(mail)
}
