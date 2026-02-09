package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/storage"
)

type FileHandler struct {
	storage  storage.Storage
	executor *db.Executor
}

func NewFileHandler(s storage.Storage, e *db.Executor) *FileHandler {
	return &FileHandler{storage: s, executor: e}
}

func (h *FileHandler) Serve(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	recordID := r.PathValue("id")
	filename := r.PathValue("filename")

	path := filepath.Join(collection, recordID, filename)
	
	file, err := h.storage.Retrieve(r.Context(), path)
	if err != nil {
		core.SendError(w, err)
		return
	}
	defer file.Close()

	// Simple content type detection
	contentType := "application/octet-stream"
	if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".gif") {
		contentType = "image/gif"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	io.Copy(w, file)
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// For this prototype, we'll implement a standalone upload that returns the filename/metadata.
	// Later we can integrate it directly into the Record Create/Update handlers.
	
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_MULTIPART", "Failed to parse multipart form"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "FILE_REQUIRED", "No file provided"))
		return
	}
	defer file.Close()

	collection := r.FormValue("collection")
	recordID := r.FormValue("recordID")
	if collection == "" || recordID == "" {
		core.SendError(w, core.NewError(http.StatusBadRequest, "MISSING_PARAMS", "collection and recordID are required"))
		return
	}

	// Generate a safe unique name
	ext := filepath.Ext(header.Filename)
	safeName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	path := filepath.Join(collection, recordID, safeName)

	if err := h.storage.Save(r.Context(), path, file); err != nil {
		core.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusCreated, map[string]any{
		"name": safeName,
		"size": header.Size,
		"mime": header.Header.Get("Content-Type"),
		"url":  fmt.Sprintf("/api/files/%s/%s/%s", collection, recordID, safeName),
	}, nil)
}
