package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/zulfikawr/vault/internal/core"
	"net/http"

	_ "modernc.org/sqlite"
)

func Connect(ctx context.Context, path string) (*sql.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, core.NewError(http.StatusInternalServerError, "DB_DIR_CREATION_FAILED", "Failed to create data directory").WithDetails(map[string]any{"error": err.Error(), "path": dir})
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, core.NewError(http.StatusInternalServerError, "DB_OPEN_FAILED", "Failed to open database").WithDetails(map[string]any{"error": err.Error(), "path": path})
	}

	// Performance and stability settings
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA foreign_keys=ON;",
	}

	for _, pragma := range pragmas {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			_ = db.Close()
			return nil, core.NewError(http.StatusInternalServerError, "DB_PRAGMA_FAILED", "Failed to execute pragma").WithDetails(map[string]any{"error": err.Error(), "pragma": pragma})
		}
	}

	// SQLite typically works best with a single writer
	db.SetMaxOpenConns(1)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, core.NewError(http.StatusInternalServerError, "DB_PING_FAILED", "Failed to ping database").WithDetails(map[string]any{"error": err.Error()})
	}

	return db, nil
}
