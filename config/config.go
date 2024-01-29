package config

import (
	"backend/config/auth"
	"backend/config/database"
	"backend/config/general"
	"backend/config/mail"
)

type Config struct {
	GitLab   general.GitLabConfig `envPrefix:"GITLAB_"`
	Database database.PsqlConfig  `envPrefix:"POSTGRES_"`
	Auth     auth.Config          `envPrefix:"AUTH_"`
	Mail     mail.Config          `envPrefix:"SMTP_"`
}
