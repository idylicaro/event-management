package list_events

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type getEventsController struct {
	Service GetEventsService
}

func NewGetEventsController(service GetEventsService) GetEventsController {
	return &getEventsController{Service: service}
}

// @Summary      List Events
// @Description  Retrieves a list of events, optionally filtered by title, date range, and paginated.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        title      query     string  false  "Filter by event title"
// @Param        start_time query     string  false  "Filter by start date (format: YYYY-MM-DD HH:MM:SS)"
// @Param        end_time   query     string  false  "Filter by end date (format: YYYY-MM-DD HH:MM:SS)"
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        limit      query     int     false  "Number of items per page (default: 10, max: 100)"
// @Param        sort       query     string  false  "Sort order (e.g., 'date:asc', 'date:desc')"
// @Success      200        {object}  []models.Event "List of events"
// @Failure      400        {object}  map[string]string   "Invalid request parameters"
// @Failure      500        {object}  map[string]string   "Internal server error"
// @Router       /events [get]
func (c *getEventsController) Handle(ctx *gin.Context) {
	var filters dto.EventFilters
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		response.Error(ctx, http.StatusBadRequest, "validation.query.failed", err.Error())
		return
	}

	events, meta, err := c.Service.Execute(filters)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "get.events.fail", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "get.events.success", events, meta)
}
