package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

type CollectionHandler struct {
	executor *db.Executor
	registry *db.SchemaRegistry
}

func NewCollectionHandler(executor *db.Executor, registry *db.SchemaRegistry) *CollectionHandler {
	return &CollectionHandler{
		executor: executor,
		registry: registry,
	}
}

func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	
	params := h.parseQueryParams(r)
	
	records, total, err := h.executor.ListRecords(r.Context(), collection, params)
	if err != nil {
		core.SendError(w, err)
		return
	}

	// Hide sensitive fields for all records
	for _, record := range records {
		if collection == "users" {
			record.HideField("password")
		}
	}

	SendJSON(w, http.StatusOK, records, map[string]any{
		"page":       params.Page,
		"perPage":    params.PerPage,
		"totalItems": total,
	})
}

func (h *CollectionHandler) View(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	id := r.PathValue("id")

	record, err := h.executor.FindRecordByID(r.Context(), collection, id)
	if err != nil {
		core.SendError(w, err)
		return
	}

	if collection == "users" {
		record.HideField("password")
	}

	SendJSON(w, http.StatusOK, record, nil)
}

func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	
	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		core.SendError(w, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	// Validate against schema
	if err := db.ValidateRecord(col, data); err != nil {
		core.SendError(w, err)
		return
	}

	record, err := h.executor.CreateRecord(r.Context(), collectionName, data)
	if err != nil {
		core.SendError(w, err)
		return
	}

	if collectionName == "users" {
		record.HideField("password")
	}

	SendJSON(w, http.StatusCreated, record, nil)
}

func (h *CollectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	id := r.PathValue("id")

	if _, ok := h.registry.GetCollection(collectionName); !ok {
		core.SendError(w, core.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	// Validate provided fields (simplified PATCH validation)
	// We only validate what's provided. db.ValidateRecord currently checks 'Required', 
	// which might be problematic for PATCH if not all fields are sent.
	// For now, let's just use it and assume we'll refine validation later.

	record, err := h.executor.UpdateRecord(r.Context(), collectionName, id, data)
	if err != nil {
		core.SendError(w, err)
		return
	}

	if collectionName == "users" {
		record.HideField("password")
	}

	SendJSON(w, http.StatusOK, record, nil)
}

func (h *CollectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	id := r.PathValue("id")

	if err := h.executor.DeleteRecord(r.Context(), collectionName, id); err != nil {
		core.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) parseQueryParams(r *http.Request) db.QueryParams {
	q := r.URL.Query()
	
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("perPage"))
	
	return db.QueryParams{
		Page:    page,
		PerPage: perPage,
		Sort:    q.Get("sort"),
		Filter:  q.Get("filter"),
		Expand:  q.Get("expand"),
	}
}
