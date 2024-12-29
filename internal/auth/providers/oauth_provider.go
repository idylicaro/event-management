package providers

import "context"

// OAuthProvider defines the behavior of an OAuth2 provider.
type OAuthProvider interface {
	GetAuthURL() string
	ExchangeCode(ctx context.Context, code string) (TokenResponse, error)
	GetUserInfo(accessToken string) (UserInfo, error)
}

// TokenResponse represents the response after exchanging the authorization code.
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

// UserInfo represents the user information retrieved from the provider.
type UserInfo struct {
	Email      string
	Name       string
	PictureURL string
}
