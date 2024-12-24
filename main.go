package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/config"
	"github.com/idylicaro/event-management/internal/events"
)

func main() {
	cfg := config.LoadConfig()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	connPool, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer connPool.Close()

	api := r.Group("/api/v1")

	eventsGroup := api.Group("/events")
	events.RegisterEventsRoutes(eventsGroup, connPool)

	r.Run(fmt.Sprintf(":%s", cfg.ServerPort))
}
