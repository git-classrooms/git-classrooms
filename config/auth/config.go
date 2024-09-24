package auth

import (
	"net/url"

	"golang.org/x/oauth2"
)

type Config interface {
	GetOAuthConfig() *oauth2.Config
	GetRedirectURL() *url.URL
}
