package api

import (
	"encoding/json"
	"net/http"

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
	SendJSON(w, http.StatusOK, map[string]any{
		"token":  token,
		"record": userRecord,
	}, nil)
}
