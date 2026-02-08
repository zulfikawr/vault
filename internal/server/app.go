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
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

type App struct {
	Config *core.Config
	DB     *sql.DB
	Server *Server
}

func NewApp() *App {
	cfg := core.LoadConfig()
	core.InitLogger(cfg.LogLevel, cfg.LogFormat)

	database, err := db.Connect(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter()
	handler := api.Chain(router,
		api.RecoveryMiddleware,
		api.LoggerMiddleware,
		api.RequestIDMiddleware,
		api.CORSMiddleware,
	)

	srv := NewServer(cfg.Port, handler)

	return &App{
		Config: cfg,
		DB:     database,
		Server: srv,
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