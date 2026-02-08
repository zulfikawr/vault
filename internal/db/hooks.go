package db

import (
	"context"

	"github.com/zulfikawr/vault/internal/models"
)

type HookFunc func(ctx context.Context, record *models.Record) error

type Hooks struct {
	BeforeCreate []HookFunc
	AfterCreate  []HookFunc
	BeforeUpdate []HookFunc
	AfterUpdate  []HookFunc
	BeforeDelete []HookFunc
	AfterDelete  []HookFunc
}

var globalHooks = make(map[string]*Hooks)

func GetHooks(collection string) *Hooks {
	if h, ok := globalHooks[collection]; ok {
		return h
	}
	h := &Hooks{}
	globalHooks[collection] = h
	return h
}

func (h *Hooks) TriggerBeforeCreate(ctx context.Context, record *models.Record) error {
	for _, fn := range h.BeforeCreate {
		if err := fn(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hooks) TriggerAfterCreate(ctx context.Context, record *models.Record) error {
	for _, fn := range h.AfterCreate {
		if err := fn(ctx, record); err != nil {
			return err
		}
	}
	return nil
}