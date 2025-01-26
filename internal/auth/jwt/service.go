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

// RefreshToken validates the provided refresh token and generates new access and refresh tokens
func (s *jwtService) RefreshToken(refreshToken string) (string, string, error) {
	// Parse the token
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.SecretKey, nil
	})

	// Handle parsing errors or invalid tokens
	if err != nil {
		return "", "", err
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from the claims
		userID, ok := claims["user_id"].(float64) // float64 because of JSON number format
		if !ok {
			return "", "", jwt.ErrInvalidKey
		}

		// Create new tokens
		user := models.User{ID: int64(userID)}
		newAccessToken, err := s.GenerateAccessToken(user)
		if err != nil {
			return "", "", err
		}

		newRefreshToken, err := s.GenerateRefreshToken(user)
		if err != nil {
			return "", "", err
		}

		return newAccessToken, newRefreshToken, nil
	}

	return "", "", jwt.ErrInvalidKey
}
