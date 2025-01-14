package jwt

import "github.com/idylicaro/event-management/internal/models"

type JWTService interface {
	GenerateAccessToken(user models.User) (string, error)
	GenerateRefreshToken(user models.User) (string, error)
}
