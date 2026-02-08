package db

import (
	"database/sql"
	"sync"

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

func (s *SchemaRegistry) AddCollection(c *models.Collection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collections[c.Name] = c
}

func (s *SchemaRegistry) LoadFromDB() error {
	// We'll implement the logic to load from _collections table later
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
