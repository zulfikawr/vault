package api

import (
	"log/slog"
	"net/http"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/realtime"
	"github.com/zulfikawr/vault/internal/storage"
	"github.com/zulfikawr/vault/ui"
)

func NewRouter(executor *db.Executor, registry *db.SchemaRegistry, store storage.Storage, hub *realtime.Hub, migration *db.MigrationEngine, config *core.Config) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(executor, config)
	crudHandler := NewCollectionHandler(executor, registry)
	fileHandler := NewFileHandler(store, executor)
	realtimeHandler := NewRealtimeHandler(hub)
	adminHandler := NewAdminHandler(executor, registry, migration)

	// Base routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			slog.Info("Redirecting root to /_/", "request_id", core.GetRequestID(r.Context()))
			http.Redirect(w, r, "/_/", http.StatusFound)
			return
		}
		
		// If it's not root and not matched by other handlers, it's a 404
		http.NotFound(w, r)
	})
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

	// File routes
	mux.HandleFunc("GET /api/files/{collection}/{id}/{filename}", fileHandler.Serve)
	mux.HandleFunc("POST /api/files", fileHandler.Upload)

	// Realtime routes
	mux.HandleFunc("GET /api/realtime", realtimeHandler.Connect)

	// Admin routes (Protected by AdminOnly)
	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("GET /collections", adminHandler.ListCollections)
	adminRouter.HandleFunc("POST /collections", adminHandler.CreateCollection)
	adminRouter.HandleFunc("GET /settings", adminHandler.GetSettings)
	adminRouter.HandleFunc("POST /backups", adminHandler.CreateBackup)

	// Mount admin router with middleware
	mux.Handle("/api/admin/", http.StripPrefix("/api/admin", AdminOnly(adminRouter)))

	// Mount UI handler
	mux.Handle("/_/", http.StripPrefix("/_", ui.Handler()))

	return mux
}
