package list_events

import (
	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/helpers"
	"github.com/idylicaro/event-management/internal/models"
)

type GetEventsController interface {
	Handle(ctx *gin.Context)
}

type GetEventsService interface {
	Execute(filters dto.EventFilters) ([]models.Event, helpers.PaginationMeta, error)
}

type GetEventsRepository interface {
	Execute(filters dto.EventFilters) ([]models.Event, helpers.PaginationMeta, error)
}
