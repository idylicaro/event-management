package callback

import (
	"context"
	"errors"

	"github.com/idylicaro/event-management/internal/auth/jwt"
	"github.com/idylicaro/event-management/internal/auth/providers"
	dto "github.com/idylicaro/event-management/internal/dto/auth"
)

type callbackService struct {
	Providers  map[string]providers.OAuthProvider
	Repository CallbackRepository
	JWTService jwt.JWTService
}

func NewCallbackService(p map[string]providers.OAuthProvider, repo CallbackRepository, jwtS jwt.JWTService) CallbackService {
	return &callbackService{Providers: p, Repository: repo, JWTService: jwtS}
}

func (s *callbackService) Execute(ctx context.Context, providerName, code string) (*dto.TokenResponse, error) {
	provider, exists := s.Providers[providerName]
	if !exists {
		return nil, errors.New("provider not supported")
	}

	tokens, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Exchange the code for a token and fetch user data
	userData, err := provider.GetUserInfo(tokens.AccessToken)
	if err != nil {
		return nil, err
	}

	// Create or find the user in the database
	user, err := s.Repository.FindOrCreateUser(userData)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.JWTService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
