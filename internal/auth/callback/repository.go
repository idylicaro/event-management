package callback

import (
	"context"

	"github.com/idylicaro/event-management/internal/auth/providers"
	"github.com/idylicaro/event-management/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindOrCreateUser(data providers.UserInfo) (models.User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

func (r *repository) FindOrCreateUser(data providers.UserInfo) (models.User, error) {
	var user models.User
	query := `SELECT id, email, name, profile_picture_url, role, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(context.Background(), query, data.Email).Scan(&user.ID, &user.Email, &user.Name, &user.ProfilePictureURL, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			query = `INSERT INTO users (email, name, profile_picture_url, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
			err = r.db.QueryRow(context.Background(), query, data.Email, data.Name, data.PictureURL, "ticket_buyer").Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				return models.User{}, err
			}
			user.Email = data.Email
			user.Name = data.Name
			user.ProfilePictureURL = data.PictureURL
			user.Role = "ticket_buyer"
		} else {
			return models.User{}, err
		}
	}
	return user, nil
}
