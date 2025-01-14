package create_event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type createEventController struct {
	Service CreateEventService
}

func NewCreateEventController(service CreateEventService) CreateEventController {
	return &createEventController{Service: service}
}

// @Summary Create a new event
// @Description Create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Success 201 {object} models.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
func (c *createEventController) Handle(ctx *gin.Context) {
	var event dto.CreateEventRequest

	if err := ctx.ShouldBindJSON(&event); err != nil {
		response.Error(ctx, http.StatusBadRequest, "validation.body.failed", err.Error())
		return
	}

	if err := c.Service.Execute(&event); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "create.event.fail", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "create.event.success", event, nil)
}
