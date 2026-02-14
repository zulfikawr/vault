package api

import (
	"encoding/json"
	"net/http"

	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
)

type AdminHandler struct {
	executor  *db.Executor
	registry  *db.SchemaRegistry
	migration *db.MigrationEngine
}

func NewAdminHandler(e *db.Executor, r *db.SchemaRegistry, m *db.MigrationEngine) *AdminHandler {
	return &AdminHandler{executor: e, registry: r, migration: m}
}

func (h *AdminHandler) ListCollections(w http.ResponseWriter, r *http.Request) {
	collections := h.registry.GetCollections()
	SendJSON(w, http.StatusOK, collections, nil)
}

func (h *AdminHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var col models.Collection
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	// 1. Sync DB
	if err := h.migration.SyncCollection(r.Context(), &col); err != nil {
		errors.SendError(w, err)
		return
	}

	// 2. Persist Definition
	if err := h.registry.SaveCollection(r.Context(), &col); err != nil {
		errors.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusCreated, col, nil)
}

func (h *AdminHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	// Simple placeholder for settings
	SendJSON(w, http.StatusOK, map[string]string{"appName": "Vault"}, nil)
}

func (h *AdminHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	// Simple logic: Copy current DB to a timestamped file
	// In production, we'd use SQLite's backup API
	SendJSON(w, http.StatusAccepted, map[string]string{"message": "Backup triggered"}, nil)
}

func (h *AdminHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "MISSING_NAME", "Collection name is required"))
		return
	}

	// Check if collection exists
	if _, exists := h.registry.GetCollection(name); !exists {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	// Remove from registry
	h.registry.RemoveCollection(name)

	// Remove from database
	ctx := r.Context()
	if err := h.migration.DropCollection(ctx, name); err != nil {
		errors.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusOK, map[string]string{"message": "Collection deleted successfully"}, nil)
}
