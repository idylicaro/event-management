package create_event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/mappers"
	"github.com/idylicaro/event-management/internal/middleware"
)

// Estrutura do serviço de eventos
type createEventService struct {
	repo CreateEventRepository
}

// Nova instância do serviço de eventos
func NewCreateEventService(repo CreateEventRepository) CreateEventService {
	return &createEventService{repo}
}

// Implementação do método CreateEvent
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
