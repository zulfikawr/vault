package cli

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/service"
)

type AdminCommand struct {
	config            *core.Config
	db                *sql.DB
	recordService     *service.RecordService
	collectionService *service.CollectionService
}

func NewAdminCommand(config *core.Config) *AdminCommand {
	return &AdminCommand{config: config}
}

func (ac *AdminCommand) Run(args []string) error {
	if len(args) < 1 {
		ac.printUsage()
		return fmt.Errorf("no admin subcommand provided")
	}

	subcommand := args[0]

	// Check for help flags before processing subcommands
	if subcommand == "-h" || subcommand == "--help" {
		ac.printUsage()
		return nil
	}

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, ac.config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	ac.db = database

	registry := db.NewSchemaRegistry(database)
	migration := db.NewMigrationEngine(database)
	repo := db.NewRepository(database, registry)

	ac.collectionService = service.NewCollectionService(registry, migration)
	ac.recordService = service.NewRecordService(repo, nil)

	// Initialize system and hooks
	service.RegisterAuthHooks()
	// Note: We don't call InitSystem here automatically for all commands,
	// but Create command might need it. Actually, Create command does bootstrap.
	// Let's defer InitSystem to Create command or do it here if it's safe.
	// Doing it here ensures schema is ready for List/Delete too.
	if err := ac.collectionService.InitSystem(ctx); err != nil {
		return fmt.Errorf("failed to initialize system: %w", err)
	}

	// Load dynamic collections
	if err := registry.LoadFromDB(ctx); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	switch subcommand {
	case "create":
		return ac.Create(ctx, args[1:])
	case "list":
		return ac.List(ctx, args[1:])
	case "delete":
		return ac.Delete(ctx, args[1:])
	case "reset-password":
		return ac.ResetPassword(ctx, args[1:])
	default:
		ac.printUsage()
		return fmt.Errorf("unknown admin subcommand: %s", subcommand)
	}
}

func (ac *AdminCommand) Create(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("admin create", flag.ContinueOnError)
	email := cmd.String("email", "", "Admin email")
	password := cmd.String("password", "", "Admin password")
	username := cmd.String("username", "", "Admin username")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *email == "" || *password == "" || *username == "" {
		return fmt.Errorf("email, password, and username are required")
	}

	// Check if user exists (email)
	records, _, err := ac.recordService.ListRecords(ctx, "users", db.QueryParams{Filter: fmt.Sprintf("email = '%s'", *email)})
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if len(records) > 0 {
		return fmt.Errorf("user with email %s already exists", *email)
	}

	// Check if user exists (username)
	records, _, err = ac.recordService.ListRecords(ctx, "users", db.QueryParams{Filter: fmt.Sprintf("username = '%s'", *username)})
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if len(records) > 0 {
		return fmt.Errorf("user with username %s already exists", *username)
	}

	// Create user
	data := map[string]any{
		"email":    *email,
		"username": *username,
		"password": *password, // Will be hashed by hook
	}

	record, err := ac.recordService.CreateRecord(ctx, "users", data)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Printf("✓ Admin user created successfully\n")
	fmt.Printf("  ID: %s\n", record.ID)
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", *username)

	return nil
}

func (ac *AdminCommand) List(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("admin list", flag.ContinueOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	records, _, err := ac.recordService.ListRecords(ctx, "users", db.QueryParams{
		Sort: "-created",
	})
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	if len(records) == 0 {
		fmt.Println("No admin users found")
		return nil
	}

	fmt.Printf("Total admins: %d\n\n", len(records))
	fmt.Printf("%-20s %-20s %-30s %-25s\n", "ID", "Username", "Email", "Created")
	fmt.Println(strings.Repeat("-", 95))

	for _, record := range records {
		createdStr := record.GetString("created")
		// Format if needed, or just print string if it's already formatted by service/repo
		// Repository usually returns time.Time or string depending on driver.
		// SQLite driver might return string or time.Time.
		// record.GetString handles type assertion.

		fmt.Printf("%-20s %-20s %-30s %-25s\n",
			record.ID[:20], // Assuming ID is long enough
			record.GetString("username"),
			record.GetString("email"),
			createdStr,
		)
	}

	return nil
}

func (ac *AdminCommand) Delete(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("admin delete", flag.ContinueOnError)
	email := cmd.String("email", "", "Admin email to delete")
	force := cmd.Bool("force", false, "Skip confirmation prompt")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *email == "" {
		return fmt.Errorf("email is required")
	}

	// Find user
	records, _, err := ac.recordService.ListRecords(ctx, "users", db.QueryParams{Filter: fmt.Sprintf("email = '%s'", *email)})
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("user with email %s not found", *email)
	}
	user := records[0]

	// Confirm deletion
	if !*force {
		fmt.Printf("Are you sure you want to delete admin user '%s' (%s)? (yes/no): ", user.GetString("username"), *email)
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		if strings.ToLower(response) != "yes" {
			fmt.Println("Deletion cancelled")
			return nil
		}
	}

	// Delete user
	if err := ac.recordService.DeleteRecord(ctx, "users", user.ID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	fmt.Printf("✓ Admin user deleted successfully\n")
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", user.GetString("username"))

	return nil
}

func (ac *AdminCommand) ResetPassword(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("admin reset-password", flag.ContinueOnError)
	email := cmd.String("email", "", "Admin email")
	newPassword := cmd.String("password", "", "New password")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *email == "" {
		return fmt.Errorf("email is required")
	}

	if *newPassword == "" {
		return fmt.Errorf("password is required")
	}

	if len(*newPassword) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Find user
	records, _, err := ac.recordService.ListRecords(ctx, "users", db.QueryParams{Filter: fmt.Sprintf("email = '%s'", *email)})
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("user with email %s not found", *email)
	}
	user := records[0]

	// Update password
	// Hook will hash it
	data := map[string]any{
		"password": *newPassword,
	}

	if _, err := ac.recordService.UpdateRecord(ctx, "users", user.ID, data); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	fmt.Printf("✓ Password reset successfully\n")
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", user.GetString("username"))

	return nil
}

func (ac *AdminCommand) printUsage() {
	fmt.Println("Usage: vault admin <subcommand> [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  create --email EMAIL --password PASSWORD --username USERNAME")
	fmt.Println("  list")
	fmt.Println("  delete --email EMAIL [--force]")
	fmt.Println("  reset-password --email EMAIL --password PASSWORD")
}
