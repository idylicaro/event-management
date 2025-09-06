package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/idylicaro/event-management/internal/helpers/response"
	"github.com/idylicaro/event-management/internal/models"
)

type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuthMiddleware(jwtSecret []byte) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

// RequireAuth middleware that requires authentication
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := a.extractUserFromToken(c)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "auth.required", "Authentication required")
			c.Abort()
			return
		}

		// Add user to context
		c.Set("user", user)
		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (a *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := a.extractUserFromToken(c)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "auth.required", "Authentication required")
			c.Abort()
			return
		}

		roleAllowed := false
		for _, role := range allowedRoles {
			if user.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			response.Error(c, http.StatusForbidden, "auth.insufficient_permissions", "Insufficient permissions")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuth middleware that extracts user if authenticated, but doesn't require it
func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := a.extractUserFromToken(c)
		if user != nil {
			c.Set("user", user)
		}
		c.Next()
	}
}

func (a *AuthMiddleware) extractUserFromToken(c *gin.Context) (*models.User, error) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, jwt.ErrInvalidKey
	}

	// Check "Bearer <token>" format
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, jwt.ErrInvalidKey
	}

	tokenString := tokenParts[1]

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	// Extract user_id
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// Here you should fetch the complete user from database
	// For now, we return a basic user
	user := &models.User{
		ID: int64(userID),
	}

	return user, nil
}

// GetUserFromContext helper to extract user from context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	userModel, ok := user.(*models.User)
	return userModel, ok
}
