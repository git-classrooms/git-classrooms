package auth

import (
	"log/slog"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type OAuthConfig struct {
	ClientID         string   `env:"CLIENT_ID"`
	ClientSecret     string   `env:"CLIENT_SECRET"`
	RedirectEndpoint string   `env:"REDIRECT_ENDPOINT" envDefault:"/api/v1/auth/gitlab/callback"`
	AuthURL          *url.URL `env:"AUTH_URL,expand" envDefault:"$GITLAB_URL/oauth/authorize"`
	TokenURL         *url.URL `env:"TOKEN_URL,expand" envDefault:"$GITLAB_URL/oauth/token"`
	Scopes           []string `env:"SCOPES" envSeparator:"," envDefault:"api"`
}

var _ slog.LogValuer = (*OAuthConfig)(nil)

func (c *OAuthConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("clientID", strings.Repeat("*", 6)),
		slog.String("clientSecret", strings.Repeat("*", 6)),
		slog.String("redirectEndpoint", c.RedirectEndpoint),
		slog.String("authURL", c.AuthURL.String()),
		slog.String("tokenURL", c.TokenURL.String()),
		slog.String("scopes", strings.Join(c.Scopes, ", ")),
	)
}

func (c *OAuthConfig) GetRedirectEndpoint() string {
	return c.RedirectEndpoint
}

func (c *OAuthConfig) GetOAuthConfig(origin string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  origin + c.RedirectEndpoint,
		Scopes:       c.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthURL.String(),
			TokenURL: c.TokenURL.String(),
		},
	}
}
