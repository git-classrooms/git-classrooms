package identityprovider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/akatranlp/identity-provider/openid"
	"github.com/akatranlp/identity-provider/openid/enums"
	"golang.org/x/oauth2"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func ExchangeAccessToken(ctx context.Context, oauthConfig *oauth2.Config, token string) (*TokenResponse, error) {
	params := make(url.Values)
	params.Set(openid.TokenFormValueGrantType.String(), enums.GrantTypeTokenExchange.String())
	params.Set(openid.TokenFormValueRequestedTokenType.String(), enums.OauthTokenTypeAccessToken.String())
	params.Set(openid.TokenFormValueRequestedIssuer.String(), "gitlab")

	params.Set(openid.TokenFormValueSubjectToken.String(), token)
	params.Set(openid.TokenFormValueSubjectTokenType.String(), enums.OauthTokenTypeAccessToken.String())

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, oauthConfig.Endpoint.TokenURL, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(res.Body)
		return nil, errors.New(string(data))
	}
	var newToken TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&newToken); err != nil {
		return nil, err
	}
	return &newToken, nil
}
