package errors

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// VaultError represents a structured error with HTTP status and details
type VaultError struct {
	Status  int            `json:"-"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *VaultError) Error() string {
	return e.Message
}

// NewError creates a new VaultError with status, code, and message
func NewError(status int, code string, message string) *VaultError {
	return &VaultError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

// WithDetails adds additional context to the error
func (e *VaultError) WithDetails(details map[string]any) *VaultError {
	e.Details = details
	return e
}

// SendError writes a VaultError as JSON response
func SendError(w http.ResponseWriter, err error) {
	ve, ok := err.(*VaultError)
	if !ok {
		ve = NewError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ve.Status)
	if err := json.NewEncoder(w).Encode(map[string]*VaultError{"error": ve}); err != nil {
		slog.Error("Failed to encode error response", "error", err)
	}
}

// Log logs an error with context and attributes if err is not nil
func Log(ctx context.Context, err error, msg string, attrs ...any) {
	if err == nil {
		return
	}

	requestID := getRequestID(ctx)
	logAttrs := []any{"error", err}
	if requestID != "" {
		logAttrs = append(logAttrs, "request_id", requestID)
	}
	logAttrs = append(logAttrs, attrs...)

	slog.ErrorContext(ctx, msg, logAttrs...)
}

// Check logs an error if not nil and returns true if no error occurred
func Check(ctx context.Context, err error, msg string, attrs ...any) bool {
	if err != nil {
		Log(ctx, err, msg, attrs...)
		return false
	}
	return true
}

// Defer handles errors in defer statements with proper logging
func Defer(ctx context.Context, fn func() error, msg string, attrs ...any) {
	if err := fn(); err != nil {
		Log(ctx, err, msg, attrs...)
	}
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	type contextKey string
	const RequestIDKey contextKey = "request_id"

	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
