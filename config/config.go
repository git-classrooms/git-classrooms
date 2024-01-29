package config

import (
	"backend/config/auth"
	"backend/config/database"
	"backend/config/general"
	"backend/config/mail"
)

type Config struct {
	Port         int                  `env:"PORT" envDefault:"3000"`
	FrontendPath string               `env:"FRONTEND_PATH" envDefault:"./public"`
	GitLab       general.GitLabConfig `envPrefix:"GITLAB_"`
	Database     database.PsqlConfig  `envPrefix:"POSTGRES_"`
	Auth         auth.Config          `envPrefix:"AUTH_"`
	Mail         mail.Config          `envPrefix:"SMTP_"`
}
