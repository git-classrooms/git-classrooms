package mail

import (
	"backend/config"
	"bytes"
	"crypto/tls"
	"html/template"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

type MailConfig struct {
	host     string
	port     int
	user     string
	password string
}

type Mail struct {
	to      string
	subject string
}

func New(to, subject string) Mail {
	return Mail{to, subject}
}

func (m *Mail) Send(templateFileName string, data interface{}) error {
	cfg, err := mailConfig()
	if err != nil {
		return err
	}

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
	mail.SetHeader("From", cfg.user)
	mail.SetHeader("To", m.to)
	mail.SetHeader("Subject", m.subject)
	mail.SetBody("text/html", result)

	// This needs to be updated when
	// https://gitlab.hs-flensburg.de/fb3-masterprojekt-gitlab-classroom/gitlab-classroom/-/merge_requests/5
	// is merged
	dailer := gomail.NewDialer(cfg.host, cfg.port, cfg.user, cfg.password)
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dailer.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}

func mailConfig() (*MailConfig, error) {
	p := config.Config("SMTP_PORT")

	port, err := strconv.Atoi(p)

	if err != nil {
		return nil, err
	}

	return &MailConfig{
		host:     config.Config("SMTP_HOST"),
		port:     port,
		user:     config.Config("SMTP_USER"),
		password: config.Config("SMTP_PASSWORD"),
	}, nil

}
