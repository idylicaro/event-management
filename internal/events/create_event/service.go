package create_event

import (
	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/mappers"
)

// Interface do serviço de eventos
type CreateEventService interface {
	Execute(event *dto.CreateEventRequest) error
}

// Estrutura do serviço de eventos
type createEventService struct {
	repo CreateEventRepository
}

// Nova instância do serviço de eventos
func NewCreateEventService(repo CreateEventRepository) CreateEventService {
	return &createEventService{repo}
}

// Implementação do método CreateEvent
func (s *createEventService) Execute(req *dto.CreateEventRequest) error {
	event := mappers.ToEventModel(req)

	if err := event.Validate(); err != nil {
		return err
	}
	return s.repo.Execute(event)
}
