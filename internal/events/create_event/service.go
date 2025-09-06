package create_event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/mappers"
	"github.com/idylicaro/event-management/internal/middleware"
)

type createEventService struct {
	repo CreateEventRepository
}

func NewCreateEventService(repo CreateEventRepository) CreateEventService {
	return &createEventService{repo}
}

func (s *createEventService) Execute(ctx *gin.Context, req *dto.CreateEventRequest) error {
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		return fmt.Errorf("user authentication required")
	}

	event := mappers.ToEventModel(req)

	// Set the user ID from the authenticated user
	event.UserID = user.ID

	if err := event.Validate(); err != nil {
		return err
	}
	return s.repo.Execute(event)
}
