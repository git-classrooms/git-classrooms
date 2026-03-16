package config

import (
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/logging"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/mail"
)

type ApplicationConfig struct {
	PublicURL      *url.URL             `env:"PUBLIC_URL" envDefault:"https://staging.hs-flensburg.dev"`
	Port           int                  `env:"PORT" envDefault:"3000"`
	FrontendPath   string               `env:"FRONTEND_PATH" envDefault:"./public"`
	TrustedProxies []string             `env:"TRUSTED_PROXIES" envSeparator:"," envDefault:""`
	GitLab         *gitlab.GitlabConfig `envPrefix:"GITLAB_"`
	Database       *database.PsqlConfig `envPrefix:"POSTGRES_"`
	Auth           *auth.OAuthConfig    `envPrefix:"AUTH_"`
	Mail           *mail.MailConfig     `envPrefix:"SMTP_"`
	Log            *logging.LogConfig   `envPrefix:"LOG_"`
}

func LoadApplicationConfig() (*ApplicationConfig, error) {
	path, _ := os.Getwd()

	godotenv.Load(filepath.Join(path, ".env"), filepath.Join(path, ".env.local"))

	config := &ApplicationConfig{
		GitLab:   &gitlab.GitlabConfig{},
		Database: &database.PsqlConfig{},
		Auth:     &auth.OAuthConfig{},
		Mail:     &mail.MailConfig{},
		Log:      &logging.LogConfig{},
	}
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

var _ slog.LogValuer = (*ApplicationConfig)(nil)

func (c *ApplicationConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("publicURL", c.PublicURL.String()),
		slog.Int("port", c.Port),
		slog.String("frontendPath", c.FrontendPath),
		slog.String("trustedProxies", strings.Join(c.TrustedProxies, ", ")),
		slog.Any("gitlab", c.GitLab),
		slog.Any("database", c.Database),
		slog.Any("auth", c.Auth),
		slog.Any("mail", c.Mail),
		slog.Any("log", c.Log),
	)
}
