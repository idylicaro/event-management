// internal/auth/providers/google_provider.go
package providers

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *GoogleProvider) GetAuthURL() string {
	// TODO: Optionally, the front-end can handle generating the code_verifier and code_challenge for PKCE security.
	// The front-end should generate a secure random code_verifier and calculate the code_challenge with SHA256.
	return g.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (g *GoogleProvider) ExchangeCode(ctx context.Context, code string) (TokenResponse, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Sub(token.Expiry).Seconds()),
	}, nil
}

func (g *GoogleProvider) GetUserInfo(accessToken string) (UserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return UserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil
}
