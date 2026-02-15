package cli

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/db"
	"golang.org/x/term"
)

func RunInit(args []string) error {
	// Check for help flags first
	if len(args) >= 1 && (args[0] == "-h" || args[0] == "--help") {
		printInitUsage()
		return nil
	}

	initCmd := flag.NewFlagSet("init", flag.ContinueOnError)
	email := initCmd.String("email", "", "Admin email address")
	username := initCmd.String("username", "", "Admin username")
	password := initCmd.String("password", "", "Admin password")
	dir := initCmd.String("dir", "./vault_data", "Data directory path")
	skipAdmin := initCmd.Bool("skip-admin", false, "Skip admin creation")
	force := initCmd.Bool("force", false, "Overwrite existing setup")

	if err := initCmd.Parse(args); err != nil {
		return err
	}

	fmt.Println("🚀 Initializing Vault project...")
	fmt.Println()

	// Check if already initialized
	dbPath := filepath.Join(*dir, "vault.db")
	if _, err := os.Stat(dbPath); err == nil && !*force {
		return fmt.Errorf("vault already initialized at %s (use --force to overwrite)", *dir)
	}

	// 1. Create directory structure
	if err := os.MkdirAll(filepath.Join(*dir, "storage"), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	fmt.Printf("✓ Created data directory: %s\n", *dir)

	// 2. Generate secure JWT secret
	jwtSecret, err := generateSecret(32)
	if err != nil {
		return fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	// 3. Create config.json
	if err := createConfigFile(jwtSecret); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	fmt.Println("✓ Generated config.json")

	// 4. Create .env.example
	if err := createEnvExample(); err != nil {
		return fmt.Errorf("failed to create .env.example: %w", err)
	}
	fmt.Println("✓ Created .env.example")

	// 5. Initialize database
	dbPath = filepath.Join(*dir, "vault.db")
	ctx := context.Background()
	database, err := db.Connect(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	fmt.Printf("✓ Initialized database: %s\n", dbPath)

	// 6. Create system collections
	if err := createSystemCollections(ctx, database); err != nil {
		return fmt.Errorf("failed to create system collections: %w", err)
	}
	fmt.Println("✓ Created system collections")
	fmt.Println()

	// 7. Create admin user
	if !*skipAdmin {
		adminEmail, adminUsername, adminPassword, err := getAdminCredentials(*email, *username, *password)
		if err != nil {
			return err
		}

		if err := createAdminUser(ctx, database, adminEmail, adminUsername, adminPassword); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
		fmt.Println("✓ Admin user created successfully")
	}

	fmt.Println()
	fmt.Println("🎉 Vault initialized successfully!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  Run: vault serve")
	fmt.Println("  Visit: http://localhost:8090/")

	return nil
}

func generateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func createConfigFile(jwtSecret string) error {
	config := fmt.Sprintf(`{
  "port": 8090,
  "data_dir": "./vault_data",
  "log_level": "INFO",
  "log_format": "text",
  "jwt_secret": "%s",
  "jwt_expiry_hours": 24,
  "refresh_token_days": 30,
  "tls_enabled": false,
  "tls_cert_path": "",
  "tls_key_path": "",
  "cors_origins": ["*"],
  "rate_limit_per_min": 100,
  "max_upload_size_mb": 10
}`, jwtSecret)

	return os.WriteFile("config.json", []byte(config), 0644)
}

func createEnvExample() error {
	envExample := `# Vault Environment Variables
# All settings can be overridden using VAULT_* prefix

# Server Configuration
VAULT_PORT=8090
VAULT_DATA_DIR=./vault_data

# Logging
VAULT_LOG_LEVEL=INFO
VAULT_LOG_FORMAT=text

# JWT Authentication
VAULT_JWT_SECRET=your-secret-key-here
VAULT_JWT_EXPIRY_HOURS=24
VAULT_REFRESH_TOKEN_DAYS=30

# TLS Configuration
VAULT_TLS_ENABLED=false
VAULT_TLS_CERT_PATH=
VAULT_TLS_KEY_PATH=

# CORS
VAULT_CORS_ORIGINS=*

# Rate Limiting
VAULT_RATE_LIMIT_PER_MIN=100

# File Upload
VAULT_MAX_UPLOAD_SIZE_MB=10
`
	return os.WriteFile(".env.example", []byte(envExample), 0644)
}

func createSystemCollections(ctx context.Context, database *sql.DB) error {
	registry := db.NewSchemaRegistry(database)
	migration := db.NewMigrationEngine(database)

	// Bootstrap system collections
	if err := registry.BootstrapSystemCollections(); err != nil {
		return err
	}

	if err := registry.BootstrapRefreshTokensCollection(); err != nil {
		return err
	}

	if err := registry.BootstrapUsersCollection(); err != nil {
		return err
	}

	if err := registry.BootstrapAuditLogsCollection(); err != nil {
		return err
	}

	// Sync tables
	systemCols := []string{"_collections", "_refresh_tokens", "_audit_logs", "users"}
	for _, name := range systemCols {
		col, ok := registry.GetCollection(name)
		if !ok || col == nil {
			return fmt.Errorf("failed to find system collection: %s", name)
		}
		if err := migration.SyncCollection(ctx, col); err != nil {
			return err
		}
	}

	return nil
}

func getAdminCredentials(email, username, password string) (string, string, string, error) {
	fmt.Println("👤 Create your first admin user:")

	// Get email
	if email == "" {
		fmt.Print("   Email: ")
		if _, err := fmt.Scanln(&email); err != nil {
			return "", "", "", fmt.Errorf("failed to read email: %w", err)
		}
	} else {
		fmt.Printf("   Email: %s\n", email)
	}

	// Get username
	if username == "" {
		fmt.Print("   Username: ")
		if _, err := fmt.Scanln(&username); err != nil {
			return "", "", "", fmt.Errorf("failed to read username: %w", err)
		}
	} else {
		fmt.Printf("   Username: %s\n", username)
	}

	// Get password
	if password == "" {
		fmt.Print("   Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", "", fmt.Errorf("failed to read password: %w", err)
		}
		password = string(passwordBytes)
		fmt.Println()
	} else {
		fmt.Println("   Password: ********")
	}

	// Validate
	if email == "" || username == "" || password == "" {
		return "", "", "", fmt.Errorf("email, username, and password are required")
	}

	if !strings.Contains(email, "@") {
		return "", "", "", fmt.Errorf("invalid email format")
	}

	if len(password) < 8 {
		return "", "", "", fmt.Errorf("password must be at least 8 characters")
	}

	fmt.Println()
	return email, username, password, nil
}

func createAdminUser(ctx context.Context, database *sql.DB, email, username, password string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Check if user already exists
	var existingID string
	err := database.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ? OR username = ?", email, username).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("user with email or username already exists")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("database error: %w", err)
	}

	hashedPassword, err := auth.HashPassword(ctx, password)
	if err != nil {
		return err
	}

	userID := generateID()
	query := `INSERT INTO users (id, username, email, password, created, updated) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	_, err = database.ExecContext(ctx, query, userID, username, email, hashedPassword, time.Now(), time.Now())
	return err
}

func generateID() string {
	return fmt.Sprintf("usr_%d", time.Now().UnixNano())
}

func printInitUsage() {
	fmt.Println("Usage: vault init [options]")
	fmt.Println("Options:")
	fmt.Println("  --email EMAIL               Admin email address")
	fmt.Println("  --username USERNAME         Admin username")
	fmt.Println("  --password PASSWORD         Admin password")
	fmt.Println("  --dir DIR                   Data directory (default: ./vault_data)")
	fmt.Println("  --skip-admin                Skip admin creation")
	fmt.Println("  --force                     Overwrite existing setup")
}
