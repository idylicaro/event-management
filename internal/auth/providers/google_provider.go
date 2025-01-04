// internal/auth/providers/google_provider.go
package providers

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
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
			Scopes:       []string{"openid", "profile", "email"},
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
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return TokenResponse{}, fmt.Errorf("ID Token not found in response")
	}
	return TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      idToken,
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

func (g *GoogleProvider) GetUserInfoFromIDToken(idToken string) (UserInfo, error) {
	var userInfo UserInfo

	// Parse and validate the ID token
	parsedToken, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		return getGooglePublicKey()
	})
	if err != nil {
		return userInfo, err
	}

	// Check if the token is valid and has the expected claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return userInfo, errors.New("invalid claims in ID Token")
	}

	// Extract user info from the claims
	userInfo.Email = claims["email"].(string)
	userInfo.Name = claims["name"].(string)
	userInfo.PictureURL = claims["picture"].(string)

	return userInfo, nil
}

func getGooglePublicKey() (*rsa.PublicKey, error) {
	const certsURL = "https://www.googleapis.com/oauth2/v3/certs"
	resp, err := http.Get(certsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var certs map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return nil, err
	}

	keyID := "some-key-id"
	if keyData, ok := certs[keyID]; ok {
		keyDataMap, ok := keyData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid key data")
		}

		publicKeyData, ok := keyDataMap["x5c"].([]interface{})
		if !ok || len(publicKeyData) == 0 {
			return nil, fmt.Errorf("x5c field missing or invalid")
		}

		certData := publicKeyData[0].(string)
		certBytes, err := json.Marshal(certData)
		if err != nil {
			return nil, err
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(certBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse public key: %w", err)
		}

		return publicKey, nil
	}

	return nil, fmt.Errorf("key not found")
}
