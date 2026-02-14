package cli

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

type AdminCommand struct {
	config *core.Config
	db     *sql.DB
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

	switch subcommand {
	case "create":
		return ac.Create(args[1:])
	case "list":
		return ac.List(args[1:])
	case "delete":
		return ac.Delete(args[1:])
	case "reset-password":
		return ac.ResetPassword(args[1:])
	default:
		ac.printUsage()
		return fmt.Errorf("unknown admin subcommand: %s", subcommand)
	}
}

func (ac *AdminCommand) Create(args []string) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Bootstrap system collections
	registry := db.NewSchemaRegistry(ac.db)
	migration := db.NewMigrationEngine(ac.db)

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

	// Sync tables
	systemCols := []string{"_collections", "_refresh_tokens", "_audit_logs", "users"}
	for _, name := range systemCols {
		col, ok := registry.GetCollection(name)
		if !ok || col == nil {
			return fmt.Errorf("failed to find system collection: %s", name)
		}
		if err := migration.SyncCollection(ctx, col); err != nil {
			return fmt.Errorf("failed to sync collection %s: %w", name, err)
		}
		if name != "_collections" {
			if err := registry.SaveCollection(ctx, col); err != nil {
				return fmt.Errorf("failed to save collection %s: %w", name, err)
			}
		}
	}

	// Check if user already exists
	var existingID string
	err := ac.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ? OR username = ?", *email, *username).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("user with email or username already exists")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("database error: %w", err)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(ctx, *password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	userID := generateID()
	_, err = ac.db.ExecContext(ctx,
		"INSERT INTO users (id, username, email, password, created, updated) VALUES (?, ?, ?, ?, ?, ?)",
		userID, *username, *email, hashedPassword, time.Now(), time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Printf("✓ Admin user created successfully\n")
	fmt.Printf("  ID: %s\n", userID)
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", *username)

	return nil
}

func (ac *AdminCommand) List(args []string) error {
	cmd := flag.NewFlagSet("admin list", flag.ContinueOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := ac.db.QueryContext(ctx, "SELECT id, username, email, created FROM users ORDER BY created DESC")
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var admins []map[string]interface{}
	for rows.Next() {
		var id, username, email, createdStr string

		if err := rows.Scan(&id, &username, &email, &createdStr); err != nil {
			return fmt.Errorf("failed to scan user: %w", err)
		}

		admins = append(admins, map[string]interface{}{
			"id":       id,
			"username": username,
			"email":    email,
			"created":  createdStr,
		})
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating users: %w", err)
	}

	if len(admins) == 0 {
		fmt.Println("No admin users found")
		return nil
	}

	fmt.Printf("Total admins: %d\n\n", len(admins))
	fmt.Printf("%-20s %-20s %-30s %-25s\n", "ID", "Username", "Email", "Created")
	fmt.Println(strings.Repeat("-", 95))

	for _, admin := range admins {
		fmt.Printf("%-20s %-20s %-30s %-25s\n",
			admin["id"].(string)[:20],
			admin["username"].(string),
			admin["email"].(string),
			admin["created"].(string),
		)
	}

	return nil
}

func (ac *AdminCommand) Delete(args []string) error {
	cmd := flag.NewFlagSet("admin delete", flag.ContinueOnError)
	email := cmd.String("email", "", "Admin email to delete")
	force := cmd.Bool("force", false, "Skip confirmation prompt")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *email == "" {
		return fmt.Errorf("email is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if user exists
	var userID, username string
	err := ac.db.QueryRowContext(ctx, "SELECT id, username FROM users WHERE email = ?", *email).Scan(&userID, &username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user with email %s not found", *email)
	}
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// Confirm deletion
	if !*force {
		fmt.Printf("Are you sure you want to delete admin user '%s' (%s)? (yes/no): ", username, *email)
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
	result, err := ac.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	fmt.Printf("✓ Admin user deleted successfully\n")
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", username)

	return nil
}

func (ac *AdminCommand) ResetPassword(args []string) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if user exists
	var userID, username string
	err := ac.db.QueryRowContext(ctx, "SELECT id, username FROM users WHERE email = ?", *email).Scan(&userID, &username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user with email %s not found", *email)
	}
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(ctx, *newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	result, err := ac.db.ExecContext(ctx, "UPDATE users SET password = ?, updated = ? WHERE id = ?", hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	fmt.Printf("✓ Password reset successfully\n")
	fmt.Printf("  Email: %s\n", *email)
	fmt.Printf("  Username: %s\n", username)

	return nil
}

func (ac *AdminCommand) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  vault admin create --email EMAIL --password PASSWORD --username USERNAME")
	fmt.Println("  vault admin list")
	fmt.Println("  vault admin delete --email EMAIL [--force]")
	fmt.Println("  vault admin reset-password --email EMAIL --password PASSWORD")
}

func generateID() string {
	return fmt.Sprintf("usr_%d", time.Now().UnixNano())
}
