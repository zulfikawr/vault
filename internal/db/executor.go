package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/models"
)

type Executor struct {
	db       *sql.DB
	registry *SchemaRegistry
}

func NewExecutor(db *sql.DB, registry *SchemaRegistry) *Executor {
	return &Executor{db: db, registry: registry}
}

func (e *Executor) CreateRecord(ctx context.Context, collectionName string, data map[string]any) (*models.Record, error) {
	col, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return nil, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	id := uuid.New().String()
	data["id"] = id

	record := &models.Record{
		ID:         id,
		Collection: collectionName,
		Data:       data,
	}

	hooks := GetHooks(collectionName)
	if err := hooks.TriggerBeforeCreate(ctx, record); err != nil {
		return nil, err
	}

	columns := []string{"id"}
	placeholders := []string{"?"}
	values := []any{id}

	for _, f := range col.Fields {
		if val, ok := record.Data[f.Name]; ok {
			columns = append(columns, f.Name)
			placeholders = append(placeholders, "?")
			values = append(values, val)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING created, updated",
		collectionName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	err := e.db.QueryRowContext(ctx, query, values...).Scan(&record.Created, &record.Updated)
	if err != nil {
		return nil, core.NewError(http.StatusInternalServerError, "RECORD_CREATE_FAILED", "Failed to create record").WithDetails(map[string]any{"error": err.Error()})
	}

	hooks.TriggerAfterCreate(ctx, record)
	return record, nil
}

func (e *Executor) FindRecordByID(ctx context.Context, collectionName string, id string) (*models.Record, error) {
	col, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return nil, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	columns := []string{"id", "created", "updated"}
	for _, f := range col.Fields {
		columns = append(columns, f.Name)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", strings.Join(columns, ", "), collectionName)
	row := e.db.QueryRowContext(ctx, query, id)

	vals := make([]any, len(columns))
	valPtrs := make([]any, len(columns))
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	if err := row.Scan(valPtrs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, core.NewError(http.StatusNotFound, "RECORD_NOT_FOUND", "Record not found")
		}
		return nil, core.NewError(http.StatusInternalServerError, "RECORD_FETCH_FAILED", "Failed to fetch record").WithDetails(map[string]any{"error": err.Error()})
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

	return record, nil
}