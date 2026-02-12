package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		runServe()
	case "admin":
		runAdmin()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  vault serve [--port PORT] [--dir DIR]")
	fmt.Println("  vault admin create --email EMAIL --password PASSWORD --username USERNAME")
}

func runServe() {
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	port := serveCmd.Int("port", 8090, "Server port")
	dir := serveCmd.String("dir", "./vault_data", "Data directory")

	serveCmd.Parse(os.Args[2:])

	os.Setenv("VAULT_PORT", fmt.Sprintf("%d", *port))
	os.Setenv("VAULT_DATA_DIR", *dir)

	// Import and run the server
	runServer()
}

func runAdmin() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault admin create --email EMAIL --password PASSWORD --username USERNAME")
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "create":
		createAdmin()
	default:
		fmt.Printf("Unknown admin subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func createAdmin() {
	adminCmd := flag.NewFlagSet("create", flag.ExitOnError)
	email := adminCmd.String("email", "", "Admin email")
	password := adminCmd.String("password", "", "Admin password")
	username := adminCmd.String("username", "", "Admin username")

	adminCmd.Parse(os.Args[3:])

	if *email == "" || *password == "" || *username == "" {
		fmt.Println("Error: email, password, and username are required")
		fmt.Println("Usage: vault admin create --email EMAIL --password PASSWORD --username USERNAME")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	core.InitLogger(cfg.LogLevel, cfg.LogFormat)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, cfg.DBPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Bootstrap system collections
	registry := db.NewSchemaRegistry(database)
	migration := db.NewMigrationEngine(database)

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

	// Sync tables
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
		// Save collection metadata to _collections table
		if name != "_collections" {
			if err := registry.SaveCollection(ctx, col); err != nil {
				slog.Error("Failed to save collection metadata", "name", name, "error", err)
				os.Exit(1)
			}
		}
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(ctx, *password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		os.Exit(1)
	}

	// Insert user
	_, err = database.ExecContext(ctx,
		"INSERT INTO users (id, username, email, password, verified, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)",
		generateID(), *username, *email, hashedPassword, true, time.Now(), time.Now(),
	)
	if err != nil {
		slog.Error("Failed to create admin user", "error", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Admin user created successfully\n")
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", *username)
}

func generateID() string {
	return fmt.Sprintf("usr_%d", time.Now().UnixNano())
}

func runServer() {
	// Import server package
	app := server.NewApp()
	app.Run()
}
