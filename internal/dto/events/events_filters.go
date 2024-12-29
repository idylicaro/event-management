package dto

// Get query filters
type EventFilters struct {
	Title     string `form:"title"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Page      int    `form:"page"`
	Limit     int    `form:"limit" binding:"lte=100"`
	Sort      string `form:"sort"`
}
