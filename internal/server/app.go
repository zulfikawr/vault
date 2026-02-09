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
)

type App struct {
	Config    *core.Config
	DB        *sql.DB
	Server    *Server
	Registry  *db.SchemaRegistry
	Migration *db.MigrationEngine
	Executor  *db.Executor
}

func NewApp() *App {
	cfg := core.LoadConfig()
	core.InitLogger(cfg.LogLevel, cfg.LogFormat)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, cfg.DBPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	registry := db.NewSchemaRegistry(database)
	migration := db.NewMigrationEngine(database)
	executor := db.NewExecutor(database, registry)

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
	}

	router := api.NewRouter(executor, cfg)
	handler := api.Chain(router,
		api.RecoveryMiddleware,
		api.LoggerMiddleware,
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
	}
}

func (a *App) Run() {
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
