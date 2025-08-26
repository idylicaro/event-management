package models

import "time"

// Payment represents the payment made by a user for a ticket.
type Payment struct {
	ID          int64     `json:"id"`
	TicketID    int64     `json:"ticket_id"` // Relationship with the ticket
	UserID      int64     `json:"user_id"`   // Relationship with the buyer
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`       // Payment status (e.g., "completed", "failed", etc.)
	PaymentDate time.Time `json:"payment_date"` // Date of the payment
	StripeID    string    `json:"stripe_id"`    // Transaction ID in Stripe
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
