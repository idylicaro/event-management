package models

import "time"

type User struct {
	ID                int64     `json:"id"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	Role              string    `json:"role"` // "admin", "event_owner","ticket_buyer".
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
