package api

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zulfikawr/vault/internal/errors"
)

type StorageHandler struct {
	basePath string
}

func NewStorageHandler(basePath string) *StorageHandler {
	return &StorageHandler{basePath: basePath}
}

type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Modified int64  `json:"modified"`
	MimeType string `json:"mime_type,omitempty"`
}

type StorageStats struct {
	TotalFiles       int   `json:"total_files"`
	TotalSize        int64 `json:"total_size"`
	TotalCollections int   `json:"total_collections"`
}

type StorageListResponse struct {
	Files   []FileInfo    `json:"files"`
	Folders []FileInfo    `json:"folders"`
	Stats   *StorageStats `json:"stats,omitempty"`
}

func (h *StorageHandler) List(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	// Prevent directory traversal
	if strings.Contains(path, "..") {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_PATH", "Invalid path"))
		return
	}

	fullPath := filepath.Join(h.basePath, path)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			SendJSON(w, http.StatusOK, StorageListResponse{Files: []FileInfo{}, Folders: []FileInfo{}}, nil)
			return
		}
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "READ_DIR_FAILED", "Failed to read directory"))
		return
	}

	var files, folders []FileInfo

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		relPath := filepath.Join(path, entry.Name())
		fileInfo := FileInfo{
			Name:     entry.Name(),
			Path:     relPath,
			Size:     info.Size(),
			IsDir:    entry.IsDir(),
			Modified: info.ModTime().Unix(),
		}

		if !entry.IsDir() {
			fileInfo.MimeType = detectMimeType(entry.Name())
			files = append(files, fileInfo)
		} else {
			folders = append(folders, fileInfo)
		}
	}

	SendJSON(w, http.StatusOK, StorageListResponse{
		Files:   files,
		Folders: folders,
	}, nil)
}

func (h *StorageHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats := &StorageStats{}

	// Count collections (top-level directories)
	entries, err := os.ReadDir(h.basePath)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				stats.TotalCollections++
			}
		}
	}

	// Walk entire storage to count files and size
	if err := filepath.WalkDir(h.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		stats.TotalFiles++
		stats.TotalSize += info.Size()
		return nil
	}); err != nil {
		errors.Log(r.Context(), err, "walk storage directory", "path", h.basePath)
	}

	SendJSON(w, http.StatusOK, stats, nil)
}

func (h *StorageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_JSON", "Invalid request body"))
		return
	}

	if req.Path == "" || strings.Contains(req.Path, "..") {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_PATH", "Invalid path"))
		return
	}

	fullPath := filepath.Join(h.basePath, req.Path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			errors.SendError(w, errors.NewError(http.StatusNotFound, "FILE_NOT_FOUND", "File not found"))
			return
		}
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete file"))
		return
	}

	SendJSON(w, http.StatusOK, map[string]string{"message": "File deleted successfully"}, nil)
}

func detectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".zip":  "application/zip",
		".mp3":  "audio/mpeg",
		".mp4":  "video/mp4",
		".txt":  "text/plain",
		".json": "application/json",
		".xml":  "application/xml",
	}

	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}
