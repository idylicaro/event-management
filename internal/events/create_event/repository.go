package create_event

import (
	"context"

	models "github.com/idylicaro/event-management/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateEventRepository defines the interface for event repository
type CreateEventRepository interface {
	CreateEvent(event *models.Event) error
}

// createEventRepository is the concrete implementation of CreateEventRepository
type createEventRepository struct {
	db *pgxpool.Pool
}

// NewEventRepository creates a new instance of createEventRepository
func NewEventRepository(db *pgxpool.Pool) CreateEventRepository {
	return &createEventRepository{db}
}

// CreateEvent inserts a new event into the database
func (r *createEventRepository) CreateEvent(event *models.Event) error {
	query := `INSERT INTO events (title, description, location, start_time, end_time, price) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(context.Background(), query, event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.Price).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
