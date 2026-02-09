package db

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/models"
)

type QueryParams struct {
	Page    int
	PerPage int
	Sort    string
	Filter  string
	Expand  string
}

func (e *Executor) ListRecords(ctx context.Context, collectionName string, params QueryParams) ([]*models.Record, int, error) {
	col, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return nil, 0, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	if params.Page <= 0 { params.Page = 1 }
	if params.PerPage <= 0 { params.PerPage = 30 }

	// Prepare WHERE clause and values
	whereClauses := []string{"1=1"}
	whereValues := []any{}

	if params.Filter != "" {
		clause, values, err := e.parseSafeFilter(col, params.Filter)
		if err != nil {
			return nil, 0, err
		}
		whereClauses = append(whereClauses, clause)
		whereValues = append(whereValues, values...)
	}

	whereQuery := strings.Join(whereClauses, " AND ")

	// Count total
	var total int
	err := e.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", collectionName, whereQuery), whereValues...).Scan(&total)
	if err != nil {
		return nil, 0, core.NewError(http.StatusInternalServerError, "DB_COUNT_FAILED", "Failed to count records").WithDetails(map[string]any{"error": err.Error()})
	}

	columns := []string{"id", "created", "updated"}
	for _, f := range col.Fields {
		columns = append(columns, f.Name)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT ? OFFSET ?", 
		strings.Join(columns, ", "), 
		collectionName,
		whereQuery,
	)

	offset := (params.Page - 1) * params.PerPage
	allValues := append(whereValues, params.PerPage, offset)
	rows, err := e.db.QueryContext(ctx, query, allValues...)
	if err != nil {
		return nil, 0, core.NewError(http.StatusInternalServerError, "RECORD_LIST_FAILED", "Failed to list records").WithDetails(map[string]any{"error": err.Error()})
	}
	defer func() { _ = rows.Close() }()

	var records []*models.Record
	for rows.Next() {
		vals := make([]any, len(columns))
		valPtrs := make([]any, len(columns))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}

		if err := rows.Scan(valPtrs...); err != nil {
			return nil, 0, core.NewError(http.StatusInternalServerError, "RECORD_SCAN_FAILED", "Failed to scan record").WithDetails(map[string]any{"error": err.Error()})
		}

		record := &models.Record{
			ID:         vals[0].(string),
			Collection: collectionName,
			Created:    vals[1].(string),
			Updated:    vals[2].(string),
			Data:       make(map[string]any),
		}

		for i, f := range col.Fields {
			record.Data[f.Name] = vals[i+3]
		}
		records = append(records, record)
	}

	if params.Expand != "" {
		e.expandRecords(ctx, col, records, params.Expand)
	}

	return records, total, nil
}

// parseSafeFilter implements basic parameterized filtering to prevent SQL injection
func (e *Executor) parseSafeFilter(col *models.Collection, filter string) (string, []any, error) {
	// Simple support for: field = 'value' or field != 'value'
	operators := []string{" != ", " = ", " > ", " < ", " >= ", " <= "}
	
	for _, op := range operators {
		if strings.Contains(filter, op) {
			parts := strings.Split(filter, op)
			if len(parts) != 2 { continue }

			fieldName := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 1. Validate field name exists in collection
			validField := false
			if fieldName == "id" {
				validField = true
			} else {
				for _, f := range col.Fields {
					if f.Name == fieldName {
						validField = true
						break
					}
				}
			}

			if !validField {
				return "", nil, core.NewError(http.StatusBadRequest, "INVALID_FILTER", fmt.Sprintf("Unknown field in filter: %s", fieldName))
			}

			// 2. Clean value (remove single quotes if present)
			value = strings.Trim(value, "'")

			// 3. Return parameterized clause
			return fmt.Sprintf("%s %s ?", fieldName, strings.TrimSpace(op)), []any{value}, nil
		}
	}

	return "", nil, core.NewError(http.StatusBadRequest, "UNSUPPORTED_FILTER", "Filter format not supported. Use 'field = value'")
}
