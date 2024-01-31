package auth

import (
	"golang.org/x/oauth2"
	"net/url"
)

type OAuthConfig struct {
	ClientID     string   `env:"CLIENT_ID"`
	ClientSecret string   `env:"CLIENT_SECRET"`
	RedirectURL  *url.URL `env:"REDIRECT_URL"`
	AuthURL      *url.URL `env:"AUTH_URL"`
	TokenURL     *url.URL `env:"TOKEN_URL"`
	Scopes       []string `env:"SCOPES" envSeparator:"," envDefault:"api"`

	config *oauth2.Config
}

func (c *OAuthConfig) GetRedirectUrl() *url.URL {
	return c.RedirectURL
}

func (c *OAuthConfig) GetOAuthConfig() *oauth2.Config {
	if c.config == nil {
		c.config = &oauth2.Config{
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
	return c.config
}
