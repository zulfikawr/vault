package service

import (
	"context"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/models"
)

func RegisterAuthHooks() {
	hooks := GetHooks("users")

	// Avoid adding duplicate hooks if called multiple times
	// This is a simple check, ideally we'd have named hooks or better management
	if len(hooks.BeforeCreate) == 0 {
		hooks.BeforeCreate = append(hooks.BeforeCreate, hashPasswordHook)
	}
	if len(hooks.BeforeUpdate) == 0 {
		hooks.BeforeUpdate = append(hooks.BeforeUpdate, hashPasswordHook)
	}
}

func hashPasswordHook(ctx context.Context, record *models.Record) error {
	password := record.GetString("password")
	if password == "" {
		return nil
	}

	hashed, err := auth.HashPassword(ctx, password)
	if err != nil {
		return err
	}
	record.Data["password"] = hashed
	return nil
}
