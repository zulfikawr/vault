package api

import (
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
	logsHandler := NewLogsHandler()
	settingsHandler := NewSettingsHandler(config)
	storageHandler := NewStorageHandler(config.StoragePath())

	// Rate limiter for collection operations (10 per minute default)
	collectionLimiter := NewRateLimiter(10)

	// Base routes
	// Mount UI handler - serves at both / and /_/ for devtunnel compatibility
	uiHandler := ui.Handler()
	mux.Handle("/", uiHandler)
	mux.Handle("/_/", http.StripPrefix("/_", uiHandler))
	
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	mux.HandleFunc("GET /api/health/collections", func(w http.ResponseWriter, r *http.Request) {
		collections := registry.GetCollections()
		SendJSON(w, http.StatusOK, map[string]any{
			"collections": len(collections),
			"status":      "healthy",
		}, nil)
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
	adminRouter.HandleFunc("GET /settings", settingsHandler.GetSettings)
	adminRouter.HandleFunc("PATCH /settings", settingsHandler.UpdateSettings)
	adminRouter.HandleFunc("POST /backups", adminHandler.CreateBackup)
	adminRouter.HandleFunc("GET /logs", logsHandler.GetLogs)
	adminRouter.HandleFunc("DELETE /logs", logsHandler.ClearLogs)
	adminRouter.HandleFunc("GET /storage", storageHandler.List)
	adminRouter.HandleFunc("GET /storage/stats", storageHandler.Stats)
	adminRouter.HandleFunc("DELETE /storage", storageHandler.Delete)

	// Apply rate limiting to collection operations
	collectionRouter := http.NewServeMux()
	collectionRouter.HandleFunc("GET /collections", adminHandler.ListCollections)
	collectionRouter.HandleFunc("POST /collections", adminHandler.CreateCollection)

	// Mount admin router with middleware
	mux.Handle("/api/admin/", http.StripPrefix("/api/admin", AdminOnly(adminRouter)))
	mux.Handle("/api/admin/collections", http.StripPrefix("/api/admin", AdminOnly(RateLimitMiddleware(collectionLimiter)(collectionRouter))))

	return mux
}
