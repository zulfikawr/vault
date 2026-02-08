package api

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

func NewRouter(executor *db.Executor, config *core.Config) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(executor, config)

	// Base routes
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	// Auth routes
	mux.HandleFunc("POST /api/collections/users/auth-with-password", authHandler.Login)
	mux.HandleFunc("POST /api/collections/users/auth-refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/collections/users/request-password-reset", authHandler.RequestPasswordReset)
	mux.HandleFunc("POST /api/collections/users/confirm-password-reset", authHandler.ConfirmPasswordReset)

	return mux
}