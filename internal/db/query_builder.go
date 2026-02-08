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
}

func (e *Executor) ListRecords(ctx context.Context, collectionName string, params QueryParams) ([]*models.Record, int, error) {
	col, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return nil, 0, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	if params.Page <= 0 { params.Page = 1 }
	if params.PerPage <= 0 { params.PerPage = 30 }

	// Count total
	var total int
	err := e.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", collectionName)).Scan(&total)
	if err != nil {
		return nil, 0, core.NewError(http.StatusInternalServerError, "DB_COUNT_FAILED", "Failed to count records").WithDetails(map[string]any{"error": err.Error()})
	}

	columns := []string{"id", "created", "updated"}
	for _, f := range col.Fields {
		columns = append(columns, f.Name)
	}

	query := fmt.Sprintf("SELECT %s FROM %s LIMIT ? OFFSET ?", 
		strings.Join(columns, ", "), 
		collectionName,
	)

	offset := (params.Page - 1) * params.PerPage
	rows, err := e.db.QueryContext(ctx, query, params.PerPage, offset)
	if err != nil {
		return nil, 0, core.NewError(http.StatusInternalServerError, "RECORD_LIST_FAILED", "Failed to list records").WithDetails(map[string]any{"error": err.Error()})
	}
	defer rows.Close()

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

	return records, total, nil
}