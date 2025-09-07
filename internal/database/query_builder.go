package database

import (
	"fmt"
	"strings"
)

// QueryBuilder helps construct SQL queries safely with proper parameterization
type QueryBuilder struct {
	baseQuery  string
	conditions []string
	args       []interface{}
	orderBy    []string
	limit      *int
	offset     *int
	argIndex   int
}

// NewQueryBuilder creates a new query builder instance
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:  baseQuery,
		conditions: make([]string, 0),
		args:       make([]interface{}, 0),
		orderBy:    make([]string, 0),
		argIndex:   1,
	}
}

// AddCondition adds a WHERE condition with proper parameterization
func (qb *QueryBuilder) AddCondition(condition string, arg interface{}) *QueryBuilder {
	placeholder := fmt.Sprintf("$%d", qb.argIndex)
	qb.conditions = append(qb.conditions, strings.Replace(condition, "?", placeholder, 1))
	qb.args = append(qb.args, arg)
	qb.argIndex++
	return qb
}

// AddLikeCondition adds a LIKE condition with proper parameterization
func (qb *QueryBuilder) AddLikeCondition(column string, value string) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s ILIKE ?", column), "%"+value+"%")
}

// AddOrderBy adds an ORDER BY clause
func (qb *QueryBuilder) AddOrderBy(column, direction string) *QueryBuilder {
	direction = strings.ToUpper(strings.TrimSpace(direction))
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}
	qb.orderBy = append(qb.orderBy, fmt.Sprintf("%s %s", column, direction))
	return qb
}

// SetLimit sets the LIMIT clause
func (qb *QueryBuilder) SetLimit(limit int) *QueryBuilder {
	qb.limit = &limit
	return qb
}

// SetOffset sets the OFFSET clause
func (qb *QueryBuilder) SetOffset(offset int) *QueryBuilder {
	qb.offset = &offset
	return qb
}

// Build constructs the final SQL query with all conditions
func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := qb.baseQuery

	if len(qb.conditions) > 0 {
		if strings.Contains(strings.ToUpper(query), "WHERE") {
			query += " AND " + strings.Join(qb.conditions, " AND ")
		} else {
			query += " WHERE " + strings.Join(qb.conditions, " AND ")
		}
	}

	if len(qb.orderBy) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderBy, ", ")
	}

	if qb.limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", qb.argIndex)
		qb.args = append(qb.args, *qb.limit)
		qb.argIndex++
	}

	if qb.offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", qb.argIndex)
		qb.args = append(qb.args, *qb.offset)
		qb.argIndex++
	}

	return query, qb.args
}
