package refresh

import (
	"context"
	"net"

	"github.com/idylicaro/event-management/internal/auth/audit"
	"github.com/idylicaro/event-management/internal/auth/jwt"
)

type refreshTokenService struct {
	JWTService   jwt.JWTService
	AuditService *audit.AuditService
}

func NewRefreshTokenService(jwtService jwt.JWTService, auditService *audit.AuditService) RefreshTokenService {
	return &refreshTokenService{JWTService: jwtService, AuditService: auditService}
}

func (s *refreshTokenService) Execute(refreshToken, userAgent, clientIP string) (string, string, error) {
	newAccessToken, newRefreshToken, err := s.JWTService.RefreshToken(refreshToken)

	// Log refresh attempt if audit service is available
	if s.AuditService != nil {
		success := err == nil
		failureReason := ""
		if err != nil {
			failureReason = err.Error()
		}

		s.AuditService.LogLoginAttempt(context.Background(), audit.LoginAttempt{
			IPAddress:     net.ParseIP(clientIP),
			UserAgent:     userAgent,
			Success:       success,
			FailureReason: failureReason,
			Provider:      "refresh_token",
		})
	}

	return newAccessToken, newRefreshToken, err
}
