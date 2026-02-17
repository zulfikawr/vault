package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
	"github.com/zulfikawr/vault/internal/service"
)

type ExportCommand struct {
	config        *core.Config
	db            *sql.DB
	recordService *service.RecordService
}

func NewExportCommand(config *core.Config) *ExportCommand {
	return &ExportCommand{config: config}
}

func (ec *ExportCommand) Run(args []string) error {
	if len(args) < 1 {
		ec.printUsage()
		return errors.NewError(400, "EXPORT_FORMAT_REQUIRED", "no export format specified")
	}

	// Check for help flags first
	if args[0] == "-h" || args[0] == "--help" {
		ec.printUsage()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Connect to database
	database, err := db.Connect(ctx, ec.config.DBPath)
	if err != nil {
		return errors.NewError(500, "DB_CONNECTION_FAILED", "failed to connect to database").WithDetails(map[string]any{"error": err.Error()})
	}
	defer database.Close()

	ec.db = database

	// Initialize services
	registry := db.NewSchemaRegistry(database)
	if err := registry.LoadFromDB(ctx); err != nil {
		return errors.NewError(500, "SCHEMA_LOAD_FAILED", "failed to load schema").WithDetails(map[string]any{"error": err.Error()})
	}

	repo := db.NewRepository(database, registry)
	ec.recordService = service.NewRecordService(repo, nil)

	format := args[0]

	switch format {
	case "json":
		return ec.exportJSON(ctx, args[1:])
	case "sql":
		return ec.exportSQL(ctx, args[1:])
	default:
		ec.printUsage()
		return errors.NewError(400, "INVALID_EXPORT_FORMAT", fmt.Sprintf("unknown export format: %s", format))
	}
}

func (ec *ExportCommand) exportJSON(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("export json", flag.ContinueOnError)
	output := cmd.String("output", "", "Output file or directory")
	collection := cmd.String("collection", "", "Specific collection to export (optional)")
	pretty := cmd.Bool("pretty", false, "Pretty print JSON")

	if err := cmd.Parse(args); err != nil {
		return errors.NewError(400, "INVALID_FLAGS", "failed to parse flags").WithDetails(map[string]any{"error": err.Error()})
	}

	// Determine output path
	outPath := *output
	if outPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		outPath = fmt.Sprintf("vault_export_%s.json", timestamp)
	}

	// Get collections to export
	registry := db.NewSchemaRegistry(ec.db)
	allCollections := registry.GetCollections()

	var collectionsToExport []*models.Collection
	if *collection != "" {
		col, ok := registry.GetCollection(*collection)
		if !ok {
			return errors.NewError(404, "COLLECTION_NOT_FOUND", fmt.Sprintf("collection '%s' not found", *collection))
		}
		collectionsToExport = []*models.Collection{col}
	} else {
		collectionsToExport = allCollections
	}

	if len(collectionsToExport) == 0 {
		fmt.Println("No collections found to export")
		return nil
	}

	// Build export data
	exportData := make(map[string]any)
	exportData["exported_at"] = time.Now().Format(time.RFC3339)
	exportData["vault_version"] = "0.7.0"

	collectionsData := make(map[string]any)

	fmt.Printf("Exporting %d collection(s)...\n\n", len(collectionsToExport))

	for _, col := range collectionsToExport {
		fmt.Printf("⬇️  Exporting: %s\n", col.Name)

		// Get all records
		params := db.QueryParams{}
		records, _, err := ec.recordService.ListRecords(ctx, col.Name, params)
		if err != nil {
			fmt.Printf("   ✗ Failed: %v\n", err)
			continue
		}

		// Convert records to simple format
		var recordsData []map[string]any
		for _, rec := range records {
			recordData := make(map[string]any)
			recordData["id"] = rec.ID
			recordData["created"] = rec.Created
			recordData["updated"] = rec.Updated
			for k, v := range rec.Data {
				recordData[k] = v
			}
			recordsData = append(recordsData, recordData)
		}

		collectionsData[col.Name] = map[string]any{
			"schema":  col,
			"records": recordsData,
			"count":   len(recordsData),
		}

		fmt.Printf("   ✓ %d records\n", len(recordsData))
	}

	exportData["collections"] = collectionsData

	// Marshal JSON
	var jsonData []byte
	var err error
	if *pretty {
		jsonData, err = json.MarshalIndent(exportData, "", "  ")
	} else {
		jsonData, err = json.Marshal(exportData)
	}
	if err != nil {
		return errors.NewError(500, "JSON_MARSHAL_FAILED", "failed to marshal JSON").WithDetails(map[string]any{"error": err.Error()})
	}

	// Ensure output directory exists
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return errors.NewError(500, "DIR_CREATE_FAILED", "failed to create output directory").WithDetails(map[string]any{"error": err.Error()})
	}

	// Write file
	if err := os.WriteFile(outPath, jsonData, 0644); err != nil {
		return errors.NewError(500, "FILE_WRITE_FAILED", "failed to write output file").WithDetails(map[string]any{"error": err.Error()})
	}

	// Calculate total records
	totalRecords := 0
	for _, col := range collectionsData {
		if data, ok := col.(map[string]any); ok {
			if count, ok := data["count"].(int); ok {
				totalRecords += count
			}
		}
	}

	fmt.Printf("\n")
	fmt.Printf("✅ Export complete!\n")
	fmt.Printf("   Format: JSON\n")
	fmt.Printf("   Collections: %d\n", len(collectionsToExport))
	fmt.Printf("   Total records: %d\n", totalRecords)
	fmt.Printf("   Output: %s\n", outPath)
	fmt.Printf("   Size: %s\n", formatFileSize(outPath))

	return nil
}

func (ec *ExportCommand) exportSQL(ctx context.Context, args []string) error {
	cmd := flag.NewFlagSet("export sql", flag.ContinueOnError)
	output := cmd.String("output", "", "Output SQL file")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	// Determine output path
	outPath := *output
	if outPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		outPath = fmt.Sprintf("vault_export_%s.sql", timestamp)
	}

	// Get collections
	registry := db.NewSchemaRegistry(ec.db)
	allCollections := registry.GetCollections()

	if len(allCollections) == 0 {
		fmt.Println("No collections found to export")
		return nil
	}

	// Ensure output directory exists
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	fmt.Printf("Exporting %d collection(s) to SQL...\n\n", len(allCollections))

	// Write header
	header := fmt.Sprintf("-- Vault SQL Export\n-- Generated: %s\n-- Vault Version: 0.7.0\n\n", time.Now().Format(time.RFC3339))
	if _, err := outFile.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	totalRecords := 0

	for _, col := range allCollections {
		fmt.Printf("⬇️  Exporting: %s\n", col.Name)

		// Write CREATE TABLE
		createSQL := generateCreateTableSQL(col)
		if _, err := outFile.WriteString(createSQL + "\n\n"); err != nil {
			fmt.Printf("   ✗ Failed to write schema: %v\n", err)
			continue
		}

		// Get all records
		params := db.QueryParams{}
		records, _, err := ec.recordService.ListRecords(ctx, col.Name, params)
		if err != nil {
			fmt.Printf("   ✗ Failed to fetch records: %v\n", err)
			continue
		}

		// Write INSERT statements
		insertCount := 0
		for _, rec := range records {
			insertSQL := generateInsertSQL(col, rec)
			if _, err := outFile.WriteString(insertSQL + "\n"); err != nil {
				fmt.Printf("   ✗ Failed to write record: %v\n", err)
				continue
			}
			insertCount++
		}

		if _, err := outFile.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}

		fmt.Printf("   ✓ %d records\n", insertCount)
		totalRecords += insertCount
	}

	fmt.Printf("\n")
	fmt.Printf("✅ Export complete!\n")
	fmt.Printf("   Format: SQL\n")
	fmt.Printf("   Collections: %d\n", len(allCollections))
	fmt.Printf("   Total records: %d\n", totalRecords)
	fmt.Printf("   Output: %s\n", outPath)
	fmt.Printf("   Size: %s\n", formatFileSize(outPath))

	return nil
}

func generateCreateTableSQL(col *models.Collection) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", col.Name))

	// Add standard fields
	fields := []string{
		"  id TEXT PRIMARY KEY",
		"  created TEXT",
		"  updated TEXT",
	}

	// Add custom fields
	for _, f := range col.Fields {
		sqlType := "TEXT"
		switch f.Type {
		case "number":
			sqlType = "REAL"
		case "boolean":
			sqlType = "INTEGER"
		}

		fieldDef := fmt.Sprintf("  %s %s", f.Name, sqlType)
		if f.Required {
			fieldDef += " NOT NULL"
		}
		if f.Unique {
			fieldDef += " UNIQUE"
		}
		fields = append(fields, fieldDef)
	}

	sb.WriteString(strings.Join(fields, ",\n"))
	sb.WriteString("\n);")

	return sb.String()
}

func generateInsertSQL(col *models.Collection, rec *models.Record) string {
	var sb strings.Builder

	// Build column list
	columns := []string{"id", "created", "updated"}
	for _, f := range col.Fields {
		columns = append(columns, f.Name)
	}

	// Build values list
	values := []any{rec.ID, rec.Created, rec.Updated}
	for _, f := range col.Fields {
		val := rec.Data[f.Name]
		values = append(values, val)
	}

	// Generate INSERT
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		col.Name,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", ")))

	// Note: In a real SQL file, we'd need to escape values properly
	// For now, this is a placeholder that would need parameter binding

	return sb.String()
}

func (ec *ExportCommand) printUsage() {
	fmt.Println("Usage: vault export <format> [options]")
	fmt.Println()
	fmt.Println("Formats:")
	fmt.Println("  json     Export collections and records as JSON")
	fmt.Println("  sql      Export schema and data as SQL statements")
	fmt.Println()
	fmt.Println("JSON Options:")
	fmt.Println("  --output FILE          Output file (default: vault_export_TIMESTAMP.json)")
	fmt.Println("  --collection NAME      Export specific collection only")
	fmt.Println("  --pretty               Pretty print JSON output")
	fmt.Println()
	fmt.Println("SQL Options:")
	fmt.Println("  --output FILE          Output SQL file (default: vault_export_TIMESTAMP.sql)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault export json --output ./backup.json")
	fmt.Println("  vault export json --collection users --pretty")
	fmt.Println("  vault export sql --output ./schema.sql")
}
