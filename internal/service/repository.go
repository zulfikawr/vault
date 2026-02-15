package service

import (
	"context"

	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/models"
)

type Repository interface {
	CreateRecord(ctx context.Context, collectionName string, data map[string]any) (*models.Record, error)
	ListRecords(ctx context.Context, collectionName string, params db.QueryParams) ([]*models.Record, int, error)
	FindRecordByID(ctx context.Context, collectionName string, id string) (*models.Record, error)
	UpdateRecord(ctx context.Context, collectionName string, id string, data map[string]any) (*models.Record, error)
	DeleteRecord(ctx context.Context, collectionName string, id string) error
	Close()
}
