package update_event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
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
	_, err := s.repo.GetByIDAndUserID(eventID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Validate only the fields that are being updated
	if err := s.validatePartialUpdate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update in repository with partial data
	updatedEvent, err := s.repo.Execute(eventID, user.ID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return updatedEvent, nil
}

// validatePartialUpdate validates only the fields that are being updated
func (s *updateEventService) validatePartialUpdate(req *dto.UpdateEventRequest) error {
	if req.Title != nil && *req.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if req.Location != nil && *req.Location == "" {
		return fmt.Errorf("location cannot be empty")
	}

	// If both start and end times are provided, validate their relationship
	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}

	// If only one time is provided, we need to get the current event to validate
	// This will be handled by the database constraints or by fetching current event

	if req.Price != nil && *req.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	return nil
}
