package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
)

type MigrationEngine struct {
	db *sql.DB
}

func NewMigrationEngine(db *sql.DB) *MigrationEngine {
	return &MigrationEngine{db: db}
}

func (m *MigrationEngine) SyncCollection(ctx context.Context, c *models.Collection) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_TX_BEGIN_FAILED", "Failed to begin transaction").WithDetails(map[string]any{"error": err.Error()})
	}
	defer func() { _ = tx.Rollback() }()

	var tableName string
	err = tx.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", c.Name).Scan(&tableName)

	if err == sql.ErrNoRows {
		if err := m.createTableTx(ctx, tx, c); err != nil {
			return err
		}
	} else if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_SYNC_FAILED", "Failed to check table existence").WithDetails(map[string]any{"error": err.Error(), "collection": c.Name})
	} else {
		if err := m.updateTableTx(ctx, tx, c); err != nil {
			return err
		}
	}

	if err := m.syncIndexes(ctx, tx, c); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MigrationEngine) syncIndexes(ctx context.Context, tx *sql.Tx, c *models.Collection) error {
	for _, idx := range c.Indexes {
		// Clean up field names
		fields := strings.Split(idx, ",")
		safeFields := make([]string, 0, len(fields))
		for _, f := range fields {
			safeFields = append(safeFields, strings.TrimSpace(f))
		}

		if len(safeFields) == 0 {
			continue
		}

		indexName := fmt.Sprintf("idx_%s_%s", c.Name, strings.Join(safeFields, "_"))
		// SQLite handles IF NOT EXISTS
		query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s)", indexName, c.Name, strings.Join(safeFields, ","))

		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return errors.NewError(http.StatusInternalServerError, "DB_CREATE_INDEX_FAILED", "Failed to create index").WithDetails(map[string]any{"error": err.Error(), "query": query})
		}
	}
	return nil
}

func (m *MigrationEngine) createTableTx(ctx context.Context, tx *sql.Tx, c *models.Collection) error {
	columns := []string{
		"id TEXT PRIMARY KEY",
		"created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))",
		"updated TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))",
	}

	for _, f := range c.Fields {
		sqlType := "TEXT"
		switch f.Type {
		case models.FieldTypeNumber:
			sqlType = "REAL"
		case models.FieldTypeBool:
			sqlType = "INTEGER"
		}

		col := fmt.Sprintf("%s %s", f.Name, sqlType)
		if f.Required {
			col += " NOT NULL"
		}
		if f.Unique {
			col += " UNIQUE"
		}
		columns = append(columns, col)
	}

	query := fmt.Sprintf("CREATE TABLE %s (%s)", c.Name, strings.Join(columns, ", "))
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_CREATE_TABLE_FAILED", "Failed to create table").WithDetails(map[string]any{"error": err.Error(), "query": query})
	}

	slog.Debug("Created table", "collection", c.Name, "request_id", core.GetRequestID(ctx))
	return nil
}

func (m *MigrationEngine) updateTableTx(ctx context.Context, tx *sql.Tx, c *models.Collection) error {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", c.Name))
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_PRAGMA_FAILED", "Failed to get table info").WithDetails(map[string]any{"error": err.Error()})
	}
	defer errors.Defer(ctx, rows.Close, "close rows")

	existingCols := make(map[string]bool)
	for rows.Next() {
		var cid sql.NullInt64
		var name, dtype sql.NullString
		var notnull sql.NullInt64
		var dfltValue sql.NullString
		var pk sql.NullInt64
		if err := rows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk); err != nil {
			slog.Error("PRAGMA scan failed", "error", err, "collection", c.Name)
			return errors.NewError(http.StatusInternalServerError, "DB_SCAN_FAILED", "Failed to scan table info").WithDetails(map[string]any{"error": err.Error(), "collection": c.Name})
		}
		if name.Valid {
			existingCols[name.String] = true
		}
	}

	for _, f := range c.Fields {
		if !existingCols[f.Name] {
			sqlType := "TEXT"
			switch f.Type {
			case models.FieldTypeNumber:
				sqlType = "REAL"
			case models.FieldTypeBool:
				sqlType = "INTEGER"
			}

			query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", c.Name, f.Name, sqlType)
			if _, err := tx.ExecContext(ctx, query); err != nil {
				return errors.NewError(http.StatusInternalServerError, "DB_ALTER_TABLE_FAILED", "Failed to add column").WithDetails(map[string]any{"error": err.Error(), "field": f.Name})
			}
			slog.Debug("Added column", "collection", c.Name, "field", f.Name, "request_id", core.GetRequestID(ctx))
		}
	}

	// Ensure system columns exist
	systemCols := map[string]string{
		"created": "TEXT DEFAULT ''",
		"updated": "TEXT DEFAULT ''",
	}

	for col, definition := range systemCols {
		if !existingCols[col] {
			query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", c.Name, col, definition)
			if _, err := tx.ExecContext(ctx, query); err != nil {
				return errors.NewError(http.StatusInternalServerError, "DB_ALTER_TABLE_FAILED", "Failed to add system column").WithDetails(map[string]any{"error": err.Error(), "field": col})
			}
			slog.Debug("Added system column", "collection", c.Name, "field", col, "request_id", core.GetRequestID(ctx))
			
			// Try to populate from created_at/updated_at if they exist
			if col == "created" && existingCols["created_at"] {
				updateQuery := fmt.Sprintf("UPDATE %s SET created = strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', created_at, 'unixepoch')", c.Name)
				_, _ = tx.ExecContext(ctx, updateQuery)
			}
			if col == "updated" && existingCols["updated_at"] {
				updateQuery := fmt.Sprintf("UPDATE %s SET updated = strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', updated_at, 'unixepoch')", c.Name)
				_, _ = tx.ExecContext(ctx, updateQuery)
			}
			
			// If still empty, use now
			updateNow := fmt.Sprintf("UPDATE %s SET %s = strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', 'now') WHERE %s = ''", c.Name, col, col)
			_, _ = tx.ExecContext(ctx, updateNow)
		}
	}

	return nil
}

func (m *MigrationEngine) DropCollection(ctx context.Context, name string) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_TX_BEGIN_FAILED", "Failed to begin transaction").WithDetails(map[string]any{"error": err.Error()})
	}
	defer func() { _ = tx.Rollback() }()

	// Drop the table
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_DROP_TABLE_FAILED", "Failed to drop table").WithDetails(map[string]any{"error": err.Error(), "query": query})
	}

	// Remove from _collections table
	_, err = tx.ExecContext(ctx, "DELETE FROM _collections WHERE name = ?", name)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_DELETE_DEFINITION_FAILED", "Failed to delete collection definition").WithDetails(map[string]any{"error": err.Error()})
	}

	err = tx.Commit()
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "DB_COMMIT_FAILED", "Failed to commit transaction").WithDetails(map[string]any{"error": err.Error()})
	}

	slog.Debug("Dropped collection", "collection", name, "request_id", core.GetRequestID(ctx))
	return nil
}
