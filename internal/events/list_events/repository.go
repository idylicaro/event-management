package list_events

import (
	"context"
	"fmt"
	"log"

	dto "github.com/idylicaro/event-management/internal/dto/events"
	"github.com/idylicaro/event-management/internal/helpers"
	"github.com/idylicaro/event-management/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type getEventsRepository struct {
	db *pgxpool.Pool
}

func NewGetEventsRepository(db *pgxpool.Pool) GetEventsRepository {
	return &getEventsRepository{db}
}

func (r *getEventsRepository) Execute(filters dto.EventFilters) ([]models.Event, helpers.PaginationMeta, error) {
	baseQuery := `
		SELECT id, title, description, location, start_time, end_time, price, created_at, updated_at
		FROM events
		WHERE 1=1
	`
	countQuery := `
		SELECT COUNT(*) FROM events WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filters.Title != "" {
		filter := ` AND title ILIKE $` + fmt.Sprint(argIndex)
		baseQuery += filter
		countQuery += filter
		args = append(args, "%"+filters.Title+"%")
		argIndex++
	}
	if filters.StartTime != "" {
		filter := ` AND start_time >= $` + fmt.Sprint(argIndex)
		baseQuery += filter
		countQuery += filter
		args = append(args, filters.StartTime)
		argIndex++
	}
	if filters.EndTime != "" {
		filter := ` AND end_time <= $` + fmt.Sprint(argIndex)
		baseQuery += filter
		countQuery += filter
		args = append(args, filters.EndTime)
		argIndex++
	}

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

	// Pagination
	paginationQuery, paginationArgs := helpers.BuildPaginationQuery(helpers.PaginationParams{
		Limit: filters.Limit,
		Page:  filters.Page,
	})
	baseQuery += paginationQuery
	args = append(args, paginationArgs...)

	log.Printf("Query: %s\n", baseQuery)
	log.Printf("Count Query: %s\n", countQuery)
	log.Printf("Args: %v\n", args)

	var totalItems int
	err := r.db.QueryRow(context.Background(), countQuery, args[:argIndex-1]...).Scan(&totalItems)
	if err != nil {
		return nil, helpers.PaginationMeta{}, err
	}

	rows, err := r.db.Query(context.Background(), baseQuery, args...)
	if err != nil {
		return nil, helpers.PaginationMeta{}, err
	}
	defer rows.Close()

	// Mapping
	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Price, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, helpers.PaginationMeta{}, err
		}
		events = append(events, event)
	}
	meta := helpers.CalculatePaginationMeta(totalItems, filters.Limit, filters.Page)

	return events, meta, nil
}
