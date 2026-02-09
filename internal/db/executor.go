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
	"github.com/zulfikawr/vault/internal/realtime"
)

type Executor struct {
	db       *sql.DB
	registry *SchemaRegistry
	hub      *realtime.Hub
}

func NewExecutor(db *sql.DB, registry *SchemaRegistry, hub *realtime.Hub) *Executor {
	return &Executor{db: db, registry: registry, hub: hub}
}

func (e *Executor) broadcast(action string, collection string, record *models.Record) {
	if e.hub != nil {
		e.hub.Broadcast(&realtime.Message{
			Action:     action,
			Collection: collection,
			Record:     record,
		})
	}
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
	e.broadcast("create", collectionName, record)
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

func (e *Executor) expandRecords(ctx context.Context, collection *models.Collection, records []*models.Record, expand string) {
	if expand == "" || len(records) == 0 {
		return
	}

	expandFields := strings.Split(expand, ",")
	for _, fieldName := range expandFields {
		fieldName = strings.TrimSpace(fieldName)

		// 1. Find the relationship metadata
		var relField *models.Field
		for _, f := range collection.Fields {
			if f.Name == fieldName && f.Type == models.FieldTypeRelation {
				relField = &f
				break
			}
		}
		if relField == nil { continue }

		targetCol := ""
		if options, ok := relField.Options.(map[string]any); ok {
			targetCol, _ = options["collection"].(string)
		}
		if targetCol == "" { continue }

		// 2. Collect all unique IDs to fetch
		relIDs := make([]string, 0)
		idMap := make(map[string]bool)
		for _, r := range records {
			if id, ok := r.Data[fieldName].(string); ok && id != "" {
				if !idMap[id] {
					relIDs = append(relIDs, id)
					idMap[id] = true
				}
			}
		}
		if len(relIDs) == 0 { continue }

		// 3. Batch fetch records (Fixing N+1)
		placeholders := make([]string, len(relIDs))
		args := make([]any, len(relIDs))
		for i, id := range relIDs {
			placeholders[i] = "?"
			args[i] = id
		}

		// Simplified batch fetch logic for this prototype
		// We use ListRecords internally with a raw IN filter for efficiency
		filter := fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ","))
		targetRecords, _, err := e.ListRecords(ctx, targetCol, QueryParams{
			Filter:  filter,
			PerPage: len(relIDs),
		})
		if err != nil { continue }

		// 4. Map back to original records
		resMap := make(map[string]*models.Record)
		for _, tr := range targetRecords {
			resMap[tr.ID] = tr
		}

		for _, r := range records {
			if id, ok := r.Data[fieldName].(string); ok {
				if expanded, exists := resMap[id]; exists {
					r.Expand[fieldName] = expanded
				}
			}
		}
	}
}

func (e *Executor) UpdateRecord(ctx context.Context, collectionName string, id string, data map[string]any) (*models.Record, error) {
	col, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return nil, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	record, err := e.FindRecordByID(ctx, collectionName, id)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		if k != "id" && k != "created" && k != "updated" {
			record.Data[k] = v
		}
	}

	columns := []string{"updated = (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))"}
	values := []any{}

	for _, f := range col.Fields {
		if val, ok := record.Data[f.Name]; ok {
			columns = append(columns, fmt.Sprintf("%s = ?", f.Name))
			values = append(values, val)
		}
	}
	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ? RETURNING updated",
		collectionName,
		strings.Join(columns, ", "),
	)

	err = e.db.QueryRowContext(ctx, query, values...).Scan(&record.Updated)
	if err != nil {
		return nil, core.NewError(http.StatusInternalServerError, "RECORD_UPDATE_FAILED", "Failed to update record").WithDetails(map[string]any{"error": err.Error()})
	}

	e.broadcast("update", collectionName, record)
	return record, nil
}

func (e *Executor) DeleteRecord(ctx context.Context, collectionName string, id string) error {
	_, ok := e.registry.GetCollection(collectionName)
	if !ok {
		return core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", collectionName)
	result, err := e.db.ExecContext(ctx, query, id)
	if err != nil {
		return core.NewError(http.StatusInternalServerError, "RECORD_DELETE_FAILED", "Failed to delete record").WithDetails(map[string]any{"error": err.Error()})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return core.NewError(http.StatusNotFound, "RECORD_NOT_FOUND", "Record not found")
	}

	e.broadcast("delete", collectionName, &models.Record{ID: id})
	return nil
}