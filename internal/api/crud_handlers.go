package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/rules"
	"github.com/zulfikawr/vault/internal/service"
)

type CollectionHandler struct {
	recordService *service.RecordService
	registry      *db.SchemaRegistry
}

func NewCollectionHandler(recordService *service.RecordService, registry *db.SchemaRegistry) *CollectionHandler {
	return &CollectionHandler{
		recordService: recordService,
		registry:      registry,
	}
}

func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	// For Phase 5, if listRule is nil/empty, it's public.
	// If it's set, we check basic auth for now.
	// (Step 3 will implement full SQL-level filtering)
	if col.ListRule != nil && *col.ListRule != "" {
		evalCtx := service.GetEvaluationContext(r, nil)
		allowed, err := rules.Evaluate(*col.ListRule, evalCtx)
		if !allowed || err != nil {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "FORBIDDEN", "You do not have permission to list this collection"))
			return
		}
	}

	params := h.parseQueryParams(r)
	records, total, err := h.recordService.ListRecords(r.Context(), collectionName, params)
	if err != nil {
		errors.SendError(w, err)
		return
	}

	for _, record := range records {
		if collectionName == "users" {
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
	collectionName := r.PathValue("collection")
	id := r.PathValue("id")

	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	record, err := h.recordService.FindRecordByID(r.Context(), collectionName, id)
	if err != nil {
		errors.SendError(w, err)
		return
	}

	// Rule Check
	if col.ViewRule != nil && *col.ViewRule != "" {
		evalCtx := service.GetEvaluationContext(r, record.Data)
		allowed, err := rules.Evaluate(*col.ViewRule, evalCtx)
		if !allowed || err != nil {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this record"))
			return
		}
	}

	if collectionName == "users" {
		record.HideField("password")
	}

	SendJSON(w, http.StatusOK, record, nil)
}

func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")

	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	// Rule Check (Pre-create check)
	if col.CreateRule != nil && *col.CreateRule != "" {
		evalCtx := service.GetEvaluationContext(r, nil)
		evalCtx.Data = data // Inject incoming data
		allowed, err := rules.Evaluate(*col.CreateRule, evalCtx)
		if !allowed || err != nil {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "FORBIDDEN", "You do not have permission to create records in this collection"))
			return
		}
	}

	if err := service.ValidateRecord(col, data); err != nil {
		errors.SendError(w, err)
		return
	}

	record, err := h.recordService.CreateRecord(r.Context(), collectionName, data)
	if err != nil {
		errors.SendError(w, err)
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

	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	// Fetch current for rule evaluation
	existing, err := h.recordService.FindRecordByID(r.Context(), collectionName, id)
	if err != nil {
		errors.SendError(w, err)
		return
	}

	// Rule Check
	if col.UpdateRule != nil && *col.UpdateRule != "" {
		evalCtx := service.GetEvaluationContext(r, existing.Data)
		evalCtx.Data = data
		allowed, err := rules.Evaluate(*col.UpdateRule, evalCtx)
		if !allowed || err != nil {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "FORBIDDEN", "You do not have permission to update this record"))
			return
		}
	}

	record, err := h.recordService.UpdateRecord(r.Context(), collectionName, id, data)
	if err != nil {
		errors.SendError(w, err)
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

	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	// Fetch current for rule evaluation
	existing, err := h.recordService.FindRecordByID(r.Context(), collectionName, id)
	if err != nil {
		errors.SendError(w, err)
		return
	}

	// Rule Check
	if col.DeleteRule != nil && *col.DeleteRule != "" {
		evalCtx := service.GetEvaluationContext(r, existing.Data)
		allowed, err := rules.Evaluate(*col.DeleteRule, evalCtx)
		if !allowed || err != nil {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "FORBIDDEN", "You do not have permission to delete this record"))
			return
		}
	}

	if err := h.recordService.DeleteRecord(r.Context(), collectionName, id); err != nil {
		errors.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	col, ok := h.registry.GetCollection(collectionName)
	if !ok {
		errors.SendError(w, errors.NewError(http.StatusNotFound, "COLLECTION_NOT_FOUND", "Collection not found"))
		return
	}

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_BODY", "Failed to decode request body"))
		return
	}

	if len(req.IDs) == 0 {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "MISSING_IDS", "At least one ID is required"))
		return
	}

	for _, id := range req.IDs {
		// Fetch current for rule evaluation
		existing, err := h.recordService.FindRecordByID(r.Context(), collectionName, id)
		if err != nil {
			continue // Skip if not found
		}

		// Rule Check
		if col.DeleteRule != nil && *col.DeleteRule != "" {
			evalCtx := service.GetEvaluationContext(r, existing.Data)
			allowed, err := rules.Evaluate(*col.DeleteRule, evalCtx)
			if !allowed || err != nil {
				continue // Skip if forbidden
			}
		}

		_ = h.recordService.DeleteRecord(r.Context(), collectionName, id)
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
