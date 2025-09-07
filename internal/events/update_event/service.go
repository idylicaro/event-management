package update_event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/mappers"
	"github.com/idylicaro/event-management/internal/middleware"
	"github.com/idylicaro/event-management/internal/models"
)

type updateEventService struct {
	repo UpdateEventRepository
}

func NewUpdateEventService(repo UpdateEventRepository) UpdateEventService {
	return &updateEventService{repo}
}

func (s *updateEventService) Execute(ctx *gin.Context, eventID int64, req *dto.UpdateEventRequest) (*models.Event, error) {
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		return nil, fmt.Errorf("user authentication required")
	}

	// Verify the event exists and belongs to the user
	existingEvent, err := s.repo.GetByIDAndUserID(eventID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Map update request to event model
	updatedEvent := mappers.ToEventModelFromUpdateWithID(req, eventID, user.ID)

	// Preserve original creation time and user ID
	updatedEvent.CreatedAt = existingEvent.CreatedAt
	updatedEvent.UserID = existingEvent.UserID

	// Validate the updated event
	if err := updatedEvent.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update in repository
	if err := s.repo.Execute(updatedEvent); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return updatedEvent, nil
}
