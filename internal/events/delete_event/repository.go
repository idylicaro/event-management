package delete_event

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type deleteEventRepository struct {
	db *pgxpool.Pool
}

func NewDeleteEventRepository(db *pgxpool.Pool) DeleteEventRepository {
	return &deleteEventRepository{db}
}

// GetByIDAndUserID checks if the event exists and belongs to the user
func (r *deleteEventRepository) GetByIDAndUserID(eventID, userID int64) error {
	query := `SELECT id FROM events WHERE id = $1 AND user_id = $2`

	var id int64
	err := r.db.QueryRow(context.Background(), query, eventID, userID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// Execute deletes the event from the database
func (r *deleteEventRepository) Execute(eventID, userID int64) error {
	query := `DELETE FROM events WHERE id = $1 AND user_id = $2`

	cmdTag, err := r.db.Exec(context.Background(), query, eventID, userID)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	if cmdTag.RowsAffected() == 0 {
		return nil // Event was already deleted or doesn't exist
	}

	return nil
}
