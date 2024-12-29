package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/config"
	"github.com/idylicaro/event-management/internal/auth/auth_url"
	"github.com/idylicaro/event-management/internal/auth/callback"
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

	authUrlService := auth_url.Service{Providers: providers}
	authUrlController := auth_url.Controller{Service: authUrlService}
	router.GET("/:provider/url", authUrlController.GetAuthURL)

	callbackRepo := callback.NewRepository(db)
	callbackService := callback.Service{Providers: providers, Repository: callbackRepo}
	callbackController := callback.Controller{Service: callbackService}
	router.GET("/:provider/callback", callbackController.HandleCallback)
}
