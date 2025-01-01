package models

import (
	"fmt"
	"time"
)

// Event representa um evento que será manipulado pela API.
type Event struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Location    string    `json:"location" db:"location"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Price       float64   `json:"price" db:"price"`
	UserID      int64     `json:"user_id"` // Relation to the user that owns the event
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (e *Event) Validate() error {
	if e.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if e.Location == "" {
		return fmt.Errorf("location cannot be empty")
	}
	if e.StartTime.IsZero() || e.EndTime.IsZero() {
		return fmt.Errorf("start time and end time must be valid")
	}
	if e.StartTime.After(e.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}
	return nil
}
