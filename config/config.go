package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/mail"
	"net/url"
	"os"
	"path/filepath"
)

type ApplicationConfig struct {
	PublicURL    *url.URL             `env:"PUBLIC_URL" envDefault:"https://staging.hs-flensburg.dev"`
	Port         int                  `env:"PORT" envDefault:"3000"`
	FrontendPath string               `env:"FRONTEND_PATH" envDefault:"./public"`
	GitLab       *gitlab.GitlabConfig `envPrefix:"GITLAB_"`
	Database     *database.PsqlConfig `envPrefix:"POSTGRES_"`
	Auth         *auth.OAuthConfig    `envPrefix:"AUTH_"`
	Mail         *mail.MailConfig     `envPrefix:"SMTP_"`
}

func LoadApplicationConfig() (*ApplicationConfig, error) {
	path, _ := os.Getwd()

	godotenv.Load(filepath.Join(path, ".env"), filepath.Join(path, ".env.local"))

	config := &ApplicationConfig{
		GitLab:   &gitlab.GitlabConfig{},
		Database: &database.PsqlConfig{},
		Auth:     &auth.OAuthConfig{},
		Mail:     &mail.MailConfig{},
	}
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
