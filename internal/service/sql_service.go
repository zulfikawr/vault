package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/zulfikawr/vault/internal/errors"
)

type SqlService struct {
	db *sql.DB
}

func NewSqlService(db *sql.DB) *SqlService {
	return &SqlService{db: db}
}

type QueryResult struct {
	Columns []string         `json:"columns"`
	Rows    []map[string]any `json:"rows"`
}

func (s *SqlService) Execute(ctx context.Context, query string) (*QueryResult, error) {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "SQL_EXECUTION_FAILED", "Failed to execute SQL query").WithDetails(map[string]any{"error": err.Error(), "query": query})
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "SQL_COLUMNS_FAILED", "Failed to get result columns").WithDetails(map[string]any{"error": err.Error()})
	}

	resultRows := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, errors.NewError(http.StatusInternalServerError, "SQL_SCAN_FAILED", "Failed to scan result row").WithDetails(map[string]any{"error": err.Error()})
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		resultRows = append(resultRows, rowMap)
	}

	return &QueryResult{
		Columns: columns,
		Rows:    resultRows,
	}, nil
}
