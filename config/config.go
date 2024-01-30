package config

import (
	"de.hs-flensburg.gitlab/gitlab-classroom/config/auth"
	"de.hs-flensburg.gitlab/gitlab-classroom/config/database"
	"de.hs-flensburg.gitlab/gitlab-classroom/config/general"
	"de.hs-flensburg.gitlab/gitlab-classroom/config/mail"
)

type Config struct {
	Port         int                  `env:"PORT" envDefault:"3000"`
	FrontendPath string               `env:"FRONTEND_PATH" envDefault:"./public"`
	GitLab       general.GitLabConfig `envPrefix:"GITLAB_"`
	Database     database.PsqlConfig  `envPrefix:"POSTGRES_"`
	Auth         auth.Config          `envPrefix:"AUTH_"`
	Mail         mail.Config          `envPrefix:"SMTP_"`
}
