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

const Version = "0.7.0"

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
	case "storage":
		runStorage()
	case "init":
		runInit()
	case "export":
		runExport()
	case "import":
		runImport()
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
	fmt.Println("  collection <subcommand>     Manage collections")
	fmt.Println("  storage <subcommand>        Manage storage")
	fmt.Println("  backup <subcommand>         Backup and restore operations")
	fmt.Println("  migrate <subcommand>        Database migration operations")
	fmt.Println("  export <format>             Export collections and data")
	fmt.Println("  import <format>             Import data from external sources")
	fmt.Println("  version                     Display version information")
	fmt.Println("  help, -h, --help            Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault init --email admin@example.com --username admin --password secret")
	fmt.Println("  vault serve --port 8090")
	fmt.Println("  vault admin create --email user@example.com --username user --password pass")
	fmt.Println("  vault backup create --output backup.zip")
	fmt.Println("  vault collection create --name posts --fields \"title:text,body:text\" --email admin@example.com --password secret")
	fmt.Println("  vault export json --output ./backup.json")
	fmt.Println("  vault import d1 wrangler-export/d1/homepage-db.sql")
	fmt.Println()
	fmt.Println("Use 'vault <command> -h' or 'vault <command> --help' for detailed command usage.")
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

	// Load base config (handles defaults, file, and env vars)
	cfg := core.LoadConfig(*configFile)

	// Apply flags overrides
	serveCmd.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "port":
			cfg.Port = *port
		case "dir":
			cfg.DataDir = *dir
		case "db-path":
			cfg.DBPath = *dbPath
		case "log-level":
			cfg.LogLevel = *logLevel
		case "log-format":
			cfg.LogFormat = *logFormat
		case "tls-cert":
			cfg.TLSCertPath = *tlsCert
		case "tls-key":
			cfg.TLSKeyPath = *tlsKey
		case "jwt-secret":
			cfg.JWTSecret = *jwtSecret
		case "cors-origins":
			cfg.CORSOrigins = *corsOrigins
		case "rate-limit":
			cfg.RateLimitPerMin = *rateLimit
		case "max-upload-size":
			if size := parseSize(*maxUploadSize); size > 0 {
				cfg.MaxFileUploadSize = size
			}
		}
	})

	runServer(cfg)
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
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	adminCmd := cli.NewAdminCommand(cfg)
	if err := adminCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Admin command failed", "error", err)
		os.Exit(1)
	}
}

func runServer(cfg *core.Config) {
	// Import server package
	app := server.NewApp(cfg)
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
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

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
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	collectionCmd := cli.NewCollectionCommand(cfg)
	if err := collectionCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Collection command failed", "error", err)
		os.Exit(1)
	}
}

func runStorage() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault storage <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  create --path PATH --file FILE")
		fmt.Println("  list [--path PATH] [--recursive]")
		fmt.Println("  get --path PATH --output FILE")
		fmt.Println("  delete --path PATH [--recursive] [--force]")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	storageCmd := cli.NewStorageCommand(cfg)
	if err := storageCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Storage command failed", "error", err)
		os.Exit(1)
	}
}

func runExport() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault export <format> [options]")
		fmt.Println("Formats:")
		fmt.Println("  json     Export collections and records as JSON")
		fmt.Println("  sql      Export schema and data as SQL statements")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	exportCmd := cli.NewExportCommand(cfg)
	if err := exportCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Export command failed", "error", err)
		os.Exit(1)
	}
}

func runImport() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: vault import <format> [options] <file>")
		fmt.Println("Formats:")
		fmt.Println("  sql      Import from generic SQL file")
		fmt.Println("  json     Import from JSON file")
		fmt.Println("  d1       Import from Cloudflare D1 SQL dump")
		os.Exit(1)
	}

	cfg := core.LoadConfig()
	if err := core.InitLogger(cfg.LogLevel, cfg.LogFormat, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	importCmd := cli.NewImportCommand(cfg)
	if err := importCmd.Run(os.Args[2:]); err != nil {
		slog.Error("Import command failed", "error", err)
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
