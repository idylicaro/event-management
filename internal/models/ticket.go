package models

import "time"

// Ticket represents a ticket purchased by a user.
type Ticket struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"` // Relationship with the event
	UserID    int64     `json:"user_id"`  // Relationship with the ticket buyer (User)
	Status    string    `json:"status"`   // Ticket status (e.g., "paid", "pending", "canceled")
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
