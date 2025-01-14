package create_event

import (
	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	models "github.com/idylicaro/event-management/internal/models"
)

type CreateEventController interface {
	Handle(ctx *gin.Context)
}

type CreateEventService interface {
	Execute(event *dto.CreateEventRequest) error
}

type CreateEventRepository interface {
	Execute(event *models.Event) error
}
