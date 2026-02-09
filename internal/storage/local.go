package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zulfikawr/vault/internal/core"
)

type Local struct {
	basePath string
}

func NewLocal(basePath string) (*Local, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base storage directory: %w", err)
	}
	return &Local{basePath: basePath}, nil
}

func (l *Local) Save(ctx context.Context, path string, data io.Reader) error {
	fullPath := filepath.Join(l.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return core.NewError(http.StatusInternalServerError, "STORAGE_DIR_CREATE_FAILED", "Failed to create directory").WithDetails(map[string]any{"error": err.Error(), "path": path})
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return core.NewError(http.StatusInternalServerError, "STORAGE_CREATE_FAILED", "Failed to create file").WithDetails(map[string]any{"error": err.Error(), "path": path})
	}
	defer func() { _ = file.Close() }()

	if _, err := io.Copy(file, data); err != nil {
		return core.NewError(http.StatusInternalServerError, "STORAGE_WRITE_FAILED", "Failed to write data").WithDetails(map[string]any{"error": err.Error(), "path": path})
	}

	return nil
}

func (l *Local) Retrieve(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.basePath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, core.NewError(http.StatusNotFound, "FILE_NOT_FOUND", "File not found")
		}
		return nil, core.NewError(http.StatusInternalServerError, "STORAGE_READ_FAILED", "Failed to read file").WithDetails(map[string]any{"error": err.Error(), "path": path})
	}
	return file, nil
}

func (l *Local) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.basePath, path)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return core.NewError(http.StatusInternalServerError, "STORAGE_DELETE_FAILED", "Failed to delete file").WithDetails(map[string]any{"error": err.Error(), "path": path})
	}
	return nil
}

func (l *Local) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(l.basePath, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, core.NewError(http.StatusInternalServerError, "STORAGE_STAT_FAILED", "Failed to stat file").WithDetails(map[string]any{"error": err.Error(), "path": path})
}
