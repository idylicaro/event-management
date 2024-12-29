package create_event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

// @Summary Create a new event
// @Description Create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Success 201 {object} models.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
func CreateEventController(svc CreateEventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event dto.CreateEventRequest

		if err := c.ShouldBindJSON(&event); err != nil {
			response.Error(c, http.StatusBadRequest, "validation.body.failed", err.Error())
			return
		}

		if err := svc.Execute(&event); err != nil {
			response.Error(c, http.StatusInternalServerError, "create.event.fail", err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "create.event.success", event, nil)
	}
}
