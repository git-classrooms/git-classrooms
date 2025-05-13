package identityprovider

import (
	"context"
	"net/http"
	"os"
	"time"

	// accmemorystore "github.com/akatranlp/identity-provider/account/memory_store"
	"github.com/akatranlp/identity-provider/openid"
	"github.com/akatranlp/identity-provider/openid/enums"
	"github.com/akatranlp/identity-provider/provider"
	// tokenmemorystore "github.com/akatranlp/identity-provider/token/memory_store"
	"github.com/akatranlp/identity-provider/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/identity_provider/stores"
	"gorm.io/gorm"
)

func CreateIdentityProvider(ctx context.Context, db *gorm.DB, config *config.ApplicationConfig) (http.Handler, error) {
	f, err := os.Open("jwk.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	userStore := stores.NewGORMUserStore(db)
	// userStore := accmemorystore.NewMemoryUserStore()
	tokenStore := stores.NewGORMTokenStore(db)
	// tokenStore := tokenmemorystore.NewMemoryTokenStore()

	sessionStore, err := stores.NewSessionStore(db)
	if err != nil {
		return nil, err
	}

	ip, err := openid.NewIdentityProvider(
		"/auth",
		userStore,
		tokenStore,
		sessionStore,
		openid.WithClients(openid.ClientRegistration{
			ClientID:            config.Auth.ClientID,
			ClientSecret:        config.Auth.ClientSecret,
			TokenExchangeSecret: "",
			Scope:               enums.ScopeValues(),
			RedirectURIs: []string{
				"http://localhost" + config.Auth.GetRedirectEndpoint(),
				"https://oidcdebugger.com/debug",
			},
		}),
		openid.WithProviders(utils.Must(provider.ProviderFactory(provider.FactoryParams{
			Type:    "gitlab",
			Name:    "GitLab",
			Slug:    "gitlab",
			BaseURL: "https://gitlab.git-classrooms.dev",
			// TODO: Use a different ClientID and Secret for the ip and not gitlabs
			ClientID:     config.Auth.ClientID,
			ClientSecret: config.Auth.ClientSecret,
		}))),
		openid.WithAccessTokenExpiration(15*time.Minute),
		openid.WithRefreshTokenExpiration(7*24*time.Hour),
		openid.WithSessionUnAuthedLifeTime(10*time.Minute),
		openid.WithSessionAuthedLifeTime(30*24*time.Hour),
		openid.WithSessionIdleTimeout(7*24*time.Hour),
		openid.WithSigningKeyReader(f),
		openid.WithAppURL("../"),
	)
	if err != nil {
		return nil, err
	}

	return ip.Handler(), nil
}
