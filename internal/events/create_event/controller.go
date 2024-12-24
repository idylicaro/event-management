package create_event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	models "github.com/idylicaro/event-management/internal/models"
)

// Controlador para criar um evento
func CreateEventController(svc CreateEventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event

		// Bind do corpo da requisição para o modelo Event
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Chama o serviço para criar o evento
		if err := svc.CreateEvent(&event); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, event)
	}
}
