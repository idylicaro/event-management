package list_events

import (
	"context"
	"log"
	"strings"

	"github.com/idylicaro/event-management/internal/database"
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
	// Build main query using QueryBuilder
	qb := database.NewQueryBuilder("SELECT id, title, description, location, start_time, end_time, price, created_at, updated_at, user_id FROM events")

	// Add filters
	if filters.Title != "" {
		qb.AddLikeCondition("title", filters.Title)
	}
	if filters.StartTime != "" {
		qb.AddCondition("start_time >= ?", filters.StartTime)
	}
	if filters.EndTime != "" {
		qb.AddCondition("end_time <= ?", filters.EndTime)
	}

	// Add sorting
	sortMap := map[string]string{
		"date:asc":   "start_time ASC",
		"date:desc":  "start_time DESC",
		"title:asc":  "title ASC",
		"title:desc": "title DESC",
	}
	if filters.Sort != "" {
		if order, ok := sortMap[filters.Sort]; ok {
			parts := strings.Split(order, " ")
			if len(parts) == 2 {
				qb.AddOrderBy(parts[0], parts[1])
			}
		}
	} else {
		qb.AddOrderBy("start_time", "ASC") // Default
	}

	// Add pagination
	if filters.Limit > 0 {
		qb.SetLimit(filters.Limit)
	}
	if filters.Page > 1 && filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		qb.SetOffset(offset)
	}

	// Build queries
	mainQuery, mainArgs := qb.Build()

	// Build count query (same conditions, no ORDER BY/LIMIT/OFFSET)
	countQb := database.NewQueryBuilder("SELECT COUNT(*) FROM events")
	if filters.Title != "" {
		countQb.AddLikeCondition("title", filters.Title)
	}
	if filters.StartTime != "" {
		countQb.AddCondition("start_time >= ?", filters.StartTime)
	}
	if filters.EndTime != "" {
		countQb.AddCondition("end_time <= ?", filters.EndTime)
	}
	countQuery, countArgs := countQb.Build()

	log.Printf("Main Query: %s\n", mainQuery)
	log.Printf("Count Query: %s\n", countQuery)
	log.Printf("Main Args: %v\n", mainArgs)
	log.Printf("Count Args: %v\n", countArgs)

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(context.Background(), countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		return nil, helpers.PaginationMeta{}, err
	}

	// Execute main query
	rows, err := r.db.Query(context.Background(), mainQuery, mainArgs...)
	if err != nil {
		return nil, helpers.PaginationMeta{}, err
	}
	defer rows.Close()

	// Scan results
	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Price, &event.CreatedAt, &event.UpdatedAt, &event.UserID)
		if err != nil {
			return nil, helpers.PaginationMeta{}, err
		}
		events = append(events, event)
	}

	// Calculate pagination metadata
	meta := helpers.CalculatePaginationMeta(totalItems, filters.Limit, filters.Page)

	return events, meta, nil
}
