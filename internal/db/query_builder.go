package db

import (
	"fmt"
	"strings"

	"github.com/zulfikawr/vault/internal/models"
)

type QueryParams struct {
	Page    int
	PerPage int
	Sort    string
	Filter  string
}

func (e *Executor) ListRecords(collectionName string, params QueryParams) ([]*models.Record, int, error) {
	col, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return nil, 0, fmt.Errorf("collection %s not found", collectionName)
	}

	if params.Page <= 0 { params.Page = 1 }
	if params.PerPage <= 0 { params.PerPage = 30 }

	// Count total
	var total int
	err := e.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", collectionName)).Scan(&total)
	if err != nil {
		return nil, 0, err
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
	rows, err := e.db.Query(query, params.PerPage, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
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
