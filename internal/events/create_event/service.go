package create_event

import (
	models "github.com/idylicaro/event-management/internal/models"
)

// Interface do serviço de eventos
type CreateEventService interface {
	CreateEvent(event *models.Event) error
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
func (s *createEventService) CreateEvent(event *models.Event) error {
	return s.repo.CreateEvent(event)
}
