package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/config"
	"github.com/idylicaro/event-management/internal/auth/auth_url"
	"github.com/idylicaro/event-management/internal/auth/callback"
	"github.com/idylicaro/event-management/internal/auth/jwt"
	"github.com/idylicaro/event-management/internal/auth/providers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthRoutes(router *gin.RouterGroup, db *pgxpool.Pool, cfg config.Config) {
	providers := map[string]providers.OAuthProvider{
		"google": providers.NewGoogleProvider(
			cfg.GoogleClientID,
			cfg.GoogleClientSecret,
			cfg.GoogleRedirectURL,
		),
	}

	authUrlService := auth_url.NewGenerateAuthURLService(providers)
	authUrlController := auth_url.NewGenerateAuthURLController(authUrlService)

	router.GET("/:provider/url", authUrlController.Handle)

	jwtService := jwt.NewJWTService([]byte(cfg.JWTSecret))
	callbackRepo := callback.NewCallbackRepository(db)
	callbackService := callback.NewCallbackService(providers, callbackRepo, jwtService)
	callbackController := callback.NewCallbackController(callbackService)

	router.GET("/:provider/callback", callbackController.Handle)
}
