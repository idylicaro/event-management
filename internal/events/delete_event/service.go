package delete_event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/middleware"
)

type deleteEventService struct {
	repo DeleteEventRepository
}

func NewDeleteEventService(repo DeleteEventRepository) DeleteEventService {
	return &deleteEventService{repo}
}

func (s *deleteEventService) Execute(ctx *gin.Context, eventID int64) error {
	// Get user from context (must be authenticated)
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		return fmt.Errorf("unauthorized")
	}

	// Check if event exists and belongs to the user
	if err := s.repo.GetByIDAndUserID(eventID, user.ID); err != nil {
		if err.Error() == "sql: no rows in result set" || err.Error() == "no rows in result set" {
			return fmt.Errorf("event not found")
		}
		return err
	}

	// Delete the event
	return s.repo.Execute(eventID, user.ID)
}
