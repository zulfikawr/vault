package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
)

type Repository struct {
	db        *sql.DB
	registry  *SchemaRegistry
	stmtCache *StatementCache
}

func NewRepository(db *sql.DB, registry *SchemaRegistry) *Repository {
	return &Repository{
		db:        db,
		registry:  registry,
		stmtCache: NewStatementCache(db),
	}
}

func (r *Repository) Close() {
	if r.stmtCache != nil {
		r.stmtCache.Close()
	}
}

type QueryParams struct {
	Page    int
	PerPage int
	Sort    string
	Filter  string
	Expand  string
}

func (r *Repository) CreateRecord(ctx context.Context, collectionName string, data map[string]any) (*models.Record, error) {
	col, ok := r.registry.GetCollection(collectionName)
	if !ok {
		return nil, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	// ID must be provided by the caller (Service layer)
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, errors.NewError(http.StatusInternalServerError, "MISSING_ID", "Record ID is required")
	}

	record := &models.Record{
		ID:         id,
		Collection: collectionName,
		Data:       data,
	}

	// Prepare data for insertion
	insertData := make(map[string]any)
	insertData["id"] = id
	for _, f := range col.Fields {
		if val, ok := record.Data[f.Name]; ok {
			insertData[f.Name] = val
		}
	}

	qb := NewQueryBuilder(collectionName)
	query, args := qb.BuildInsert(insertData, "created", "updated")

	stmt, err := r.stmtCache.Prepare(query)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "DB_PREPARE_ERROR", "Failed to prepare statement").WithDetails(map[string]any{"error": err.Error()})
	}

	err = stmt.QueryRowContext(ctx, args...).Scan(&record.Created, &record.Updated)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "RECORD_CREATE_FAILED", "Failed to create record").WithDetails(map[string]any{"error": err.Error()})
	}

	return record, nil
}

func (r *Repository) FindRecordByID(ctx context.Context, collectionName string, id string) (*models.Record, error) {
	col, ok := r.registry.GetCollection(collectionName)
	if !ok {
		return nil, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	columns := []string{"id", "created", "updated"}
	for _, f := range col.Fields {
		columns = append(columns, f.Name)
	}

	qb := NewQueryBuilder(collectionName)
	query, args := qb.Select(columns...).Where("id = ?", id).BuildSelect()

	stmt, err := r.stmtCache.Prepare(query)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "DB_PREPARE_ERROR", "Failed to prepare statement").WithDetails(map[string]any{"error": err.Error()})
	}

	row := stmt.QueryRowContext(ctx, args...)

	vals := make([]any, len(columns))
	valPtrs := make([]any, len(columns))
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	if err := row.Scan(valPtrs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewError(http.StatusNotFound, "RECORD_NOT_FOUND", "Record not found")
		}
		return nil, errors.NewError(http.StatusInternalServerError, "RECORD_FETCH_FAILED", "Failed to fetch record").WithDetails(map[string]any{"error": err.Error()})
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

func (r *Repository) UpdateRecord(ctx context.Context, collectionName string, id string, data map[string]any) (*models.Record, error) {
	col, ok := r.registry.GetCollection(collectionName)
	if !ok {
		return nil, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	record, err := r.FindRecordByID(ctx, collectionName, id)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		if k != "id" && k != "created" && k != "updated" {
			record.Data[k] = v
		}
	}

	updateData := make(map[string]any)
	updateData["updated"] = time.Now().UTC().Format(time.RFC3339)

	for _, f := range col.Fields {
		if val, ok := record.Data[f.Name]; ok {
			updateData[f.Name] = val
		}
	}

	qb := NewQueryBuilder(collectionName)
	query, args := qb.Where("id = ?", id).BuildUpdate(updateData, "updated")

	stmt, err := r.stmtCache.Prepare(query)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "DB_PREPARE_ERROR", "Failed to prepare statement").WithDetails(map[string]any{"error": err.Error()})
	}

	err = stmt.QueryRowContext(ctx, args...).Scan(&record.Updated)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "RECORD_UPDATE_FAILED", "Failed to update record").WithDetails(map[string]any{"error": err.Error()})
	}

	return record, nil
}

func (r *Repository) DeleteRecord(ctx context.Context, collectionName string, id string) error {
	_, ok := r.registry.GetCollection(collectionName)
	if !ok {
		return errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	qb := NewQueryBuilder(collectionName)
	query, args := qb.Where("id = ?", id).BuildDelete()

	stmt, err := r.stmtCache.Prepare(query)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_PREPARE_ERROR", "Failed to prepare statement").WithDetails(map[string]any{"error": err.Error()})
	}

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "RECORD_DELETE_FAILED", "Failed to delete record").WithDetails(map[string]any{"error": err.Error()})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NewError(http.StatusNotFound, "RECORD_NOT_FOUND", "Record not found")
	}

	return nil
}

func (r *Repository) ListRecords(ctx context.Context, collectionName string, params QueryParams) ([]*models.Record, int, error) {
	col, ok := r.registry.GetCollection(collectionName)
	if !ok {
		return nil, 0, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", fmt.Sprintf("Collection %s not found", collectionName))
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 30
	}

	qb := NewQueryBuilder(collectionName)

	if params.Filter != "" {
		clause, values, err := r.parseSafeFilter(col, params.Filter)
		if err != nil {
			return nil, 0, err
		}
		qb.Where(clause, values...)
	}

	// Validate and apply sort
	sortField, sortDir, err := r.validateSortField(col, params.Sort)
	if err != nil {
		return nil, 0, errors.NewError(http.StatusBadRequest, "INVALID_SORT", err.Error())
	}
	qb.OrderBy(fmt.Sprintf("%s %s", sortField, sortDir))

	// Count total
	countQuery, countArgs := qb.BuildCount()
	var total int

	stmt, err := r.stmtCache.Prepare(countQuery)
	if err != nil {
		return nil, 0, errors.NewError(http.StatusInternalServerError, "DB_PREPARE_ERROR", "Failed to prepare count statement").WithDetails(map[string]any{"error": err.Error()})
	}

	err = stmt.QueryRowContext(ctx, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.NewError(http.StatusInternalServerError, "DB_COUNT_FAILED", "Failed to count records").WithDetails(map[string]any{"error": err.Error()})
	}

	columns := []string{"id", "created", "updated"}
	for _, f := range col.Fields {
		columns = append(columns, f.Name)
	}

	offset := (params.Page - 1) * params.PerPage
	query, args := qb.Select(columns...).Limit(params.PerPage).Offset(offset).BuildSelect()

	stmt, err = r.stmtCache.Prepare(query)
	if err != nil {
		return nil, 0, errors.NewError(http.StatusInternalServerError, "DB_PREPARE_ERROR", "Failed to prepare list statement").WithDetails(map[string]any{"error": err.Error()})
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, errors.NewError(http.StatusInternalServerError, "RECORD_LIST_FAILED", "Failed to list records").WithDetails(map[string]any{"error": err.Error()})
	}
	defer errors.Defer(ctx, rows.Close, "close rows")

	var records []*models.Record
	for rows.Next() {
		vals := make([]any, len(columns))
		valPtrs := make([]any, len(columns))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}

		if err := rows.Scan(valPtrs...); err != nil {
			return nil, 0, errors.NewError(http.StatusInternalServerError, "RECORD_SCAN_FAILED", "Failed to scan record").WithDetails(map[string]any{"error": err.Error()})
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
		r.expandRecords(ctx, col, records, params.Expand)
	}

	return records, total, nil
}

func (r *Repository) validateSortField(col *models.Collection, sortParam string) (string, string, error) {
	if sortParam == "" {
		return "created", "DESC", nil // Default sort
	}

	direction := "ASC"
	fieldName := sortParam
	if strings.HasPrefix(sortParam, "-") {
		direction = "DESC"
		fieldName = strings.TrimPrefix(sortParam, "-")
	}

	// Validate field name
	isValid := false
	if fieldName == "id" || fieldName == "created" || fieldName == "updated" {
		isValid = true
	} else {
		for _, f := range col.Fields {
			if f.Name == fieldName {
				isValid = true
				break
			}
		}
	}

	if !isValid {
		return "", "", fmt.Errorf("invalid sort field: %s", fieldName)
	}

	return fieldName, direction, nil
}

// parseSafeFilter implements basic parameterized filtering to prevent SQL injection
func (r *Repository) parseSafeFilter(col *models.Collection, filter string) (string, []any, error) {
	// Simple support for: field = 'value' or field != 'value'
	operators := []string{" != ", " = ", " > ", " < ", " >= ", " <= "}

	for _, op := range operators {
		if strings.Contains(filter, op) {
			parts := strings.Split(filter, op)
			if len(parts) != 2 {
				continue
			}

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
				return "", nil, errors.NewError(http.StatusBadRequest, "INVALID_FILTER", fmt.Sprintf("Unknown field in filter: %s", fieldName))
			}

			// 2. Clean value (remove single quotes if present)
			value = strings.Trim(value, "'")

			// 3. Return parameterized clause
			return fmt.Sprintf("%s %s ?", fieldName, strings.TrimSpace(op)), []any{value}, nil
		}
	}

	return "", nil, errors.NewError(http.StatusBadRequest, "UNSUPPORTED_FILTER", "Filter format not supported. Use 'field = value'")
}

func (r *Repository) expandRecords(ctx context.Context, collection *models.Collection, records []*models.Record, expand string) {
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
		if relField == nil {
			continue
		}

		targetCol := ""
		if options, ok := relField.Options.(map[string]any); ok {
			targetCol, ok = options["collection"].(string)
			if !ok {
				errors.Log(ctx, nil, "relation field missing collection option", "field", relField.Name, "collection", collection.Name)
			}
		}
		if targetCol == "" {
			continue
		}

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
		if len(relIDs) == 0 {
			continue
		}

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
		targetRecords, _, err := r.ListRecords(ctx, targetCol, QueryParams{
			Filter:  filter,
			PerPage: len(relIDs),
		})
		if err != nil {
			continue
		}

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
