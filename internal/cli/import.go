package cli

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
	"github.com/zulfikawr/vault/internal/service"
)

type ImportCommand struct {
	config            *core.Config
	db                *sql.DB
	registry          *db.SchemaRegistry
	migration         *db.MigrationEngine
	collectionService *service.CollectionService
	recordService     *service.RecordService
}

func NewImportCommand(config *core.Config) *ImportCommand {
	return &ImportCommand{config: config}
}

func (ic *ImportCommand) Run(args []string) error {
	if len(args) < 1 {
		ic.printUsage()
		return errors.NewError(400, "IMPORT_FORMAT_REQUIRED", "no import format specified")
	}

	// Check for help flags first
	if args[0] == "-h" || args[0] == "--help" {
		ic.printUsage()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Connect to database
	database, err := db.Connect(ctx, ic.config.DBPath)
	if err != nil {
		return errors.NewError(500, "DB_CONNECTION_FAILED", "failed to connect to database").WithDetails(map[string]any{"error": err.Error()})
	}
	defer database.Close()

	ic.db = database

	format := args[0]

	// D1 import doesn't need schema registry initially - it creates tables directly
	if format == "d1" {
		// Initialize minimal services for D1 import
		registry := db.NewSchemaRegistry(database)
		// Try to load, but don't fail if _collections doesn't exist yet
		_ = registry.LoadFromDB(ctx)
		ic.registry = registry
		ic.migration = db.NewMigrationEngine(database)
		return ic.importD1(ctx, args[1:])
	}

	// For other formats, we need full schema registry
	registry := db.NewSchemaRegistry(database)
	if err := registry.LoadFromDB(ctx); err != nil {
		return errors.NewError(500, "SCHEMA_LOAD_FAILED", "failed to load schema").WithDetails(map[string]any{"error": err.Error()})
	}
	ic.registry = registry

	migration := db.NewMigrationEngine(database)
	ic.migration = migration

	repo := db.NewRepository(database, registry)
	ic.recordService = service.NewRecordService(repo, nil)
	ic.collectionService = service.NewCollectionService(registry, migration)

	switch format {
	case "sql":
		return ic.importSQL(ctx, args[1:])
	case "json":
		return ic.importJSON(ctx, args[1:])
	default:
		ic.printUsage()
		return errors.NewError(400, "INVALID_IMPORT_FORMAT", fmt.Sprintf("unknown import format: %s", format))
	}
}

func (ic *ImportCommand) importSQL(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("import sql", flag.ContinueOnError)
	dryRun := cmd.Bool("dry-run", false, "Validate without importing")

	if err := cmd.Parse(args); err != nil {
		return errors.NewError(400, "INVALID_FLAGS", "failed to parse flags").WithDetails(map[string]any{"error": err.Error()})
	}

	if cmd.NArg() < 1 {
		return errors.NewError(400, "SQL_FILE_REQUIRED", "SQL file path required")
	}

	inputFile := cmd.Arg(0)

	// Check if file exists
	if _, err := os.Stat(inputFile); err != nil {
		return errors.NewError(404, "FILE_NOT_FOUND", fmt.Sprintf("input file not found: %s", inputFile))
	}

	fmt.Printf("📥 Importing SQL file: %s\n", inputFile)
	if *dryRun {
		fmt.Printf("⚠️  Dry-run mode: no changes will be made\n\n")
	}

	// Read and parse SQL file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return errors.NewError(500, "FILE_READ_FAILED", fmt.Sprintf("failed to read SQL file: %s", inputFile)).WithDetails(map[string]any{"error": err.Error()})
	}

	// Parse SQL statements
	statements := parseSQLStatements(string(content))

	fmt.Printf("Found %d SQL statements\n\n", len(statements))

	createCount := 0
	insertCount := 0
	errorCount := 0

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Detect statement type
		upperStmt := strings.ToUpper(stmt)

		if strings.HasPrefix(upperStmt, "CREATE TABLE") {
			createCount++
			tableName := extractTableName(stmt)

			if *dryRun {
				fmt.Printf("✓ Would create table: %s\n", tableName)
				continue
			}

			// Execute CREATE TABLE
			if _, err := ic.db.ExecContext(ctx, stmt); err != nil {
				fmt.Printf("✗ Failed to create table %s: %v\n", tableName, err)
				errorCount++
				continue
			}
			fmt.Printf("✓ Created table: %s\n", tableName)

		} else if strings.HasPrefix(upperStmt, "INSERT INTO") {
			insertCount++

			if *dryRun {
				if insertCount <= 5 {
					fmt.Printf("✓ Would insert record\n")
				}
				continue
			}

			// Execute INSERT
			if _, err := ic.db.ExecContext(ctx, stmt); err != nil {
				// Log but continue
				if insertCount <= 10 {
					fmt.Printf("⚠️  Insert warning: %v\n", err)
				}
				errorCount++
				continue
			}

			if insertCount%100 == 0 {
				fmt.Printf("  ... %d records inserted\n", insertCount)
			}
		}
	}

	fmt.Printf("\n")
	fmt.Printf("✅ Import complete!\n")
	fmt.Printf("   Tables created: %d\n", createCount)
	fmt.Printf("   Records inserted: %d\n", insertCount)
	if errorCount > 0 {
		fmt.Printf("   Errors: %d\n", errorCount)
	}

	return nil
}

func (ic *ImportCommand) importJSON(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("import json", flag.ContinueOnError)
	collection := cmd.String("collection", "", "Target collection name")
	dryRun := cmd.Bool("dry-run", false, "Validate without importing")

	if err := cmd.Parse(args); err != nil {
		return errors.NewError(400, "INVALID_FLAGS", "failed to parse flags").WithDetails(map[string]any{"error": err.Error()})
	}

	if cmd.NArg() < 1 {
		return errors.NewError(400, "JSON_FILE_REQUIRED", "JSON file path required")
	}

	inputFile := cmd.Arg(0)

	// Check if file exists
	if _, err := os.Stat(inputFile); err != nil {
		return errors.NewError(404, "FILE_NOT_FOUND", fmt.Sprintf("input file not found: %s", inputFile))
	}

	fmt.Printf("📥 Importing JSON file: %s\n", inputFile)
	if *dryRun {
		fmt.Printf("⚠️  Dry-run mode: no changes will be made\n\n")
	}

	// Read JSON file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return errors.NewError(500, "FILE_READ_FAILED", "failed to read JSON file").WithDetails(map[string]any{"error": err.Error()})
	}

	// Try to parse as Vault export format first
	var vaultExport map[string]any
	if err := json.Unmarshal(content, &vaultExport); err != nil {
		return errors.NewError(400, "JSON_PARSE_FAILED", "failed to parse JSON").WithDetails(map[string]any{"error": err.Error()})
	}

	// Check if it's a Vault export format
	if collectionsData, ok := vaultExport["collections"].(map[string]any); ok {
		return ic.importVaultJSON(ctx, collectionsData, *dryRun)
	}

	// Otherwise, treat as simple records array for a single collection
	if *collection == "" {
		return errors.NewError(400, "COLLECTION_REQUIRED", "--collection flag required for simple JSON import")
	}

	// Parse as array of records
	var records []map[string]any
	if err := json.Unmarshal(content, &records); err != nil {
		return fmt.Errorf("failed to parse records array: %w", err)
	}

	return ic.importRecordsToCollection(ctx, *collection, records, *dryRun)
}

func (ic *ImportCommand) importVaultJSON(ctx context.Context, collectionsData map[string]any, dryRun bool) error {
	totalCollections := 0
	totalRecords := 0
	errorCount := 0

	for colName, data := range collectionsData {
		totalCollections++
		fmt.Printf("⬇️  Importing: %s\n", colName)

		colData, ok := data.(map[string]any)
		if !ok {
			fmt.Printf("   ✗ Invalid collection data format\n")
			errorCount++
			continue
		}

		// Get schema
		schemaData, hasSchema := colData["schema"].(map[string]any)
		recordsData, hasRecords := colData["records"].([]any)

		if !hasRecords {
			fmt.Printf("   ✗ No records found\n")
			continue
		}

		// Create collection if schema exists
		if hasSchema && !dryRun {
			if err := ic.createCollectionFromSchema(colName, schemaData); err != nil {
				fmt.Printf("   ✗ Failed to create collection: %v\n", err)
				errorCount++
				continue
			}
			fmt.Printf("   ✓ Created collection schema\n")
		}

		// Import records
		recordCount := 0
		for _, recData := range recordsData {
			recMap, ok := recData.(map[string]any)
			if !ok {
				continue
			}

			if dryRun {
				recordCount++
				continue
			}

			// Create record
			if err := ic.createRecord(ctx, colName, recMap); err != nil {
				fmt.Printf("   ⚠️  Failed to create record: %v\n", err)
				errorCount++
				continue
			}
			recordCount++

			if recordCount%100 == 0 {
				fmt.Printf("   ... %d records imported\n", recordCount)
			}
		}

		fmt.Printf("   ✓ %d records\n", recordCount)
		totalRecords += recordCount
	}

	fmt.Printf("\n")
	fmt.Printf("✅ Import complete!\n")
	fmt.Printf("   Collections: %d\n", totalCollections)
	fmt.Printf("   Records: %d\n", totalRecords)
	if errorCount > 0 {
		fmt.Printf("   Errors: %d\n", errorCount)
	}

	return nil
}

func (ic *ImportCommand) importRecordsToCollection(ctx context.Context, collectionName string, records []map[string]any, dryRun bool) error {
	fmt.Printf("⬇️  Importing to collection: %s\n", collectionName)

	// Check if collection exists
	_, exists := ic.registry.GetCollection(collectionName)
	if !exists {
		return fmt.Errorf("collection '%s' does not exist", collectionName)
	}

	imported := 0
	errors := 0

	for _, recData := range records {
		if dryRun {
			imported++
			continue
		}

		if err := ic.createRecord(ctx, collectionName, recData); err != nil {
			fmt.Printf("   ⚠️  Failed to create record: %v\n", err)
			errors++
			continue
		}
		imported++

		if imported%100 == 0 {
			fmt.Printf("   ... %d records imported\n", imported)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("✅ Import complete!\n")
	fmt.Printf("   Records imported: %d\n", imported)
	if errors > 0 {
		fmt.Printf("   Errors: %d\n", errors)
	}

	return nil
}

func (ic *ImportCommand) importD1(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("import d1", flag.ContinueOnError)
	dryRun := cmd.Bool("dry-run", false, "Validate without importing")

	if err := cmd.Parse(args); err != nil {
		return errors.NewError(400, "INVALID_FLAGS", "failed to parse flags").WithDetails(map[string]any{"error": err.Error()})
	}

	if cmd.NArg() < 1 {
		return errors.NewError(400, "D1_FILE_REQUIRED", "D1 SQL dump file path required")
	}

	inputFile := cmd.Arg(0)

	// Check if file exists
	if _, err := os.Stat(inputFile); err != nil {
		return errors.NewError(404, "FILE_NOT_FOUND", fmt.Sprintf("input file not found: %s", inputFile))
	}

	fmt.Printf("📥 Importing D1 dump: %s\n", inputFile)
	if *dryRun {
		fmt.Printf("⚠️  Dry-run mode: no changes will be made\n\n")
	}

	// Read SQL file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return errors.NewError(500, "FILE_READ_FAILED", "failed to read SQL file").WithDetails(map[string]any{"error": err.Error()})
	}

	// Parse and import D1 SQL - don't need registry for raw SQL import
	return ic.parseAndImportD1SQL(ctx, string(content), *dryRun)
}

func (ic *ImportCommand) parseAndImportD1SQL(ctx context.Context, sqlContent string, dryRun bool) error {
	// Parse SQL statements
	statements := parseSQLStatements(sqlContent)

	fmt.Printf("Found %d SQL statements\n\n", len(statements))

	tablesCreated := 0
	recordsInserted := 0
	errorCount := 0

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		upperStmt := strings.ToUpper(stmt)

		if strings.HasPrefix(upperStmt, "CREATE TABLE") {
			tableName := extractTableName(stmt)

			// Skip system tables
			if strings.HasPrefix(tableName, "_") || tableName == "sqlite_sequence" {
				continue
			}

			tablesCreated++

			if dryRun {
				fmt.Printf("✓ Would create table: %s\n", tableName)
				continue
			}

			// Check if table exists
			var exists string
			err := ic.db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&exists)

			if err == sql.ErrNoRows {
				// Create table
				if _, err := ic.db.ExecContext(ctx, stmt); err != nil {
					fmt.Printf("✗ Failed to create table %s: %v\n", tableName, err)
					errorCount++
					continue
				}
				fmt.Printf("✓ Created table: %s\n", tableName)
			} else if err != nil {
				fmt.Printf("⚠️  Error checking table %s: %v\n", tableName, err)
			} else {
				fmt.Printf("✓ Table exists: %s\n", tableName)
			}

		} else if strings.HasPrefix(upperStmt, "INSERT INTO") {
			if dryRun {
				recordsInserted++
				continue
			}

			// Execute INSERT (ignore duplicates)
			if _, err := ic.db.ExecContext(ctx, stmt); err != nil {
				// Try with OR REPLACE
				stmt = strings.Replace(stmt, "INSERT INTO", "INSERT OR REPLACE INTO", 1)
				if _, err := ic.db.ExecContext(ctx, stmt); err != nil {
					errorCount++
					continue
				}
			}
			recordsInserted++

			if recordsInserted%500 == 0 {
				fmt.Printf("  ... %d records imported\n", recordsInserted)
			}
		}
	}

	fmt.Printf("\n")
	fmt.Printf("✅ D1 Import complete!\n")
	fmt.Printf("   Tables: %d\n", tablesCreated)
	fmt.Printf("   Records: %d\n", recordsInserted)
	if errorCount > 0 {
		fmt.Printf("   Errors: %d\n", errorCount)
	}

	// Register tables as Vault collections
	fmt.Printf("\n📋 Registering collections in Vault...\n")

	// Bootstrap _collections table first
	if err := ic.registry.BootstrapSystemCollections(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to bootstrap system collections: %v\n", err)
	}

	if err := ic.registerImportedTables(ctx); err != nil {
		fmt.Printf("⚠️  Warning: Failed to register some collections: %v\n", err)
	}

	return nil
}

func (ic *ImportCommand) registerImportedTables(ctx context.Context) error {
	// Get all tables from database - include all user tables
	rows, err := ic.db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite\\_%' ESCAPE '\\'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		// Skip system tables
		if strings.HasPrefix(name, "_") {
			continue
		}
		tables = append(tables, name)
	}

	fmt.Printf("   Found %d user tables: %v\n", len(tables), tables)

	registered := 0
	for _, tableName := range tables {
		// Check if already registered
		_, exists := ic.registry.GetCollection(tableName)
		if exists {
			continue
		}

		// Get table schema
		fieldRows, err := ic.db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", tableName))
		if err != nil {
			fmt.Printf("   ⚠️  Failed to get schema for %s: %v\n", tableName, err)
			continue
		}

		var fields []models.Field
		for fieldRows.Next() {
			var cid int
			var name, dtype string
			var notnull, pk int
			var dfltValue sql.NullString
			if err := fieldRows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk); err != nil {
				continue
			}

			// Skip standard fields
			if name == "id" || name == "created" || name == "updated" {
				continue
			}

			field := models.Field{
				Name:     name,
				Type:     sqlTypeToVaultType(dtype),
				Required: notnull == 1,
			}
			fields = append(fields, field)
		}
		fieldRows.Close()

		if len(fields) == 0 {
			continue
		}

		// Create collection
		col := &models.Collection{
			ID:      fmt.Sprintf("col_%s", tableName),
			Name:    tableName,
			Type:    models.CollectionTypeBase,
			Fields:  fields,
			Created: time.Now().Format(time.RFC3339),
			Updated: time.Now().Format(time.RFC3339),
		}

		// Save to registry and DB
		if err := ic.registry.SaveCollection(ctx, col); err != nil {
			fmt.Printf("   ⚠️  Failed to register %s: %v\n", tableName, err)
			continue
		}

		registered++
		fmt.Printf("   ✓ Registered: %s (%d fields)\n", tableName, len(fields))
	}

	if registered > 0 {
		fmt.Printf("   Registered %d collections\n", registered)
	} else {
		fmt.Printf("   No new collections to register\n")
	}

	return nil
}

func sqlTypeToVaultType(sqlType string) models.FieldType {
	sqlType = strings.ToUpper(sqlType)
	switch {
	case strings.Contains(sqlType, "INT") || strings.Contains(sqlType, "BOOL"):
		return models.FieldTypeBool
	case strings.Contains(sqlType, "REAL") || strings.Contains(sqlType, "FLOAT") || strings.Contains(sqlType, "DOUBLE"):
		return models.FieldTypeNumber
	default:
		return models.FieldTypeText
	}
}

func (ic *ImportCommand) createCollectionFromSchema(name string, schemaData map[string]any) error {
	// Check if already exists
	_, exists := ic.registry.GetCollection(name)
	if exists {
		return nil
	}

	// Build collection from schema
	fields := []models.Field{}

	if fieldsData, ok := schemaData["fields"].([]any); ok {
		for _, f := range fieldsData {
			fMap, ok := f.(map[string]any)
			if !ok {
				continue
			}

			field := models.Field{
				Name:     getString(fMap, "name"),
				Type:     models.FieldType(getString(fMap, "type")),
				Required: getBool(fMap, "required"),
				Unique:   getBool(fMap, "unique"),
			}
			fields = append(fields, field)
		}
	}

	col := &models.Collection{
		ID:      fmt.Sprintf("col_%s", name),
		Name:    name,
		Type:    models.CollectionTypeBase,
		Fields:  fields,
		Created: time.Now().Format(time.RFC3339),
		Updated: time.Now().Format(time.RFC3339),
	}

	return ic.collectionService.CreateCollection(context.Background(), col)
}

func (ic *ImportCommand) createRecord(ctx context.Context, collectionName string, data map[string]any) error {
	// Extract standard fields
	id, hasID := data["id"].(string)

	if !hasID || id == "" {
		id = fmt.Sprintf("rec_%d", time.Now().UnixNano())
	}

	// Build record data
	recordData := make(map[string]any)
	for k, v := range data {
		if k != "id" && k != "created" && k != "updated" {
			recordData[k] = v
		}
	}

	// Add ID back for repository
	recordData["id"] = id

	// Use repository to insert
	repo := db.NewRepository(ic.db, ic.registry)
	_, err := repo.CreateRecord(ctx, collectionName, recordData)
	return err
}

func parseSQLStatements(content string) []string {
	var statements []string
	var current strings.Builder
	inMultiline := false

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}

		// Check for multiline statements
		if inMultiline {
			current.WriteString(" ")
			current.WriteString(line)
			if strings.Contains(line, ";") {
				inMultiline = false
				stmt := strings.TrimSpace(current.String())
				stmt = strings.TrimSuffix(stmt, ";")
				statements = append(statements, stmt)
				current.Reset()
			}
			continue
		}

		// Check if line ends with semicolon
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			statements = append(statements, strings.TrimSpace(strings.TrimSuffix(line, ";")))
		} else if strings.TrimSpace(line) != "" {
			current.WriteString(line)
			if !strings.HasSuffix(line, ")") {
				inMultiline = true
			}
		}
	}

	// Handle last statement without semicolon
	if current.Len() > 0 {
		statements = append(statements, strings.TrimSpace(current.String()))
	}

	return statements
}

func extractTableName(createStmt string) string {
	// Simple extraction: CREATE TABLE name (...)
	start := strings.Index(strings.ToUpper(createStmt), "TABLE") + 6
	if start <= 5 {
		return ""
	}

	rest := strings.TrimSpace(createStmt[start:])
	end := strings.IndexAny(rest, " (")
	if end == -1 {
		return rest
	}

	return strings.TrimSpace(rest[:end])
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBool(m map[string]any, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func (ic *ImportCommand) printUsage() {
	fmt.Println("Usage: vault import <format> [options] <file>")
	fmt.Println()
	fmt.Println("Formats:")
	fmt.Println("  sql      Import from generic SQL file")
	fmt.Println("  json     Import from JSON file (Vault export or records array)")
	fmt.Println("  d1       Import from Cloudflare D1 SQL dump")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --dry-run              Validate without making changes")
	fmt.Println("  --collection NAME      Target collection (for simple JSON import)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault import d1 wrangler-export/d1/homepage-db.sql")
	fmt.Println("  vault import d1 dump.sql --dry-run")
	fmt.Println("  vault import json users.json --collection users")
	fmt.Println("  vault import sql schema.sql")
}
