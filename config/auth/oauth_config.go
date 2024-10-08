package auth

import (
	"net/url"

	"golang.org/x/oauth2"
)

type OAuthConfig struct {
	ClientID     string   `env:"CLIENT_ID"`
	ClientSecret string   `env:"CLIENT_SECRET"`
	RedirectURL  *url.URL `env:"REDIRECT_URL,expand" envDefault:"$PUBLIC_URL/api/v1/auth/gitlab/callback"`
	AuthURL      *url.URL `env:"AUTH_URL,expand" envDefault:"$GITLAB_URL/oauth/authorize"`
	TokenURL     *url.URL `env:"TOKEN_URL,expand" envDefault:"$GITLAB_URL/oauth/token"`
	Scopes       []string `env:"SCOPES" envSeparator:"," envDefault:"api"`
}

func (c *OAuthConfig) GetRedirectUrl() *url.URL {
	return c.RedirectURL
}

func (c *OAuthConfig) GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL.String(),
		Scopes:       c.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthURL.String(),
			TokenURL: c.TokenURL.String(),
		},
	}
}
