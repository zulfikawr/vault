package db

import "github.com/zulfikawr/vault/internal/models"

type HookFunc func(record *models.Record) error

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

func (h *Hooks) TriggerBeforeCreate(record *models.Record) error {
	for _, fn := range h.BeforeCreate {
		if err := fn(record); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hooks) TriggerAfterCreate(record *models.Record) error {
	for _, fn := range h.AfterCreate {
		if err := fn(record); err != nil {
			return err
		}
	}
	return nil
}
