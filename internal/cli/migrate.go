package cli

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"time"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/models"
)

type MigrateCommand struct {
	config *core.Config
	db     *sql.DB
}

func NewMigrateCommand(config *core.Config) *MigrateCommand {
	return &MigrateCommand{config: config}
}

func (mc *MigrateCommand) Run(args []string) error {
	if len(args) < 1 {
		mc.printUsage()
		return fmt.Errorf("no migrate subcommand provided")
	}

	subcommand := args[0]

	// Check for help flags before processing subcommands
	if subcommand == "-h" || subcommand == "--help" {
		mc.printUsage()
		return nil
	}

	switch subcommand {
	case "sync":
		return mc.Sync(args[1:])
	case "status":
		return mc.Status(args[1:])
	default:
		mc.printUsage()
		return fmt.Errorf("unknown migrate subcommand: %s", subcommand)
	}
}

func (mc *MigrateCommand) Sync(args []string) error {
	cmd := flag.NewFlagSet("migrate sync", flag.ContinueOnError)
	collection := cmd.String("collection", "", "Specific collection to sync (optional)")
	verbose := cmd.Bool("verbose", false, "Verbose output")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Connect to database
	database, err := db.Connect(ctx, mc.config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	mc.db = database

	// Initialize schema registry and migration engine
	registry := db.NewSchemaRegistry(mc.db)
	migration := db.NewMigrationEngine(mc.db)

	// Bootstrap system collections
	if err := registry.BootstrapSystemCollections(); err != nil {
		return fmt.Errorf("failed to bootstrap system collections: %w", err)
	}

	if err := registry.BootstrapRefreshTokensCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap refresh tokens collection: %w", err)
	}

	if err := registry.BootstrapUsersCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap users collection: %w", err)
	}

	if err := registry.BootstrapAuditLogsCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap audit logs collection: %w", err)
	}

	// Load existing collections from database
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load collections from database: %w", err)
	}

	// Get all collections
	allCollections := registry.GetCollections()

	if len(allCollections) == 0 {
		fmt.Println("No collections found to migrate")
		return nil
	}

	// Filter by collection name if specified
	var collectionsToSync []*models.Collection
	if *collection != "" {
		found := false
		for _, col := range allCollections {
			if col.Name == *collection {
				collectionsToSync = append(collectionsToSync, col)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("collection '%s' not found", *collection)
		}
	} else {
		collectionsToSync = allCollections
	}

	// Sync collections
	fmt.Printf("Syncing %d collection(s)...\n\n", len(collectionsToSync))

	successCount := 0
	errorCount := 0

	for _, col := range collectionsToSync {
		if *verbose {
			fmt.Printf("Syncing collection: %s\n", col.Name)
		}

		if err := migration.SyncCollection(ctx, col); err != nil {
			fmt.Printf("✗ Failed to sync '%s': %v\n", col.Name, err)
			errorCount++
			continue
		}

		if *verbose {
			fmt.Printf("✓ Synced '%s' with %d fields\n", col.Name, len(col.Fields))
		} else {
			fmt.Printf("✓ %s\n", col.Name)
		}

		successCount++
	}

	fmt.Printf("\n")
	fmt.Printf("Migration Summary:\n")
	fmt.Printf("  Total: %d\n", len(collectionsToSync))
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Failed: %d\n", errorCount)

	if errorCount > 0 {
		return fmt.Errorf("migration completed with %d error(s)", errorCount)
	}

	fmt.Printf("\n✓ All collections synced successfully\n")
	return nil
}

func (mc *MigrateCommand) Status(args []string) error {
	cmd := flag.NewFlagSet("migrate status", flag.ContinueOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to database
	database, err := db.Connect(ctx, mc.config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	mc.db = database

	// Initialize schema registry
	registry := db.NewSchemaRegistry(mc.db)

	// Bootstrap system collections
	if err := registry.BootstrapSystemCollections(); err != nil {
		return fmt.Errorf("failed to bootstrap system collections: %w", err)
	}

	if err := registry.BootstrapRefreshTokensCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap refresh tokens collection: %w", err)
	}

	if err := registry.BootstrapUsersCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap users collection: %w", err)
	}

	if err := registry.BootstrapAuditLogsCollection(); err != nil {
		return fmt.Errorf("failed to bootstrap audit logs collection: %w", err)
	}

	// Load existing collections from database
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load collections from database: %w", err)
	}

	// Get all collections
	allCollections := registry.GetCollections()

	if len(allCollections) == 0 {
		fmt.Println("No collections found")
		return nil
	}

	fmt.Printf("Database Status:\n")
	fmt.Printf("  Database: %s\n", mc.config.DBPath)
	fmt.Printf("  Collections: %d\n\n", len(allCollections))

	fmt.Printf("%-30s %-15s %-10s\n", "Collection", "Type", "Fields")
	fmt.Println("-----------------------------------------------------------")

	for _, col := range allCollections {
		fmt.Printf("%-30s %-15s %-10d\n", col.Name, col.Type, len(col.Fields))
	}

	return nil
}

func (mc *MigrateCommand) printUsage() {
	fmt.Println("Usage: vault migrate <subcommand> [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  sync [--collection NAME] [--verbose]")
	fmt.Println("  status")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  sync     - Synchronize database schema with collections")
	fmt.Println("  status   - Show current database and collection status")
}
