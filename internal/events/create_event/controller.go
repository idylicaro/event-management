package create_event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	models "github.com/idylicaro/event-management/internal/models"
)

// @Summary Create a new event
// @Description Create a new event
// @Tags events
// @Accept json
// @Produce json
// @Success 201 {object} models.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
func CreateEventController(svc CreateEventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event

		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := svc.CreateEvent(&event); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, event)
	}
}
