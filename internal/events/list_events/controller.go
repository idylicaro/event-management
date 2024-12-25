package list_events

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/models"
)

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
func GetEventsController(svc GetEventsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters models.EventFilters
		if err := c.ShouldBindQuery(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		events, err := svc.Execute(filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"events": events})
	}
}
