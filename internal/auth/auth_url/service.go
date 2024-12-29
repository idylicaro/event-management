package auth_url

import (
	"errors"

	"github.com/idylicaro/event-management/internal/auth/providers"
)

type Service struct {
	Providers map[string]providers.OAuthProvider // Ex: "google": GoogleProvider
}

func (s *Service) GenerateAuthURL(providerName string) (string, error) {
	provider, exists := s.Providers[providerName]
	if !exists {
		return "", errors.New("provider not supported")
	}
	return provider.GetAuthURL(), nil
}
