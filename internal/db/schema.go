package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
)

type SchemaRegistry struct {
	mu          sync.RWMutex
	collections map[string]*models.Collection
	db          *sql.DB
}

func NewSchemaRegistry(db *sql.DB) *SchemaRegistry {
	return &SchemaRegistry{
		collections: make(map[string]*models.Collection),
		db:          db,
	}
}

func (s *SchemaRegistry) GetCollection(name string) (*models.Collection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.collections[name]
	return c, ok
}

func (s *SchemaRegistry) GetCollections() []*models.Collection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Collection, 0, len(s.collections))
	for _, c := range s.collections {
		list = append(list, c)
	}

	// Sort collections by Created field in descending order (most recent first)
	sort.Slice(list, func(i, j int) bool {
		// Parse the created timestamps to compare them
		timeI, errI := time.Parse(time.RFC3339, list[i].Created)
		timeJ, errJ := time.Parse(time.RFC3339, list[j].Created)

		// If parsing fails, fallback to string comparison
		if errI != nil || errJ != nil {
			return list[i].Created > list[j].Created
		}

		// Compare times (descending order - most recent first)
		return timeI.After(timeJ)
	})

	return list
}

func (s *SchemaRegistry) AddCollection(c *models.Collection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collections[c.Name] = c
}

func (s *SchemaRegistry) RemoveCollection(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.collections, name)
}

func (s *SchemaRegistry) LoadFromDB(ctx context.Context) error {
	// First, ensure all existing records have IDs (migration for old data)
	_, err := s.db.ExecContext(ctx, `UPDATE _collections SET id = 'col_' || name WHERE id IS NULL`)
	if err != nil {
		slog.Warn("Failed to migrate _collections IDs", "error", err)
	}

	rows, err := s.db.QueryContext(ctx, "SELECT id, name, type, fields, list_rule, view_rule, create_rule, update_rule, delete_rule, created, updated FROM _collections")
	if err != nil {
		return err
	}
	defer errors.Defer(ctx, rows.Close, "close rows")

	for rows.Next() {
		var id, name, ctype, fieldsJSON, created, updated string
		var listRule, viewRule, createRule, updateRule, deleteRule sql.NullString
		if err := rows.Scan(&id, &name, &ctype, &fieldsJSON, &listRule, &viewRule, &createRule, &updateRule, &deleteRule, &created, &updated); err != nil {
			return err
		}

		var fields []models.Field
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			return err
		}

		col := &models.Collection{
			ID:      id,
			Name:    name,
			Type:    models.CollectionType(ctype),
			Fields:  fields,
			Created: created,
			Updated: updated,
		}

		if listRule.Valid {
			col.ListRule = &listRule.String
		}
		if viewRule.Valid {
			col.ViewRule = &viewRule.String
		}
		if createRule.Valid {
			col.CreateRule = &createRule.String
		}
		if updateRule.Valid {
			col.UpdateRule = &updateRule.String
		}
		if deleteRule.Valid {
			col.DeleteRule = &deleteRule.String
		}

		s.AddCollection(col)
	}
	return nil
}

func (s *SchemaRegistry) SaveCollection(ctx context.Context, c *models.Collection) error {
	fieldsJSON, _ := json.Marshal(c.Fields)

	// Generate ID if not present
	if c.ID == "" {
		c.ID = "col_" + c.Name
	}

	query := `INSERT INTO _collections (id, name, type, fields, list_rule, view_rule, create_rule, update_rule, delete_rule) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) 
	          ON CONFLICT(name) DO UPDATE SET 
			  type=excluded.type, 
			  fields=excluded.fields,
			  list_rule=excluded.list_rule,
			  view_rule=excluded.view_rule,
			  create_rule=excluded.create_rule,
			  update_rule=excluded.update_rule,
			  delete_rule=excluded.delete_rule`

	_, err := s.db.ExecContext(ctx, query,
		c.ID, c.Name, c.Type, string(fieldsJSON),
		c.ListRule, c.ViewRule, c.CreateRule, c.UpdateRule, c.DeleteRule,
	)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_SAVE_COLLECTION_FAILED", "Failed to persist collection definition").WithDetails(map[string]any{"error": err.Error()})
	}

	s.AddCollection(c)
	return nil
}

// BootstrapSystemCollections initializes the internal meta tables
func (s *SchemaRegistry) BootstrapSystemCollections() error {
	adminOnly := "@request.auth.id != ''"
	collectionsTable := &models.Collection{
		ID:   "system_collections",
		Name: "_collections",
		Type: models.CollectionTypeSystem,
		Fields: []models.Field{
			{Name: "name", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "type", Type: models.FieldTypeText, Required: true},
			{Name: "fields", Type: models.FieldTypeJSON, Required: true},
			{Name: "list_rule", Type: models.FieldTypeText},
			{Name: "view_rule", Type: models.FieldTypeText},
			{Name: "create_rule", Type: models.FieldTypeText},
			{Name: "update_rule", Type: models.FieldTypeText},
			{Name: "delete_rule", Type: models.FieldTypeText},
		},
		ListRule:   &adminOnly,
		ViewRule:   &adminOnly,
		CreateRule: &adminOnly,
		UpdateRule: &adminOnly,
		DeleteRule: &adminOnly,
		Created:    time.Now().Format(time.RFC3339),
		Updated:    time.Now().Format(time.RFC3339),
	}

	s.AddCollection(collectionsTable)
	return nil
}

func (s *SchemaRegistry) BootstrapUsersCollection() error {
	adminOnly := "@request.auth.id != ''"
	usersTable := &models.Collection{
		ID:   "system_users",
		Name: "users",
		Type: models.CollectionTypeAuth,
		Fields: []models.Field{
			{Name: "username", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "email", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "password", Type: models.FieldTypeText, Required: true},
			{Name: "lastLogin", Type: models.FieldTypeDate},
		},
		ListRule:   &adminOnly,
		ViewRule:   &adminOnly,
		CreateRule: &adminOnly,
		UpdateRule: &adminOnly,
		DeleteRule: &adminOnly,
	}

	s.AddCollection(usersTable)
	return nil
}

func (s *SchemaRegistry) BootstrapRefreshTokensCollection() error {
	adminOnly := "@request.auth.id != ''"
	tokensTable := &models.Collection{
		ID:   "system_refresh_tokens",
		Name: "_refresh_tokens",
		Type: models.CollectionTypeSystem,
		Fields: []models.Field{
			{Name: "token", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "user_id", Type: models.FieldTypeText, Required: true},
			{Name: "expires", Type: models.FieldTypeDate, Required: true},
		},
		ListRule:   &adminOnly,
		ViewRule:   &adminOnly,
		CreateRule: &adminOnly,
		UpdateRule: &adminOnly,
		DeleteRule: &adminOnly,
	}

	s.AddCollection(tokensTable)
	return nil
}

func (s *SchemaRegistry) BootstrapAuditLogsCollection() error {
	adminOnly := "@request.auth.id != ''"
	auditTable := &models.Collection{
		ID:   "system_audit_logs",
		Name: "_audit_logs",
		Type: models.CollectionTypeSystem,
		Fields: []models.Field{
			{Name: "action", Type: models.FieldTypeText, Required: true},
			{Name: "resource", Type: models.FieldTypeText, Required: true},
			{Name: "admin_id", Type: models.FieldTypeText, Required: true},
			{Name: "details", Type: models.FieldTypeJSON},
			{Name: "timestamp", Type: models.FieldTypeDate, Required: true},
		},
		ListRule:   &adminOnly,
		ViewRule:   &adminOnly,
		CreateRule: &adminOnly,
		UpdateRule: &adminOnly,
		DeleteRule: &adminOnly,
		Created:    time.Now().Format(time.RFC3339),
		Updated:    time.Now().Format(time.RFC3339),
	}

	s.AddCollection(auditTable)
	return nil
}
