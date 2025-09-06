package events

import (
	"github.com/gin-gonic/gin"

	"github.com/idylicaro/event-management/internal/events/create_event"
	"github.com/idylicaro/event-management/internal/events/list_events"
	"github.com/idylicaro/event-management/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes generates the routes for the events module
func RegisterEventsRoutes(router *gin.RouterGroup, db *pgxpool.Pool, authMiddleware *middleware.AuthMiddleware) {
	createEventRepo := create_event.NewEventRepository(db)
	createEventService := create_event.NewCreateEventService(createEventRepo)
	createEventController := create_event.NewCreateEventController(createEventService)

	// POST requires authentication
	router.POST("/", authMiddleware.RequireAuth(), createEventController.Handle)

	listEventsRepo := list_events.NewGetEventsRepository(db)
	listEventsService := list_events.NewGetEventsService(listEventsRepo)
	listEventsController := list_events.NewGetEventsController(listEventsService)

	// GET can be optional auth (for public event listing)
	router.GET("/", authMiddleware.OptionalAuth(), listEventsController.Handle)
}
