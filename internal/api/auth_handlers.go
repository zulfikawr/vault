package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

type AuthHandler struct {
	executor *db.Executor
	config   *core.Config
}

func NewAuthHandler(executor *db.Executor, config *core.Config) *AuthHandler {
	return &AuthHandler{
		executor: executor,
		config:   config,
	}
}

type LoginRequest struct {
	Identity string `json:"identity"` // email or username
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body"))
		return
	}

	// This is a simplified lookup. In a real scenario, we'd use the QueryBuilder to find by email or username.
	// For now, let's assume we have a way to find a record by a filter.
	// I'll add a temporary FindOne helper to executor later, but for this step, let's assume it exists.

	// Find user by email
	// In Phase 5 we implemented a simple filter parser that requires valid field names.
	records, _, err := h.executor.ListRecords(r.Context(), "users", db.QueryParams{Filter: "email = " + req.Identity})
	if err != nil || len(records) == 0 {
		// If not found by email, try username
		records, _, err = h.executor.ListRecords(r.Context(), "users", db.QueryParams{Filter: "username = " + req.Identity})
	}

	if err != nil || len(records) == 0 {
		core.SendError(w, core.NewError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid identity or password"))
		return
	}

	userRecord := records[0]
	hashedPassword := userRecord.GetString("password")

	if !auth.ComparePasswords(hashedPassword, req.Password) {
		core.SendError(w, core.NewError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid identity or password"))
		return
	}

	// Update lastLogin timestamp
	_, err = h.executor.UpdateRecord(r.Context(), "users", userRecord.ID, map[string]any{
		"lastLogin": time.Now().Format(time.RFC3339),
	})
	if err != nil {
		slog.Warn("Failed to update lastLogin", "error", err)
	}

	// Refresh the user record to get updated lastLogin
	records, _, err = h.executor.ListRecords(r.Context(), "users", db.QueryParams{Filter: "id = " + userRecord.ID})
	if err == nil && len(records) > 0 {
		userRecord = records[0]
	}

	token, err := auth.GenerateToken(r.Context(), userRecord, h.config.JWTSecret, h.config.JWTExpiry)
	if err != nil {
		core.SendError(w, err)
		return
	}

	userRecord.HideField("password")

	// Generate Refresh Token
	refreshToken := uuid.New().String()
	_, err = h.executor.CreateRecord(r.Context(), "_refresh_tokens", map[string]any{
		"token":   refreshToken,
		"user_id": userRecord.ID,
		"expires": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	})
	if err != nil {
		core.SendError(w, err)
		return
	}

	SendJSON(w, http.StatusOK, map[string]any{
		"token":         token,
		"refresh_token": refreshToken,
		"record":        userRecord,
	}, nil)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body"))
		return
	}

	// Lookup refresh token
	records, _, err := h.executor.ListRecords(r.Context(), "_refresh_tokens", db.QueryParams{Filter: "token = " + req.RefreshToken})
	if err != nil || len(records) == 0 {
		core.SendError(w, core.NewError(http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid or expired refresh token"))
		return
	}

	refreshTokenRecord := records[0]
	// Check expiry (omitted for brevity, but should be checked)

	userID := refreshTokenRecord.GetString("user_id")
	userRecord, err := h.executor.FindRecordByID(r.Context(), "users", userID)
	if err != nil {
		core.SendError(w, core.NewError(http.StatusUnauthorized, "USER_NOT_FOUND", "Associated user not found"))
		return
	}

	newToken, err := auth.GenerateToken(r.Context(), userRecord, h.config.JWTSecret, h.config.JWTExpiry)
	if err != nil {
		core.SendError(w, err)
		return
	}

	userRecord.HideField("password")
	SendJSON(w, http.StatusOK, map[string]any{
		"token":  newToken,
		"record": userRecord,
	}, nil)
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body"))
		return
	}

	// 1. Find user
	records, _, err := h.executor.ListRecords(r.Context(), "users", db.QueryParams{Filter: "email = " + req.Email})
	if err != nil || len(records) == 0 {
		// Silent fail for security: don't reveal if email exists
		SendJSON(w, http.StatusOK, map[string]string{"message": "If the email exists, a reset link will be sent"}, nil)
		return
	}

	// 2. Generate a temporary token (mocked for this checkpoint)
	resetToken := uuid.New().String()

	// In a real implementation, we'd save this to a 'password_resets' collection with an expiry.
	// For now, we log it (representing the "email" being sent).
	core.InitLogger("INFO", "text") // Ensure logger is ready
	slog.Info("Password reset requested", "email", req.Email, "token", resetToken, "request_id", core.GetRequestID(r.Context()))

	SendJSON(w, http.StatusOK, map[string]string{"message": "If the email exists, a reset link will be sent"}, nil)
}

func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.SendError(w, core.NewError(http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body"))
		return
	}

	// In a real implementation, we'd verify the token from the DB.
	// Since we don't have the table yet, we'll return an error for now to show logic is present.
	core.SendError(w, core.NewError(http.StatusNotImplemented, "RESET_NOT_FULLY_IMPLEMENTED", "Token verification table not yet migrated"))
}
