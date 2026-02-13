package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
	"net/http"
)

type AuditLog struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	AdminID   string    `json:"admin_id"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

func LogAuditEvent(ctx context.Context, db *sql.DB, action, resource, adminID string, details map[string]any) error {
	detailsJSON, _ := json.Marshal(details)

	query := `INSERT INTO _audit_logs (action, resource, admin_id, details, timestamp) 
	          VALUES (?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query, action, resource, adminID, string(detailsJSON), time.Now().Format(time.RFC3339))
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "AUDIT_LOG_FAILED", "Failed to log audit event").WithDetails(map[string]any{"error": err.Error()})
	}

	return nil
}

func BootstrapAuditLogsCollection(registry *SchemaRegistry) error {
	auditTable := &models.Collection{
		ID:   "system_audit_logs",
		Name: "_audit_logs",
		Type: models.CollectionTypeSystem,
		Fields: []models.Field{
			{Name: "action", Type: models.FieldTypeText, Required: true},
			{Name: "resource", Type: models.FieldTypeText, Required: true},
			{Name: "admin_id", Type: models.FieldTypeText, Required: true},
			{Name: "details", Type: models.FieldTypeJSON},
			{Name: "timestamp", Type: models.FieldTypeDate, Required: true},
		},
		Created: time.Now().Format(time.RFC3339),
		Updated: time.Now().Format(time.RFC3339),
	}

	registry.AddCollection(auditTable)
	return nil
}
