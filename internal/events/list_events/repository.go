package list_events

import (
	"context"
	"fmt"

	"github.com/idylicaro/event-management/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GetEventsRepository interface {
	Execute(filters models.EventFilters) ([]models.Event, error)
}

type getEventsRepository struct {
	db *pgxpool.Pool
}

func NewGetEventsRepository(db *pgxpool.Pool) GetEventsRepository {
	return &getEventsRepository{db}
}

func (r *getEventsRepository) Execute(filters models.EventFilters) ([]models.Event, error) {
	baseQuery := `
		SELECT id, title, description, location, start_time, end_time, price, created_at, updated_at
		FROM events
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Adicionar filtros dinâmicos
	if filters.Title != "" {
		baseQuery += ` AND title ILIKE $` + fmt.Sprint(argIndex)
		args = append(args, "%"+filters.Title+"%")
		argIndex++
	}
	if filters.StartTime != "" {
		baseQuery += ` AND start_time >= $` + fmt.Sprint(argIndex)
		args = append(args, filters.StartTime)
		argIndex++
	}
	if filters.EndTime != "" {
		baseQuery += ` AND end_time <= $` + fmt.Sprint(argIndex)
		args = append(args, filters.EndTime)
		argIndex++
	}

	// Ordenação
	sortMap := map[string]string{
		"date:asc":   "start_time ASC",
		"date:desc":  "start_time DESC",
		"title:asc":  "title ASC",
		"title:desc": "title DESC",
	}
	if order, ok := sortMap[filters.Sort]; ok {
		baseQuery += ` ORDER BY ` + order
	} else {
		baseQuery += ` ORDER BY start_time ASC` // Default
	}

	// Paginação
	limit := filters.Limit
	if limit == 0 {
		limit = 10 // Valor padrão
	}
	offset := (filters.Page - 1) * limit
	baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Executar query
	rows, err := r.db.Query(context.Background(), baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Processar resultados
	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Price, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
