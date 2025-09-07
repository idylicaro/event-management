package update_event

import (
	"context"
	"fmt"

	"github.com/idylicaro/event-management/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type updateEventRepository struct {
	db *pgxpool.Pool
}

func NewUpdateEventRepository(db *pgxpool.Pool) UpdateEventRepository {
	return &updateEventRepository{db}
}

// Update updates an existing event in the database
func (r *updateEventRepository) Execute(event *models.Event) error {
	query := `
		UPDATE events 
		SET title=$1, description=$2, location=$3, start_time=$4, end_time=$5, price=$6, updated_at=NOW() 
		WHERE id=$7 AND user_id=$8
		RETURNING updated_at`

	err := r.db.QueryRow(
		context.Background(),
		query,
		event.Title,
		event.Description,
		event.Location,
		event.StartTime,
		event.EndTime,
		event.Price,
		event.ID,
		event.UserID,
	).Scan(&event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// GetByIDAndUserID retrieves an event by ID and user ID to verify ownership
func (r *updateEventRepository) GetByIDAndUserID(eventID, userID int64) (*models.Event, error) {
	query := `
		SELECT id, title, description, location, start_time, end_time, price, user_id, created_at, updated_at
		FROM events 
		WHERE id=$1 AND user_id=$2`

	event := &models.Event{}
	err := r.db.QueryRow(context.Background(), query, eventID, userID).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.Price,
		&event.UserID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("event not found or unauthorized: %w", err)
	}

	return event, nil
}
