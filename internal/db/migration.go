package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/models"
)

type MigrationEngine struct {
	db *sql.DB
}

func NewMigrationEngine(db *sql.DB) *MigrationEngine {
	return &MigrationEngine{db: db}
}

func (m *MigrationEngine) SyncCollection(ctx context.Context, c *models.Collection) error {
	var tableName string
	err := m.db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", c.Name).Scan(&tableName)
	
	if err == sql.ErrNoRows {
		return m.createTable(ctx, c)
	} else if err != nil {
		return core.NewError(http.StatusInternalServerError, "DB_SYNC_FAILED", "Failed to check table existence").WithDetails(map[string]any{"error": err.Error(), "collection": c.Name})
	}

	return m.updateTable(ctx, c)
}

func (m *MigrationEngine) createTable(ctx context.Context, c *models.Collection) error {
	columns := []string{
		"id TEXT PRIMARY KEY",
		"created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))",
		"updated TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))",
	}

	for _, f := range c.Fields {
		sqlType := "TEXT"
		if f.Type == models.FieldTypeNumber {
			sqlType = "REAL"
		} else if f.Type == models.FieldTypeBool {
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
	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return core.NewError(http.StatusInternalServerError, "DB_CREATE_TABLE_FAILED", "Failed to create table").WithDetails(map[string]any{"error": err.Error(), "query": query})
	}

	slog.Info("Created table", "collection", c.Name, "request_id", core.GetRequestID(ctx))
	return nil
}

func (m *MigrationEngine) updateTable(ctx context.Context, c *models.Collection) error {
	rows, err := m.db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", c.Name))
	if err != nil {
		return core.NewError(http.StatusInternalServerError, "DB_PRAGMA_FAILED", "Failed to get table info").WithDetails(map[string]any{"error": err.Error()})
	}
	defer rows.Close()

	existingCols := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, dtype string
		var notnull, pk int
		var dfltValue any
		rows.Scan(&cid, &name, &dtype, &notnull, &pk, &dfltValue)
		existingCols[name] = true
	}

	for _, f := range c.Fields {
		if !existingCols[f.Name] {
			sqlType := "TEXT"
			if f.Type == models.FieldTypeNumber {
				sqlType = "REAL"
			} else if f.Type == models.FieldTypeBool {
				sqlType = "INTEGER"
			}
			
			query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", c.Name, f.Name, sqlType)
			if _, err := m.db.ExecContext(ctx, query); err != nil {
				return core.NewError(http.StatusInternalServerError, "DB_ALTER_TABLE_FAILED", "Failed to add column").WithDetails(map[string]any{"error": err.Error(), "field": f.Name})
			}
			slog.Info("Added column", "collection", c.Name, "field", f.Name, "request_id", core.GetRequestID(ctx))
		}
	}

	return nil
}