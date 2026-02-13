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
	case "backup":
		runBackup()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  vault serve [--port PORT] [--dir DIR]")
	fmt.Println("  vault admin <subcommand> [options]")
	fmt.Println("  vault backup <subcommand> [options]")
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
