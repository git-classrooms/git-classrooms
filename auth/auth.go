package auth

import (
	"backend/config"
	"fmt"

	"golang.org/x/oauth2"
)

func ConfigGitlab() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.GetConfig().Auth.ClientID,
		ClientSecret: config.GetConfig().Auth.ClientSecret,
		RedirectURL:  config.GetConfig().Auth.RedirectURL.String(),
		Scopes:       []string{"api"}, // you can use other scopes to get more data
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", config.GetConfig().GitLab.URL), //"https://gitlab.com/oauth/authorize"
			TokenURL: fmt.Sprintf("%s/oauth/token", config.GetConfig().GitLab.URL),     //"https://gitlab.com/oauth/token",
		},
	}
	return conf
}
