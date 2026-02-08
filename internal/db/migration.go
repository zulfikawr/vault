package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zulfikawr/vault/internal/models"
)

type MigrationEngine struct {
	db *sql.DB
}

func NewMigrationEngine(db *sql.DB) *MigrationEngine {
	return &MigrationEngine{db: db}
}

func (m *MigrationEngine) SyncCollection(c *models.Collection) error {
	// 1. Check if table exists
	var tableName string
	err := m.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", c.Name).Scan(&tableName)
	
	if err == sql.ErrNoRows {
		return m.createTable(c)
	} else if err != nil {
		return err
	}

	return m.updateTable(c)
}

func (m *MigrationEngine) createTable(c *models.Collection) error {
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
	_, err := m.db.Exec(query)
	return err
}

func (m *MigrationEngine) updateTable(c *models.Collection) error {
	// For simplicity in Phase 2, we'll only handle adding new columns
	rows, err := m.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", c.Name))
	if err != nil {
		return err
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
			if _, err := m.db.Exec(query); err != nil {
				return err
			}
		}
	}

	return nil
}
