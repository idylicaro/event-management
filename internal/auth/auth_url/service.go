package auth_url

import (
	"errors"

	"github.com/idylicaro/event-management/internal/auth/providers"
	"github.com/idylicaro/event-management/internal/auth/security"
)

type generateAuthURLService struct {
	Providers map[string]providers.OAuthProvider
}

func NewGenerateAuthURLService(providers map[string]providers.OAuthProvider) GenerateAuthURLService {
	return &generateAuthURLService{Providers: providers}
}

func (s *generateAuthURLService) Execute(providerName, userAgent, clientIP string) (string, error) {
	provider, exists := s.Providers[providerName]
	if !exists {
		return "", errors.New("provider not supported")
	}

	// Generate secure state for CSRF protection
	state, err := security.GenerateState(providerName, userAgent, clientIP)
	if err != nil {
		return "", err
	}

	return provider.GetAuthURL(state), nil
}
