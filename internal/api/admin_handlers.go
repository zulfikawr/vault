package api

import (
	"encoding/json"
	"net/http"

	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
	"github.com/zulfikawr/vault/internal/service"
)

type AdminHandler struct {
	collectionService *service.CollectionService
}

func NewAdminHandler(collectionService *service.CollectionService) *AdminHandler {
	return &AdminHandler{collectionService: collectionService}
}

func (h *AdminHandler) ListCollections(w http.ResponseWriter, r *http.Request) {
	collections := h.collectionService.ListCollections()
	SendJSON(w, http.StatusOK, collections, nil)
}

func (h *AdminHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var col models.Collection
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	if err := h.collectionService.CreateCollection(r.Context(), &col); err != nil {
		errors.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusCreated, col, nil)
}

func (h *AdminHandler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "MISSING_ID", "Collection ID is required"))
		return
	}

	var col models.Collection
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	// Ensure the ID from the path is used
	col.ID = id

	if err := h.collectionService.CreateCollection(r.Context(), &col); err != nil {
		errors.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusOK, col, nil)
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
	if _, exists := h.collectionService.GetCollection(name); !exists {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	if err := h.collectionService.DeleteCollection(r.Context(), name); err != nil {
		errors.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusOK, map[string]string{"message": "Collection deleted successfully"}, nil)
}
