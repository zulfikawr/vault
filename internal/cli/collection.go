package cli

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/models"
)

type CollectionCommand struct {
	config *core.Config
	db     *sql.DB
}

func NewCollectionCommand(config *core.Config) *CollectionCommand {
	return &CollectionCommand{config: config}
}

func (cc *CollectionCommand) Run(args []string) error {
	if len(args) < 1 {
		cc.printUsage()
		return fmt.Errorf("no collection subcommand provided")
	}

	subcommand := args[0]

	// Check for help flags before processing subcommands
	if subcommand == "-h" || subcommand == "--help" {
		cc.printUsage()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, cc.config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	cc.db = database

	switch subcommand {
	case "create":
		return cc.Create(ctx, args[1:])
	case "list":
		return cc.List(ctx, args[1:])
	case "get":
		return cc.Get(ctx, args[1:])
	case "delete":
		return cc.Delete(ctx, args[1:])
	default:
		cc.printUsage()
		return fmt.Errorf("unknown collection subcommand: %s", subcommand)
	}
}

func (cc *CollectionCommand) printUsage() {
	fmt.Println("Usage: vault collection <subcommand> [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  create --name NAME --fields FIELDS --email EMAIL --password PASSWORD")
	fmt.Println("  list --email EMAIL --password PASSWORD")
	fmt.Println("  get --name NAME --email EMAIL --password PASSWORD")
	fmt.Println("  delete --name NAME --email EMAIL --password PASSWORD [--force]")
	fmt.Println()
	fmt.Println("Fields format: name:type[,name:type,...]")
	fmt.Println("Field types: text, number, boolean, date, json")
}

func (cc *CollectionCommand) authenticateAdmin(ctx context.Context, email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password are required")
	}

	var storedPassword string
	err := cc.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("admin user not found")
		}
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if !auth.ComparePasswords(storedPassword, password) {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func (cc *CollectionCommand) Create(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	name := fs.String("name", "", "Collection name")
	fieldsStr := fs.String("fields", "", "Fields (comma-separated: name:type)")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" || *fieldsStr == "" || *email == "" || *password == "" {
		fmt.Println("Error: --name, --fields, --email, and --password are required")
		cc.printUsage()
		return fmt.Errorf("missing required flags")
	}

	if err := cc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fields, err := cc.parseFields(*fieldsStr)
	if err != nil {
		return fmt.Errorf("invalid fields: %w", err)
	}

	col := &models.Collection{
		ID:      "col_" + *name,
		Name:    *name,
		Type:    models.CollectionTypeBase,
		Fields:  fields,
		Created: time.Now().Format(time.RFC3339),
		Updated: time.Now().Format(time.RFC3339),
	}

	registry := db.NewSchemaRegistry(cc.db)
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	migration := db.NewMigrationEngine(cc.db)
	if err := migration.SyncCollection(ctx, col); err != nil {
		return fmt.Errorf("failed to sync collection: %w", err)
	}

	if err := registry.SaveCollection(ctx, col); err != nil {
		return fmt.Errorf("failed to save collection: %w", err)
	}

	// Log audit event
	db.LogAuditEvent(ctx, cc.db, "collection_created", col.Name, *email, map[string]any{
		"fields": len(fields),
		"type":   col.Type,
	})

	slog.Info("collection_created", "collection", *name, "fields", len(fields), "email", *email)
	fmt.Printf("✓ Collection '%s' created with %d fields\n", *name, len(fields))
	return nil
}

func (cc *CollectionCommand) List(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *email == "" || *password == "" {
		fmt.Println("Error: --email and --password are required")
		cc.printUsage()
		return fmt.Errorf("missing required flags")
	}

	if err := cc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	registry := db.NewSchemaRegistry(cc.db)
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	collections := registry.GetCollections()
	if len(collections) == 0 {
		fmt.Println("No collections found")
		return nil
	}

	fmt.Println("\nCollections:")
	fmt.Println("─────────────────────────────────────────────────────")
	for _, col := range collections {
		fmt.Printf("  %s (%s) - %d fields\n", col.Name, col.Type, len(col.Fields))
	}
	fmt.Println("─────────────────────────────────────────────────────")
	fmt.Printf("Total: %d collections\n\n", len(collections))

	return nil
}

func (cc *CollectionCommand) Get(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	name := fs.String("name", "", "Collection name")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" || *email == "" || *password == "" {
		fmt.Println("Error: --name, --email, and --password are required")
		cc.printUsage()
		return fmt.Errorf("missing required flags")
	}

	if err := cc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	registry := db.NewSchemaRegistry(cc.db)
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	col, ok := registry.GetCollection(*name)
	if !ok {
		return fmt.Errorf("collection '%s' not found", *name)
	}

	fmt.Printf("\nCollection: %s\n", col.Name)
	fmt.Printf("Type: %s\n", col.Type)
	fmt.Printf("Created: %s\n", col.Created)
	fmt.Printf("Updated: %s\n", col.Updated)
	fmt.Println("\nFields:")
	for _, field := range col.Fields {
		fmt.Printf("  - %s (%s) required=%v unique=%v\n", field.Name, field.Type, field.Required, field.Unique)
	}
	fmt.Println()

	return nil
}

func (cc *CollectionCommand) Delete(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	name := fs.String("name", "", "Collection name")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")
	force := fs.Bool("force", false, "Skip confirmation")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" || *email == "" || *password == "" {
		fmt.Println("Error: --name, --email, and --password are required")
		cc.printUsage()
		return fmt.Errorf("missing required flags")
	}

	if err := cc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	registry := db.NewSchemaRegistry(cc.db)
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	_, ok := registry.GetCollection(*name)
	if !ok {
		return fmt.Errorf("collection '%s' not found", *name)
	}

	if !*force {
		fmt.Printf("Are you sure you want to delete collection '%s'? (yes/no): ", *name)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "yes" {
			fmt.Println("Deletion cancelled")
			return nil
		}
	}

	if err := cc.deleteCollection(ctx, *name); err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	// Log audit event
	db.LogAuditEvent(ctx, cc.db, "collection_deleted", *name, *email, map[string]any{})

	slog.Info("collection_deleted", "collection", *name, "email", *email)
	fmt.Printf("✓ Collection '%s' deleted\n", *name)
	return nil
}

func (cc *CollectionCommand) deleteCollection(ctx context.Context, name string) error {
	tx, err := cc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", name)); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM _collections WHERE name = ?", name); err != nil {
		return err
	}

	return tx.Commit()
}

func (cc *CollectionCommand) parseFields(fieldsStr string) ([]models.Field, error) {
	var fields []models.Field
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		fieldParts := strings.Split(part, ":")
		if len(fieldParts) < 2 {
			return nil, fmt.Errorf("invalid field format: %s (expected name:type[:constraint1:constraint2])", part)
		}

		fieldName := strings.TrimSpace(fieldParts[0])
		fieldType := strings.TrimSpace(fieldParts[1])

		field := models.Field{
			Name: fieldName,
			Type: models.FieldType(fieldType),
		}

		// Parse constraints
		for i := 2; i < len(fieldParts); i++ {
			constraint := strings.TrimSpace(fieldParts[i])
			switch constraint {
			case "required":
				field.Required = true
			case "unique":
				field.Unique = true
			}
		}

		fields = append(fields, field)
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields provided")
	}

	return fields, nil
}
