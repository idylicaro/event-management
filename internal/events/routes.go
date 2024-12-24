package events

import (
	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/events/create_event"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes generates the routes for the events module
func RegisterEventsRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	router.POST("/", create_event.CreateEventController(create_event.NewCreateEventService(create_event.NewEventRepository(db))))
}
