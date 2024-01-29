package auth

import "net/url"

type Config struct {
	ClientID     string   `env:"CLIENT_ID"`
	ClientSecret string   `env:"CLIENT_SECRET"`
	RedirectURL  *url.URL `env:"REDIRECT_URL"`
}
