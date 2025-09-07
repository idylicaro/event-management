package update_event

import (
	"context"
	"fmt"
	"strings"
	"time"

	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type updateEventRepository struct {
	db *pgxpool.Pool
}

func NewUpdateEventRepository(db *pgxpool.Pool) UpdateEventRepository {
	return &updateEventRepository{db}
}

// Execute updates an existing event in the database with only provided fields
func (r *updateEventRepository) Execute(eventID, userID int64, req *dto.UpdateEventRequest) (*models.Event, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	// Build dynamic SET clause based on provided fields
	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title=$%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description=$%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Location != nil {
		setParts = append(setParts, fmt.Sprintf("location=$%d", argIndex))
		args = append(args, *req.Location)
		argIndex++
	}
	if req.StartTime != nil {
		setParts = append(setParts, fmt.Sprintf("start_time=$%d", argIndex))
		args = append(args, *req.StartTime)
		argIndex++
	}
	if req.EndTime != nil {
		setParts = append(setParts, fmt.Sprintf("end_time=$%d", argIndex))
		args = append(args, *req.EndTime)
		argIndex++
	}
	if req.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price=$%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	// If no fields to update, return error
	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Always update the updated_at field
	setParts = append(setParts, fmt.Sprintf("updated_at=$%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE conditions
	args = append(args, eventID, userID)

	query := fmt.Sprintf(`
		UPDATE events 
		SET %s 
		WHERE id=$%d AND user_id=$%d
		RETURNING id, title, description, location, start_time, end_time, price, user_id, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex, argIndex+1)

	event := &models.Event{}
	err := r.db.QueryRow(context.Background(), query, args...).Scan(
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
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
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
