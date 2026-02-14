package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/zulfikawr/vault/internal/cli"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/server"
)

const Version = "0.5.1"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	// Handle global help flags
	if command == "-h" || command == "--help" || command == "help" {
		printUsage()
		os.Exit(0)
	}

	switch command {
	case "serve":
		runServe()
	case "admin":
		runAdmin()
	case "backup":
		runBackup()
	case "migrate":
		runMigrate()
	case "collection":
		runCollection()
	case "init":
		runInit()
	case "version", "-v", "--version":
		runVersion()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Vault - Self-contained Backend Framework")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vault <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init                        Initialize new Vault project")
	fmt.Println("  serve                       Start the HTTP server")
	fmt.Println("  admin <subcommand>          Manage admin users")
	fmt.Println("  backup <subcommand>         Backup and restore operations")
	fmt.Println("  migrate <subcommand>        Database migration operations")
	fmt.Println("  version                     Display version information")
	fmt.Println("  help, -h, --help            Show this help message")
	fmt.Println()
	fmt.Println("Init Options:")
	fmt.Println("  --email EMAIL               Admin email address")
	fmt.Println("  --username USERNAME         Admin username")
	fmt.Println("  --password PASSWORD         Admin password")
	fmt.Println("  --dir DIR                   Data directory (default: ./vault_data)")
	fmt.Println("  --skip-admin                Skip admin creation")
	fmt.Println("  --force                     Overwrite existing setup")
	fmt.Println()
	fmt.Println("Serve Options:")
	fmt.Println("  --port PORT                 Server port (default: 8090)")
	fmt.Println("  --dir DIR                   Data directory (default: ./vault_data)")
	fmt.Println("  --db-path PATH              Database path")
	fmt.Println("  --log-level LEVEL           Log level (DEBUG/INFO/WARN/ERROR)")
	fmt.Println("  --log-format FORMAT         Log format (text/json)")
	fmt.Println("  --tls-cert PATH             TLS certificate path")
	fmt.Println("  --tls-key PATH              TLS key path")
	fmt.Println("  --jwt-secret SECRET         JWT secret")
	fmt.Println("  --cors-origins ORIGINS      CORS origins (comma-separated)")
	fmt.Println("  --rate-limit NUM            Rate limit per minute")
	fmt.Println("  --max-upload-size SIZE      Max upload size (e.g., 10MB, 1GB)")
	fmt.Println("  --config FILE               Config file path")
	fmt.Println()
	fmt.Println("Admin Subcommands:")
	fmt.Println("  create                      Create new admin user")
	fmt.Println("  list                        List all admin users")
	fmt.Println("  delete                      Delete admin user")
	fmt.Println("  reset-password              Reset admin password")
	fmt.Println()
	fmt.Println("Backup Subcommands:")
	fmt.Println("  create                      Create backup")
	fmt.Println("  list                        List all backups")
	fmt.Println("  restore                     Restore from backup")
	fmt.Println()
	fmt.Println("Migrate Subcommands:")
	fmt.Println("  sync                        Synchronize database schema")
	fmt.Println("  status                      Show migration status")
	fmt.Println()
	fmt.Println("Collection Subcommands:")
	fmt.Println("  create                      Create new collection")
	fmt.Println("  list                        List all collections")
	fmt.Println("  get                         Get collection details")
	fmt.Println("  delete                      Delete collection")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault init --email admin@example.com --username admin --password secret")
	fmt.Println("  vault serve --port 8090")
	fmt.Println("  vault admin create --email user@example.com --username user --password pass")
	fmt.Println("  vault backup create --output backup.zip")
	fmt.Println("  vault collection create --name posts --fields \"title:text,body:text\" --email admin@example.com --password secret")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/zulfikawr/vault")
}

func runServe() {
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	port := serveCmd.Int("port", 8090, "Server port")
	dir := serveCmd.String("dir", "./vault_data", "Data directory")
	dbPath := serveCmd.String("db-path", "", "Database path (default: {dir}/vault.db)")
	logLevel := serveCmd.String("log-level", "", "Log level (DEBUG/INFO/WARN/ERROR)")
	logFormat := serveCmd.String("log-format", "", "Log format (text/json)")
	tlsCert := serveCmd.String("tls-cert", "", "TLS certificate path")
	tlsKey := serveCmd.String("tls-key", "", "TLS key path")
	jwtSecret := serveCmd.String("jwt-secret", "", "JWT secret")
	corsOrigins := serveCmd.String("cors-origins", "", "CORS origins (comma-separated)")
	rateLimit := serveCmd.Int("rate-limit", 0, "Rate limit per minute")
	maxUploadSize := serveCmd.String("max-upload-size", "", "Max upload size (e.g., 10MB, 1GB)")
	configFile := serveCmd.String("config", "", "Config file path")

	if err := serveCmd.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables for config loading
	os.Setenv("VAULT_PORT", fmt.Sprintf("%d", *port))
	os.Setenv("VAULT_DATA_DIR", *dir)

	if *dbPath != "" {
		os.Setenv("VAULT_DB_PATH", *dbPath)
	}
	if *logLevel != "" {
		os.Setenv("VAULT_LOG_LEVEL", *logLevel)
	}
	if *logFormat != "" {
		os.Setenv("VAULT_LOG_FORMAT", *logFormat)
	}
	if *tlsCert != "" {
		os.Setenv("VAULT_TLS_CERT_PATH", *tlsCert)
	}
	if *tlsKey != "" {
		os.Setenv("VAULT_TLS_KEY_PATH", *tlsKey)
	}
	if *jwtSecret != "" {
		os.Setenv("VAULT_JWT_SECRET", *jwtSecret)
	}
	if *corsOrigins != "" {
		os.Setenv("VAULT_CORS_ORIGINS", *corsOrigins)
	}
	if *rateLimit > 0 {
		os.Setenv("VAULT_RATE_LIMIT_PER_MIN", fmt.Sprintf("%d", *rateLimit))
	}
	if *maxUploadSize != "" {
		bytes := parseSize(*maxUploadSize)
		os.Setenv("VAULT_MAX_FILE_UPLOAD_SIZE", fmt.Sprintf("%d", bytes))
	}
	if *configFile != "" {
		os.Setenv("VAULT_CONFIG_FILE", *configFile)
	}

	runServer()
}

func runAdmin() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault admin <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  create --email EMAIL --password PASSWORD --username USERNAME")
		fmt.Println("  list")
		fmt.Println("  delete --email EMAIL [--force]")
		fmt.Println("  reset-password --email EMAIL --password PASSWORD")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	core.InitLogger(cfg.LogLevel, cfg.LogFormat)

	adminCmd := cli.NewAdminCommand(cfg)
	if err := adminCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Admin command failed", "error", err)
		os.Exit(1)
	}
}

func runServer() {
	// Import server package
	app := server.NewApp()
	app.Run()
}

func runBackup() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault backup <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  create [--output FILE]")
		fmt.Println("  list")
		fmt.Println("  restore --input FILE [--force]")
		os.Exit(1)
	}

	cfg := core.LoadConfig()

	backupCmd := cli.NewBackupCommand(cfg.DataDir, cfg.DBPath)
	if err := backupCmd.Run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runVersion() {
	fmt.Printf("Vault version %s\n", Version)
}

func runMigrate() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault migrate <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  sync [--collection NAME] [--verbose]")
		fmt.Println("  status")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	core.InitLogger(cfg.LogLevel, cfg.LogFormat)

	migrateCmd := cli.NewMigrateCommand(cfg)
	if err := migrateCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Migrate command failed", "error", err)
		os.Exit(1)
	}
}

func runInit() {
	if err := cli.RunInit(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCollection() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault collection <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  create --name NAME --fields FIELDS --email EMAIL --password PASSWORD")
		fmt.Println("  list --email EMAIL --password PASSWORD")
		fmt.Println("  get --name NAME --email EMAIL --password PASSWORD")
		fmt.Println("  delete --name NAME --email EMAIL --password PASSWORD [--force]")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	core.InitLogger(cfg.LogLevel, cfg.LogFormat)

	collectionCmd := cli.NewCollectionCommand(cfg)
	if err := collectionCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Collection command failed", "error", err)
		os.Exit(1)
	}
}

func parseSize(size string) int64 {
	if size == "" {
		return 0
	}

	// Parse size with suffix (e.g., "10MB", "1GB")
	var num int64
	var suffix string

	_, err := fmt.Sscanf(size, "%d%s", &num, &suffix)
	if err != nil {
		return 0
	}

	switch suffix {
	case "B":
		return num
	case "KB":
		return num * 1024
	case "MB":
		return num * 1024 * 1024
	case "GB":
		return num * 1024 * 1024 * 1024
	default:
		return 0
	}
}
