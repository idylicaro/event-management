package delete_event

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type deleteEventController struct {
	Service DeleteEventService
}

func NewDeleteEventController(service DeleteEventService) DeleteEventController {
	return &deleteEventController{Service: service}
}

// @Summary Delete an event
// @Description Delete an event that belongs to the authenticated user
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /events/{id} [delete]
func (c *deleteEventController) Handle(ctx *gin.Context) {
	// Extract and validate event ID from path
	idParam := ctx.Param("id")
	if idParam == "" {
		response.Error(ctx, http.StatusBadRequest, "validation.id.required", "Event ID is required")
		return
	}

	eventID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "validation.id.invalid", "Invalid event ID format")
		return
	}

	// Execute delete service
	if err := c.Service.Execute(ctx, eventID); err != nil {
		if err.Error() == "event not found" {
			response.Error(ctx, http.StatusNotFound, "delete.event.not_found", "Event not found")
			return
		}
		if err.Error() == "unauthorized" {
			response.Error(ctx, http.StatusForbidden, "delete.event.unauthorized", "You are not authorized to delete this event")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "delete.event.fail", err.Error())
		return
	}

	// Return 204 No Content on successful deletion
	ctx.Status(http.StatusNoContent)
}
