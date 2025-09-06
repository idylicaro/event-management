package callback

import (
	"context"
	"errors"
	"net"

	"github.com/idylicaro/event-management/internal/auth/audit"
	"github.com/idylicaro/event-management/internal/auth/jwt"
	"github.com/idylicaro/event-management/internal/auth/providers"
	"github.com/idylicaro/event-management/internal/auth/security"
	dto "github.com/idylicaro/event-management/internal/dto/auth"
)

type callbackService struct {
	Providers    map[string]providers.OAuthProvider
	Repository   CallbackRepository
	JWTService   jwt.JWTService
	AuditService *audit.AuditService
}

func NewCallbackService(p map[string]providers.OAuthProvider, repo CallbackRepository, jwtS jwt.JWTService, auditS *audit.AuditService) CallbackService {
	return &callbackService{Providers: p, Repository: repo, JWTService: jwtS, AuditService: auditS}
}

func (s *callbackService) Execute(ctx context.Context, providerName, code, state, userAgent, clientIP string) (*dto.TokenResponse, error) {
	provider, exists := s.Providers[providerName]
	if !exists {
		// Log failed attempt
		if s.AuditService != nil {
			s.AuditService.LogLoginAttempt(ctx, audit.LoginAttempt{
				IPAddress:     net.ParseIP(clientIP),
				UserAgent:     userAgent,
				Success:       false,
				FailureReason: "provider not supported",
				Provider:      providerName,
			})
		}
		return nil, errors.New("provider not supported")
	}

	// Validate state for CSRF protection
	err := security.ValidateState(state, providerName, userAgent, clientIP)
	if err != nil {
		// Log failed attempt
		if s.AuditService != nil {
			s.AuditService.LogLoginAttempt(ctx, audit.LoginAttempt{
				IPAddress:     net.ParseIP(clientIP),
				UserAgent:     userAgent,
				Success:       false,
				FailureReason: "invalid state",
				Provider:      providerName,
			})
		}
		return nil, err
	}

	tokens, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		// Log failed attempt
		if s.AuditService != nil {
			s.AuditService.LogLoginAttempt(ctx, audit.LoginAttempt{
				IPAddress:     net.ParseIP(clientIP),
				UserAgent:     userAgent,
				Success:       false,
				FailureReason: "code exchange failed",
				Provider:      providerName,
			})
		}
		return nil, err
	}

	// Exchange the code for a token and fetch user data
	userData, err := provider.GetUserInfo(tokens.AccessToken)
	if err != nil {
		// Log failed attempt
		if s.AuditService != nil {
			s.AuditService.LogLoginAttempt(ctx, audit.LoginAttempt{
				IPAddress:     net.ParseIP(clientIP),
				UserAgent:     userAgent,
				Success:       false,
				FailureReason: "user info fetch failed",
				Provider:      providerName,
			})
		}
		return nil, err
	}

	// Create or find the user in the database
	user, err := s.Repository.FindOrCreateUser(userData)
	if err != nil {
		// Log failed attempt
		if s.AuditService != nil {
			s.AuditService.LogLoginAttempt(ctx, audit.LoginAttempt{
				Email:         userData.Email,
				IPAddress:     net.ParseIP(clientIP),
				UserAgent:     userAgent,
				Success:       false,
				FailureReason: "user creation/lookup failed",
				Provider:      providerName,
			})
		}
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

	// Log successful login
	if s.AuditService != nil {
		s.AuditService.LogLoginAttempt(ctx, audit.LoginAttempt{
			Email:     userData.Email,
			IPAddress: net.ParseIP(clientIP),
			UserAgent: userAgent,
			Success:   true,
			Provider:  providerName,
		})
		s.AuditService.UpdateLastLogin(ctx, user.ID)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
