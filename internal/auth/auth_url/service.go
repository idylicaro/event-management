package auth_url

import (
	"errors"

	"github.com/idylicaro/event-management/internal/auth/providers"
)

type generateAuthURLService struct {
	Providers map[string]providers.OAuthProvider
}

func NewGenerateAuthURLService(providers map[string]providers.OAuthProvider) GenerateAuthURLService {
	return &generateAuthURLService{Providers: providers}
}

func (s *generateAuthURLService) Execute(providerName string) (string, error) {
	provider, exists := s.Providers[providerName]
	if !exists {
		return "", errors.New("provider not supported")
	}
	return provider.GetAuthURL(), nil
}
