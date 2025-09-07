package update_event

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/helpers/response"
	"github.com/idylicaro/event-management/internal/mappers"
)

type updateEventController struct {
	Service UpdateEventService
}

func NewUpdateEventController(service UpdateEventService) UpdateEventController {
	return &updateEventController{Service: service}
}

// @Summary Update an existing event
// @Description Update an existing event
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body dto.UpdateEventRequest true "Event update data"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id} [put]
func (c *updateEventController) Handle(ctx *gin.Context) {
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

	// Bind and validate request body
	var updateReq dto.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, "validation.body.failed", err.Error())
		return
	}

	// Execute service
	updatedEvent, err := c.Service.Execute(ctx, eventID, &updateReq)
	if err != nil {
		switch err.Error() {
		case "event not found":
			response.Error(ctx, http.StatusNotFound, "event.not_found", "Event not found")
		case "unauthorized":
			response.Error(ctx, http.StatusForbidden, "event.unauthorized", "You are not authorized to update this event")
		case "user authentication required":
			response.Error(ctx, http.StatusUnauthorized, "auth.required", "Authentication required")
		default:
			response.Error(ctx, http.StatusInternalServerError, "update.event.fail", err.Error())
		}
		return
	}

	// Convert to response DTO
	eventResponse := mappers.ToEventResponse(updatedEvent)
	response.Success(ctx, http.StatusOK, "update.event.success", eventResponse, nil)
}
