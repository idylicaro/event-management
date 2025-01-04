package callback

import (
	"context"
	"errors"

	"github.com/idylicaro/event-management/internal/auth/providers"
	"github.com/idylicaro/event-management/internal/models"
)

type Service struct {
	Providers  map[string]providers.OAuthProvider
	Repository Repository
}

func (s *Service) ProcessCallback(ctx context.Context, providerName, code string) (models.User, error) {
	provider, exists := s.Providers[providerName]
	if !exists {
		return models.User{}, errors.New("provider not supported")
	}

	tokens, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return models.User{}, err
	}

	// Exchange the code for a token and fetch user data
	userData, err := provider.GetUserInfo(tokens.AccessToken)
	if err != nil {
		return models.User{}, err
	}

	// Create or find the user in the database
	user, err := s.Repository.FindOrCreateUser(userData)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
