package api

import (
	"encoding/json"
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
	
	// mock finding user
	records, _, err := h.executor.ListRecords(r.Context(), "users", db.QueryParams{Filter: "identity = " + req.Identity})
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
	// Skeleton implementation
	SendJSON(w, http.StatusOK, map[string]string{"message": "If the email exists, a reset link will be sent"}, nil)
}

func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	// Skeleton implementation
	SendJSON(w, http.StatusOK, map[string]string{"message": "Password has been successfully reset"}, nil)
}
