package api

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/api/middleware"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/realtime"
	"github.com/zulfikawr/vault/internal/service"
	"github.com/zulfikawr/vault/internal/storage"
	"github.com/zulfikawr/vault/ui"
)

func NewRouter(
	recordService *service.RecordService,
	collectionService *service.CollectionService,
	sqlService *service.SqlService,
	registry *db.SchemaRegistry,
	store storage.Storage,
	hub *realtime.Hub,
	config *core.Config,
) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(recordService, config)
	crudHandler := NewCollectionHandler(recordService, registry)
	fileHandler := NewFileHandler(store, config.MaxFileUploadSize)
	realtimeHandler := NewRealtimeHandler(hub)
	adminHandler := NewAdminHandler(collectionService, sqlService)
	logsHandler := NewLogsHandler()
	settingsHandler := NewSettingsHandler(config)
	storageHandler := NewStorageHandler(config.DataDir + "/storage")

	// Base routes
	uiHandler := ui.Handler()
	mux.Handle("/", uiHandler)
	mux.Handle("/_/", http.StripPrefix("/_", uiHandler))

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	mux.HandleFunc("GET /api/health/collections", func(w http.ResponseWriter, r *http.Request) {
		collections := collectionService.ListCollections()
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
	adminRouter.HandleFunc("PATCH /collections/{id}", adminHandler.UpdateCollection)
	adminRouter.HandleFunc("DELETE /collections/{name}", adminHandler.DeleteCollection)
	adminRouter.HandleFunc("GET /settings", settingsHandler.GetSettings)
	adminRouter.HandleFunc("PATCH /settings", settingsHandler.UpdateSettings)
	adminRouter.HandleFunc("POST /backups", adminHandler.CreateBackup)
	adminRouter.HandleFunc("GET /logs", logsHandler.GetLogs)
	adminRouter.HandleFunc("DELETE /logs", logsHandler.ClearLogs)
	adminRouter.HandleFunc("GET /storage", storageHandler.List)
	adminRouter.HandleFunc("GET /storage/stats", storageHandler.Stats)
	adminRouter.HandleFunc("DELETE /storage", storageHandler.Delete)
	adminRouter.HandleFunc("POST /query", adminHandler.ExecuteQuery)

	// Apply rate limiting to admin operations
	mux.Handle("/api/admin/", http.StripPrefix("/api/admin", middleware.RateLimitMiddleware(config.RateLimitPerMin)(middleware.AdminOnly(adminRouter))))

	return mux
}
