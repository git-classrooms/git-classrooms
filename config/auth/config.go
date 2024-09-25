package auth

import (
	"golang.org/x/oauth2"
	"net/url"
)

type Config interface {
	GetOAuthConfig() *oauth2.Config
	GetRedirectUrl() *url.URL
}
