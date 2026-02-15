package service

import (
	"context"
	"fmt"

	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/models"
)

type CollectionService struct {
	registry  *db.SchemaRegistry
	migration *db.MigrationEngine
}

func NewCollectionService(registry *db.SchemaRegistry, migration *db.MigrationEngine) *CollectionService {
	return &CollectionService{
		registry:  registry,
		migration: migration,
	}
}

func (s *CollectionService) InitSystem(ctx context.Context) error {
	if err := s.registry.BootstrapSystemCollections(); err != nil {
		return fmt.Errorf("failed to bootstrap system collections: %w", err)
	}
	if err := s.registry.BootstrapRefreshTokensCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap refresh tokens collection: %w", err)
	}
	if err := s.registry.BootstrapUsersCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap users collection: %w", err)
	}
	if err := s.registry.BootstrapAuditLogsCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap audit logs collection: %w", err)
	}

	systemCols := []string{"_collections", "_refresh_tokens", "_audit_logs", "users"}
	for _, name := range systemCols {
		col, ok := s.registry.GetCollection(name)
		if !ok || col == nil {
			return fmt.Errorf("system collection %s not found", name)
		}
		if err := s.migration.SyncCollection(ctx, col); err != nil {
			return fmt.Errorf("failed to sync collection %s: %w", name, err)
		}
		if name != "_collections" {
			if err := s.registry.SaveCollection(ctx, col); err != nil {
				return fmt.Errorf("failed to save collection %s: %w", name, err)
			}
		}
	}
	return nil
}

func (s *CollectionService) ListCollections() []*models.Collection {
	return s.registry.GetCollections()
}

func (s *CollectionService) GetCollection(name string) (*models.Collection, bool) {
	return s.registry.GetCollection(name)
}

func (s *CollectionService) CreateCollection(ctx context.Context, col *models.Collection) error {
	// 1. Sync DB
	if err := s.migration.SyncCollection(ctx, col); err != nil {
		return err
	}

	// 2. Persist Definition
	if err := s.registry.SaveCollection(ctx, col); err != nil {
		return err
	}

	return nil
}

func (s *CollectionService) DeleteCollection(ctx context.Context, name string) error {
	// Remove from registry
	s.registry.RemoveCollection(name)

	// Remove from database
	if err := s.migration.DropCollection(ctx, name); err != nil {
		return err
	}

	return nil
}
