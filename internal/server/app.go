package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/zulfikawr/vault/internal/api"
	"github.com/zulfikawr/vault/internal/api/middleware"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/realtime"
	"github.com/zulfikawr/vault/internal/service"
	"github.com/zulfikawr/vault/internal/storage"
)

type App struct {
	Config            *core.Config
	DB                *sql.DB
	Server            *Server
	Registry          *db.SchemaRegistry
	Migration         *db.MigrationEngine
	RecordService     *service.RecordService
	CollectionService *service.CollectionService
	Storage           storage.Storage
	Hub               *realtime.Hub
}

func NewApp(cfg *core.Config) *App {
	// Ensure data directory exists
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		slog.Error("Failed to create data directory", "error", err, "path", cfg.DataDir)
	}

	logPath := filepath.Join(cfg.DataDir, "vault.log")
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, logPath); err != nil {
		_ = core.InitLogger(cfg.LogLevel, cfg.LogFormat, "")
		slog.Error("Failed to initialize file logger, using stdout", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, cfg.DBPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Initialize Storage
	store, err := storage.NewLocal(cfg.DataDir + "/storage")
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Initialize Realtime Hub
	hub := realtime.NewHub()
	// Hub is started in Run()

	registry := db.NewSchemaRegistry(database)
	migration := db.NewMigrationEngine(database)
	repo := db.NewRepository(database, registry)

	recordService := service.NewRecordService(repo, hub)
	collectionService := service.NewCollectionService(registry, migration)

	// Register Auth Hooks
	service.RegisterAuthHooks()

	// Bootstrap system
	if err := collectionService.InitSystem(ctx); err != nil {
		slog.Error("Failed to initialize system", "error", err)
		os.Exit(1)
	}

	// Load dynamic collections from DB
	if err := registry.LoadFromDB(ctx); err != nil {
		slog.Error("Failed to load dynamic collections", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(recordService, collectionService, registry, store, hub, cfg)
	handler := middleware.Chain(router,
		middleware.RecoveryMiddleware,
		middleware.LoggerMiddleware,
		middleware.SecurityMiddleware,
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RequestIDMiddleware,
		middleware.CORSMiddleware,
	)

	srv := NewServer(cfg.Port, handler)

	return &App{
		Config:            cfg,
		DB:                database,
		Server:            srv,
		Registry:          registry,
		Migration:         migration,
		RecordService:     recordService,
		CollectionService: collectionService,
		Storage:           store,
		Hub:               hub,
	}
}

func (a *App) Run() {
	// Root context for the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Realtime Hub
	go a.Hub.Run(ctx)

	// Check if any users exist
	var count int
	err := a.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err == nil && count == 0 {
		fmt.Println("⚠️  No users found! Create an admin user with:")
		fmt.Println("./vault admin create --email <email> --password <password> --username <username>")
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverError := make(chan error, 1)
	go func() {
		if err := a.Server.Start(); err != nil {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		fmt.Fprintf(os.Stderr, "\n❌ Error: Could not start server\n")
		if strings.Contains(err.Error(), "bind: address already in use") {
			fmt.Fprintf(os.Stderr, "   Port %d is already in use by another process.\n", a.Config.Port)
			fmt.Fprintf(os.Stderr, "   Please use a different port with the --port flag.\n\n")
		} else {
			fmt.Fprintf(os.Stderr, "   %v\n\n", err)
		}
		slog.Error("Server failed to start", "error", err)
		cancel() // Stop Hub
		os.Exit(1)
	case <-stop:
		fmt.Println("\nShutting down server...")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.Server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}

	// Cancel root context to stop Hub and other background tasks
	cancel()

	a.RecordService.Close()

	if err := a.DB.Close(); err != nil {
		slog.Error("Database close failed", "error", err)
	}

	fmt.Println("Vault stopped cleanly")
}
