package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/config"
	"github.com/idylicaro/event-management/internal/auth"
	"github.com/idylicaro/event-management/internal/auth/security"
	"github.com/idylicaro/event-management/internal/events"
	"github.com/idylicaro/event-management/internal/helpers/response"
	"github.com/idylicaro/event-management/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/idylicaro/event-management/docs"
)

// @title Event Management API
// @version 1.0
// @description API to manage events

// @contact.name Idyl Santos
// @contact.url https://github.com/idylicaro
// @contact.email suport@example.com

// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.LoadConfig()

	r := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CorsAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(response.ErrorHandlerMiddleware())

	// Global rate limiting
	globalRateLimit := middleware.NewRateLimiter(100, time.Minute) // 100 req/min
	r.Use(globalRateLimit.RateLimit())

	connPool, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer connPool.Close()

	// Authentication middleware
	authMiddleware := middleware.NewAuthMiddleware([]byte(cfg.JWTSecret))

	// Auth-specific rate limiting
	authRateLimit := middleware.NewAuthRateLimiter()

	api := r.Group("/api/v1")

	// Events routes with optional protection
	eventsGroup := api.Group("/events")
	eventsGroup.Use(authMiddleware.OptionalAuth())
	events.RegisterEventsRoutes(eventsGroup, connPool)

	// Protected admin routes
	adminGroup := api.Group("/admin")
	adminGroup.Use(authMiddleware.RequireRole("admin"))
	// Register admin routes here

	// Authentication routes with rate limiting
	authGroup := api.Group("/auth")
	authGroup.Use(authRateLimit.RateLimit())
	auth.RegisterAuthRoutes(authGroup, connPool, *cfg)

	// Swagger (disable in production)
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// Periodic cleanup of expired states and PKCEs
	go startCleanupRoutines()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	r.Run(fmt.Sprintf(":%s", cfg.ServerPort))
}

func startCleanupRoutines() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Cleanup expired states and PKCEs
		security.CleanupExpiredStates()
		security.CleanupExpiredPKCE()
	}
}
