package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/zulfikawr/vault/internal/core"
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
	return list
}

func (s *SchemaRegistry) AddCollection(c *models.Collection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collections[c.Name] = c
}

func (s *SchemaRegistry) LoadFromDB(ctx context.Context) error {
	// First, ensure all existing records have IDs (migration for old data)
	_, err := s.db.ExecContext(ctx, `UPDATE _collections SET id = 'col_' || name WHERE id IS NULL`)
	if err != nil {
		slog.Warn("Failed to migrate _collections IDs", "error", err)
	}

	rows, err := s.db.QueryContext(ctx, "SELECT id, name, type, fields, created, updated FROM _collections")
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var id, name, ctype, fieldsJSON, created, updated string
		if err := rows.Scan(&id, &name, &ctype, &fieldsJSON, &created, &updated); err != nil {
			return err
		}

		var fields []models.Field
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			return err
		}

		s.AddCollection(&models.Collection{
			ID:      id,
			Name:    name,
			Type:    models.CollectionType(ctype),
			Fields:  fields,
			Created: created,
			Updated: updated,
		})
	}
	return nil
}

func (s *SchemaRegistry) SaveCollection(ctx context.Context, c *models.Collection) error {
	fieldsJSON, _ := json.Marshal(c.Fields)

	// Generate ID if not present
	if c.ID == "" {
		c.ID = "col_" + c.Name
	}

	query := `INSERT INTO _collections (id, name, type, fields) VALUES (?, ?, ?, ?) 
	          ON CONFLICT(name) DO UPDATE SET type=excluded.type, fields=excluded.fields`

	_, err := s.db.ExecContext(ctx, query, c.ID, c.Name, c.Type, string(fieldsJSON))
	if err != nil {
		return core.NewError(http.StatusInternalServerError, "DB_SAVE_COLLECTION_FAILED", "Failed to persist collection definition").WithDetails(map[string]any{"error": err.Error()})
	}

	s.AddCollection(c)
	return nil
}

// BootstrapSystemCollections initializes the internal meta tables
func (s *SchemaRegistry) BootstrapSystemCollections() error {
	collectionsTable := &models.Collection{
		ID:   "system_collections",
		Name: "_collections",
		Type: models.CollectionTypeSystem,
		Fields: []models.Field{
			{Name: "name", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "type", Type: models.FieldTypeText, Required: true},
			{Name: "fields", Type: models.FieldTypeJSON, Required: true},
			{Name: "listRule", Type: models.FieldTypeText},
			{Name: "viewRule", Type: models.FieldTypeText},
			{Name: "createRule", Type: models.FieldTypeText},
			{Name: "updateRule", Type: models.FieldTypeText},
			{Name: "deleteRule", Type: models.FieldTypeText},
		},
		Created: time.Now().Format(time.RFC3339),
		Updated: time.Now().Format(time.RFC3339),
	}

	s.AddCollection(collectionsTable)
	return nil
}

func (s *SchemaRegistry) BootstrapUsersCollection() error {
	usersTable := &models.Collection{
		ID:   "system_users",
		Name: "users",
		Type: models.CollectionTypeAuth,
		Fields: []models.Field{
			{Name: "username", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "email", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "password", Type: models.FieldTypeText, Required: true},
			{Name: "verified", Type: models.FieldTypeBool},
			{Name: "lastLogin", Type: models.FieldTypeDate},
		},
	}

	s.AddCollection(usersTable)
	return nil
}

func (s *SchemaRegistry) BootstrapRefreshTokensCollection() error {
	tokensTable := &models.Collection{
		ID:   "system_refresh_tokens",
		Name: "_refresh_tokens",
		Type: models.CollectionTypeSystem,
		Fields: []models.Field{
			{Name: "token", Type: models.FieldTypeText, Required: true, Unique: true},
			{Name: "user_id", Type: models.FieldTypeText, Required: true},
			{Name: "expires", Type: models.FieldTypeDate, Required: true},
		},
	}

	s.AddCollection(tokensTable)
	return nil
}
