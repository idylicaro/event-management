package delete_event

import (
	"github.com/gin-gonic/gin"
)

type DeleteEventController interface {
	Handle(ctx *gin.Context)
}

type DeleteEventService interface {
	Execute(ctx *gin.Context, eventID int64) error
}

type DeleteEventRepository interface {
	Execute(eventID, userID int64) error
	GetByIDAndUserID(eventID, userID int64) error
}
