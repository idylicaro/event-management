package helpers

import "fmt"

// PaginationParams holds parameters for pagination
type PaginationParams struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type PaginationMeta struct {
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	NextPage    int `json:"next_page,omitempty"`
	PrevPage    int `json:"prev_page,omitempty"`
}

// BuildPaginationQuery builds the LIMIT, OFFSET by SQL clauses.
func BuildPaginationQuery(params PaginationParams) (string, []interface{}) {
	query := ""
	args := []interface{}{}
	argIndex := 1

	// Pagination
	limit := params.Limit
	if limit == 0 {
		limit = 10
	}
	if params.Page == 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	return query, args
}

func CalculatePaginationMeta(totalItems, limit, page int) PaginationMeta {
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	totalPages := totalItems / limit
	if totalItems%limit != 0 {
		totalPages++
	}

	meta := PaginationMeta{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
	}
	if page < totalPages {
		meta.NextPage = page + 1
	}
	if page > 1 {
		meta.PrevPage = page - 1
	}

	return meta
}
