package list_events

import "github.com/idylicaro/event-management/internal/models"

// GetEventsService defines the service for getting events
type GetEventsService interface {
	Execute(filters models.EventFilters) ([]models.Event, error)
}

// getEventsService is the concrete implementation of GetEventsService
type getEventsService struct {
	repo GetEventsRepository
}

// NewGetEventsService creates a new instance of getEventsService
func NewGetEventsService(repo GetEventsRepository) GetEventsService {
	return &getEventsService{repo}
}

// Execute gets events based on filters
func (s *getEventsService) Execute(filters models.EventFilters) ([]models.Event, error) {
	return s.repo.Execute(filters)
}
