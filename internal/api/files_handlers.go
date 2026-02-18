package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/storage"
)

type FileHandler struct {
	storage       storage.Storage
	maxUploadSize int64
}

func NewFileHandler(s storage.Storage, maxUploadSize int64) *FileHandler {
	return &FileHandler{
		storage:       s,
		maxUploadSize: maxUploadSize,
	}
}

func (h *FileHandler) Serve(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	recordID := r.PathValue("id")
	filename := r.PathValue("filename")

	path := filepath.Join(collection, recordID, filename)

	file, err := h.storage.Retrieve(r.Context(), path)
	if err != nil {
		errors.SendError(w, err)
		return
	}
	defer errors.Defer(r.Context(), file.Close, "close file", "path", path)

	// Read first 512 bytes for content type detection
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	contentType := http.DetectContentType(buffer[:n])

	// Reset file pointer
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			errors.Log(r.Context(), err, "seek file", "path", path)
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	if _, err := io.Copy(w, file); err != nil {
		errors.Log(r.Context(), err, "copy file to response", "collection", collection, "file", filename)
	}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Enforce max upload size
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadSize)
	if err := r.ParseMultipartForm(h.maxUploadSize); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "FILE_TOO_LARGE", "File exceeds maximum allowed size"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "FILE_REQUIRED", "No file provided"))
		return
	}
	defer errors.Defer(r.Context(), file.Close, "close uploaded file", "filename", header.Filename)

	// Validate content type
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "FILE_READ_ERROR", "Failed to read file header"))
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "FILE_SEEK_ERROR", "Failed to seek file"))
		return
	}

	// Validate mime type
	mimeType := http.DetectContentType(buff)
	if isBlockedMimeType(mimeType) {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "FILE_TYPE_NOT_ALLOWED", fmt.Sprintf("File type %s is not allowed", mimeType)))
		return
	}

	collection := r.FormValue("collection")
	recordID := r.FormValue("recordID")
	preserveName := r.FormValue("preserve_name") == "true"

	if collection == "" || recordID == "" {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "MISSING_PARAMS", "collection and recordID are required"))
		return
	}

	// Generate a safe unique name or use original sanitized name
	var safeName string
	if preserveName {
		// Basic sanitization: remove path traversal, special characters
		safeName = sanitizeFileName(header.Filename)
	} else {
		ext := filepath.Ext(header.Filename)
		safeName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
	}
	path := filepath.Join(collection, recordID, safeName)

	if err := h.storage.Save(r.Context(), path, file); err != nil {
		errors.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusCreated, map[string]any{
		"name": safeName,
		"size": header.Size,
		"mime": header.Header.Get("Content-Type"),
		"url":  fmt.Sprintf("/api/files/%s/%s/%s", collection, recordID, safeName),
	}, nil)
}

func isBlockedMimeType(mimeType string) bool {
	// Block potentially dangerous file types
	blocked := map[string]bool{
		"application/x-dosexec":    true, // Windows Executables
		"application/x-msdownload": true, // DLLs, etc.
		"application/x-sh":         true, // Shell scripts
		"text/html":                true, // Prevent Stored XSS
		"text/javascript":          true, // JS files
	}
	// Check exact match or prefix
	if blocked[mimeType] {
		return true
	}
	return false
}

func sanitizeFileName(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "..", "")
	// Remove non-alphanumeric (keep ., -, _)
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}
	return result.String()
}
