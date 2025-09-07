package update_event

import (
	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/models"
)

type UpdateEventController interface {
	Handle(ctx *gin.Context)
}

type UpdateEventService interface {
	Execute(ctx *gin.Context, eventID int64, req *dto.UpdateEventRequest) (*models.Event, error)
}

type UpdateEventRepository interface {
	Execute(eventID, userID int64, req *dto.UpdateEventRequest) (*models.Event, error)
	GetByIDAndUserID(eventID, userID int64) (*models.Event, error)
}
