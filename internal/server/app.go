package server

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zulfikawr/vault/internal/api"
	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/models"
	"github.com/zulfikawr/vault/internal/realtime"
	"github.com/zulfikawr/vault/internal/storage"
)

type App struct {
	Config    *core.Config
	DB        *sql.DB
	Server    *Server
	Registry  *db.SchemaRegistry
	Migration *db.MigrationEngine
	Executor  *db.Executor
	Storage   storage.Storage
	Hub       *realtime.Hub
}

func NewApp() *App {
	cfg := core.LoadConfig()
	
	if err := core.InitLoggerWithFile(cfg.LogLevel, cfg.LogFormat, "./vault_data/vault.log"); err != nil {
		core.InitLogger(cfg.LogLevel, cfg.LogFormat)
		slog.Error("Failed to initialize file logger, using stdout", "error", err)
	}
	
	if err := core.InitFileLogger("./vault_data"); err != nil {
		slog.Error("Failed to initialize file logger", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, cfg.DBPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Initialize Storage
	store, err := storage.NewLocal("./vault_data/storage")
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Initialize Realtime Hub
	hub := realtime.NewHub()
	go hub.Run(context.Background())

	registry := db.NewSchemaRegistry(database)
	migration := db.NewMigrationEngine(database)
	executor := db.NewExecutor(database, registry, hub)

	// Register Auth Hooks
	db.GetHooks("users").BeforeCreate = append(db.GetHooks("users").BeforeCreate, func(ctx context.Context, record *models.Record) error {
		password := record.GetString("password")
		if password == "" {
			return nil
		}
		hashed, err := auth.HashPassword(ctx, password)
		if err != nil {
			return err
		}
		record.Data["password"] = hashed
		return nil
	})

	// Bootstrap system
	if err := registry.BootstrapSystemCollections(); err != nil {
		slog.Error("Failed to bootstrap system collections", "error", err)
		os.Exit(1)
	}

	if err := registry.BootstrapRefreshTokensCollection(); err != nil {
		slog.Error("Failed to bootstrap refresh tokens collection", "error", err)
		os.Exit(1)
	}

	if err := registry.BootstrapUsersCollection(); err != nil {
		slog.Error("Failed to bootstrap users collection", "error", err)
		os.Exit(1)
	}

	// Sync system tables
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	systemCols := []string{"_collections", "_refresh_tokens", "users"}
	for _, name := range systemCols {
		col, ok := registry.GetCollection(name)
		if !ok || col == nil {
			slog.Error("Failed to find system collection in registry", "name", name)
			os.Exit(1)
		}
		if err := migration.SyncCollection(ctx, col); err != nil {
			slog.Error("Failed to sync system collection", "name", name, "error", err)
			os.Exit(1)
		}
		// Save system collection metadata to _collections table (except _collections itself to avoid recursion)
		if name != "_collections" {
			if err := registry.SaveCollection(ctx, col); err != nil {
				slog.Warn("Failed to save system collection metadata", "name", name, "error", err)
			}
		}
	}

	// Load dynamic collections from DB
	if err := registry.LoadFromDB(ctx); err != nil {
		slog.Error("Failed to load dynamic collections", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(executor, registry, store, hub, migration, cfg)
	handler := api.Chain(router,
		api.RecoveryMiddleware,
		api.LoggerMiddleware,
		api.AuthMiddleware(cfg.JWTSecret), // Added this line
		api.RequestIDMiddleware,
		api.CORSMiddleware,
	)

	srv := NewServer(cfg.Port, handler)

	return &App{
		Config:    cfg,
		DB:        database,
		Server:    srv,
		Registry:  registry,
		Migration: migration,
		Executor:  executor,
		Storage:   store,
		Hub:       hub,
	}
}

func (a *App) Run() {
	// Check if any users exist
	ctx := context.Background()
	var count int
	err := a.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err == nil && count == 0 {
		slog.Warn("⚠️  No users found!")
		slog.Warn("Create an admin user with: ./vault admin create --email <email> --password <password> --username <username>")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.Server.Start(); err != nil {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}

	if err := a.DB.Close(); err != nil {
		slog.Error("Database close failed", "error", err)
	}

	slog.Info("Vault stopped cleanly")
}
