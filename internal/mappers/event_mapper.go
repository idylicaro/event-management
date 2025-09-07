package mappers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/models"
)

// EventMapper provides specialized mapping for event-related operations
type EventMapper struct {
	toResponseMapper *FieldMapper[*models.Event, *dto.EventResponse]
	fromCreateMapper *FieldMapper[*dto.CreateEventRequest, *models.Event]
	fromUpdateMapper *FieldMapper[*dto.UpdateEventRequest, *models.Event]
}

// NewEventMapper creates a new event mapper with configured field mappings
func NewEventMapper() *EventMapper {
	return &EventMapper{
		toResponseMapper: NewFieldMapper[*models.Event, *dto.EventResponse](),
		fromCreateMapper: NewFieldMapper[*dto.CreateEventRequest, *models.Event](),
		fromUpdateMapper: NewFieldMapper[*dto.UpdateEventRequest, *models.Event](),
	}
}

// Singleton instance for backward compatibility
var defaultEventMapper = NewEventMapper()

// ToEventResponse converts a model Event to EventResponse DTO
func ToEventResponse(event *models.Event) *dto.EventResponse {
	if event == nil {
		return nil
	}

	return defaultEventMapper.toResponseMapper.Map(event)
}

// ToEventModel converts CreateEventRequest DTO to Event model
func ToEventModel(req *dto.CreateEventRequest) *models.Event {
	if req == nil {
		return nil
	}

	return defaultEventMapper.fromCreateMapper.Map(req)
}

// ToEventModelFromUpdate converts UpdateEventRequest DTO to Event model with ID extraction
func ToEventModelFromUpdate(ctx *gin.Context, req *dto.UpdateEventRequest) (*models.Event, error) {
	if req == nil {
		return nil, nil
	}

	// Extract ID from path parameter
	idParam := ctx.Param("id")
	if idParam == "" {
		return nil, fmt.Errorf("event ID is required")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	// Map fields from request
	event := defaultEventMapper.fromUpdateMapper.Map(req)
	event.ID = id
	event.UpdatedAt = time.Now()

	return event, nil
}

// ToEventModelFromUpdateWithID converts UpdateEventRequest to Event model with explicit ID
func ToEventModelFromUpdateWithID(req *dto.UpdateEventRequest, id int64, userID int64) *models.Event {
	if req == nil {
		return nil
	}

	event := defaultEventMapper.fromUpdateMapper.Map(req)
	event.ID = id
	event.UserID = userID
	event.UpdatedAt = time.Now()

	return event
}

// MapEventSliceToResponse converts a slice of Event models to EventResponse DTOs
func MapEventSliceToResponse(events []*models.Event) []*dto.EventResponse {
	if events == nil {
		return nil
	}

	responses := make([]*dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = ToEventResponse(event)
	}
	return responses
}
