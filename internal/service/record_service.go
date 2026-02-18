package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
	"github.com/zulfikawr/vault/internal/realtime"
)

type RecordService struct {
	repo Repository
	hub  *realtime.Hub
}

func NewRecordService(repo Repository, hub *realtime.Hub) *RecordService {
	return &RecordService{repo: repo, hub: hub}
}

func (s *RecordService) Close() {
	s.repo.Close()
}

func (s *RecordService) broadcast(action string, collection string, record *models.Record) {
	if s.hub != nil {
		s.hub.Broadcast(&realtime.Message{
			Action:     action,
			Collection: collection,
			Record:     record,
		})
	}
}

func (s *RecordService) CreateRecord(ctx context.Context, collectionName string, data map[string]any) (*models.Record, error) {
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

	// Use record.Data which might have been modified by hooks
	createdRecord, err := s.repo.CreateRecord(ctx, collectionName, record.Data)
	if err != nil {
		return nil, err
	}

	if err := hooks.TriggerAfterCreate(ctx, createdRecord); err != nil {
		errors.Log(ctx, err, "after create hook failed", "collection", collectionName, "record_id", createdRecord.ID)
	}

	s.broadcast("create", collectionName, createdRecord)
	return createdRecord, nil
}

func (s *RecordService) ListRecords(ctx context.Context, collectionName string, params db.QueryParams) ([]*models.Record, int, error) {
	return s.repo.ListRecords(ctx, collectionName, params)
}

func (s *RecordService) FindRecordByID(ctx context.Context, collectionName string, id string) (*models.Record, error) {
	return s.repo.FindRecordByID(ctx, collectionName, id)
}

func (s *RecordService) UpdateRecord(ctx context.Context, collectionName string, id string, data map[string]any) (*models.Record, error) {
	record, err := s.repo.FindRecordByID(ctx, collectionName, id)
	if err != nil {
		return nil, err
	}

	// Merge incoming data into existing record data
	for k, v := range data {
		if k != "id" && k != "created" && k != "updated" {
			record.Data[k] = v
		}
	}

	hooks := GetHooks(collectionName)
	if err := hooks.TriggerBeforeUpdate(ctx, record); err != nil {
		return nil, err
	}

	// Use record.Data which now contains merged data and potential hook modifications.
	// The repository will handle filtering fields that are no longer in the schema.
	updatedRecord, err := s.repo.UpdateRecord(ctx, collectionName, id, record.Data)
	if err != nil {
		return nil, err
	}

	if err := hooks.TriggerAfterUpdate(ctx, updatedRecord); err != nil {
		errors.Log(ctx, err, "after update hook failed", "collection", collectionName, "record_id", updatedRecord.ID)
	}

	s.broadcast("update", collectionName, updatedRecord)
	return updatedRecord, nil
}

func (s *RecordService) DeleteRecord(ctx context.Context, collectionName string, id string) error {
	record, err := s.repo.FindRecordByID(ctx, collectionName, id)
	if err != nil {
		return err
	}

	hooks := GetHooks(collectionName)
	if err := hooks.TriggerBeforeDelete(ctx, record); err != nil {
		return err
	}

	if err := s.repo.DeleteRecord(ctx, collectionName, id); err != nil {
		return err
	}

	if err := hooks.TriggerAfterDelete(ctx, record); err != nil {
		errors.Log(ctx, err, "after delete hook failed", "collection", collectionName, "record_id", record.ID)
	}

	s.broadcast("delete", collectionName, record)
	return nil
}
