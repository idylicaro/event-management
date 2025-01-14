package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/idylicaro/event-management/internal/models"
)

// Constants for token durations
const (
	AccessTokenDuration  = time.Hour * 1
	RefreshTokenDuration = time.Hour * 24 * 7
)

type jwtService struct {
	SecretKey []byte
}

// NewService creates a new instance of the JWT Service
func NewJWTService(secretKey []byte) JWTService {
	return &jwtService{SecretKey: secretKey}
}

// GenerateAccessToken generates an access token for the user
func (s *jwtService) GenerateAccessToken(user models.User) (string, error) {
	return s.generateToken(user, AccessTokenDuration)
}

// GenerateRefreshToken generates a refresh token for the user
func (s *jwtService) GenerateRefreshToken(user models.User) (string, error) {
	return s.generateToken(user, RefreshTokenDuration)
}

// generateToken is a helper function to generate JWT tokens
func (s *jwtService) generateToken(user models.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
