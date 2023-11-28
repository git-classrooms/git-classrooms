package auth

import (
	"backend/config"
	"fmt"

	"golang.org/x/oauth2"
)

func ConfigGitlab() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("OAUTH2_CLIENT_ID"),
		ClientSecret: config.Config("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  config.Config("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"api"}, // you can use other scopes to get more data
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", config.Config("GITLAB_URL")),  //"https://gitlab.com/oauth/authorize"
			TokenURL: fmt.Sprintf("%s/oauth/token", config.Config("GITLAB_URL")), //"https://gitlab.com/oauth/token",
		},
	}
	return conf
}
