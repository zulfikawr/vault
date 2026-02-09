package api

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

func NewRouter(executor *db.Executor, registry *db.SchemaRegistry, config *core.Config) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(executor, config)
	crudHandler := NewCollectionHandler(executor, registry)

	// Base routes
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	// Auth routes
	mux.HandleFunc("POST /api/collections/users/auth-with-password", authHandler.Login)
	mux.HandleFunc("POST /api/collections/users/auth-refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/collections/users/request-password-reset", authHandler.RequestPasswordReset)
	mux.HandleFunc("POST /api/collections/users/confirm-password-reset", authHandler.ConfirmPasswordReset)

	// CRUD routes (Dynamic)
	mux.HandleFunc("GET /api/collections/{collection}/records", crudHandler.List)
	mux.HandleFunc("POST /api/collections/{collection}/records", crudHandler.Create)
	mux.HandleFunc("GET /api/collections/{collection}/records/{id}", crudHandler.View)
	mux.HandleFunc("PATCH /api/collections/{collection}/records/{id}", crudHandler.Update)
	mux.HandleFunc("DELETE /api/collections/{collection}/records/{id}", crudHandler.Delete)

	return mux
}