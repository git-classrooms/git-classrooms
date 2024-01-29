package auth

import (
	"backend/config"
	"fmt"

	"golang.org/x/oauth2"
)

func ConfigGitlab(applicationConfig *config.Config) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     applicationConfig.Auth.ClientID,
		ClientSecret: applicationConfig.Auth.ClientSecret,
		RedirectURL:  applicationConfig.Auth.RedirectURL.String(),
		Scopes:       []string{"api"}, // you can use other scopes to get more data
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", applicationConfig.GitLab.URL), //"https://gitlab.com/oauth/authorize"
			TokenURL: fmt.Sprintf("%s/oauth/token", applicationConfig.GitLab.URL),     //"https://gitlab.com/oauth/token",
		},
	}
	return conf
}
