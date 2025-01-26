package refresh

import "github.com/idylicaro/event-management/internal/auth/jwt"

type refreshTokenService struct {
	JWTService jwt.JWTService
}

func NewRefreshTokenService(jwtService jwt.JWTService) RefreshTokenService {
	return &refreshTokenService{JWTService: jwtService}
}

func (s *refreshTokenService) Execute(refreshToken string) (string, string, error) {
	return s.JWTService.RefreshToken(refreshToken)
}
