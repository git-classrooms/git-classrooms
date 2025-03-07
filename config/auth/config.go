package auth

import (
	"golang.org/x/oauth2"
)

type Config interface {
	GetOAuthConfig(origin string) *oauth2.Config
	GetRedirectEndpoint() string
}
