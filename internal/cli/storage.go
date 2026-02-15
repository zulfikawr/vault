package cli

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/db"
)

type StorageCommand struct {
	config *core.Config
	db     *sql.DB
}

type FileEntry struct {
	Name     string
	Path     string
	Size     int64
	IsDir    bool
	Modified time.Time
}

func NewStorageCommand(config *core.Config) *StorageCommand {
	return &StorageCommand{config: config}
}

func (sc *StorageCommand) Run(args []string) error {
	if len(args) < 1 {
		sc.printUsage()
		return fmt.Errorf("no storage subcommand provided")
	}

	subcommand := args[0]

	// Check for help flags before processing subcommands
	if subcommand == "-h" || subcommand == "--help" {
		sc.printUsage()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to database for all commands
	database, err := db.Connect(ctx, sc.config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	sc.db = database

	switch subcommand {
	case "create":
		return sc.Create(ctx, args[1:])
	case "list":
		return sc.List(ctx, args[1:])
	case "get":
		return sc.Get(ctx, args[1:])
	case "delete":
		return sc.Delete(ctx, args[1:])
	default:
		sc.printUsage()
		return fmt.Errorf("unknown storage subcommand: %s", subcommand)
	}
}

func (sc *StorageCommand) printUsage() {
	fmt.Println("Usage: vault storage <subcommand> [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  create --path PATH --file FILE")
	fmt.Println("  list [--path PATH] [--recursive]")
	fmt.Println("  get --path PATH --output FILE")
	fmt.Println("  delete --path PATH [--recursive] [--force]")
}

func (sc *StorageCommand) Create(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("storage create", flag.ExitOnError)
	path := fs.String("path", "", "Storage path")
	file := fs.String("file", "", "File to upload")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *path == "" || *file == "" || *email == "" || *password == "" {
		fmt.Println("Error: --path, --file, --email, and --password are required")
		fs.Usage()
		return fmt.Errorf("missing required flags")
	}

	if err := sc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Prevent directory traversal
	if strings.Contains(*path, "..") {
		return fmt.Errorf("invalid path: contains '..'")
	}

	// Read the file to upload
	fileData, err := os.ReadFile(*file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", *file)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create destination path
	fullPath := filepath.Join(sc.config.StoragePath(), *path)
	destDir := filepath.Dir(fullPath)

	// Create directories if they don't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file to storage
	if err := os.WriteFile(fullPath, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fileSize := sc.formatSize(int64(len(fileData)))
	fmt.Printf("✓ File uploaded successfully (%s) to %s\n", fileSize, *path)
	return nil
}

func (sc *StorageCommand) List(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("storage list", flag.ExitOnError)
	path := fs.String("path", "", "Storage path to list")
	recursive := fs.Bool("recursive", false, "List recursively")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *email == "" || *password == "" {
		fmt.Println("Error: --email and --password are required")
		fs.Usage()
		return fmt.Errorf("missing required flags")
	}

	if err := sc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Default to root storage
	listPath := *path
	if listPath == "" {
		listPath = "."
	}

	// Prevent directory traversal
	if strings.Contains(listPath, "..") {
		return fmt.Errorf("invalid path: contains '..'")
	}

	fullPath := filepath.Join(sc.config.StoragePath(), listPath)

	// Check if path exists
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path not found: %s", listPath)
		}
		return fmt.Errorf("failed to access path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", listPath)
	}

	var entries []FileEntry
	var totalSize int64

	if *recursive {
		// Recursive listing
		err = filepath.WalkDir(fullPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(fullPath, path)
			if relPath == "." {
				return nil
			}

			entry := FileEntry{
				Name:     d.Name(),
				Path:     filepath.ToSlash(relPath),
				Size:     info.Size(),
				IsDir:    d.IsDir(),
				Modified: info.ModTime(),
			}

			entries = append(entries, entry)
			if !d.IsDir() {
				totalSize += info.Size()
			}

			return nil
		})
	} else {
		// Non-recursive listing
		dirEntries, err := os.ReadDir(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		for _, d := range dirEntries {
			info, err := d.Info()
			if err != nil {
				continue
			}

			entry := FileEntry{
				Name:     d.Name(),
				Path:     filepath.ToSlash(filepath.Join(listPath, d.Name())),
				Size:     info.Size(),
				IsDir:    d.IsDir(),
				Modified: info.ModTime(),
			}

			entries = append(entries, entry)
			if !d.IsDir() {
				totalSize += info.Size()
			}
		}
	}

	if err != nil {
		return fmt.Errorf("failed to list directory: %w", err)
	}

	// Display results
	if len(entries) == 0 {
		fmt.Println("No files or folders found")
		return nil
	}

	// Calculate max width for alignment (with max limits)
	maxNameLen := len("Name")
	maxTypeLen := len("Type")
	maxSizeLen := len("Size")
	maxModifiedLen := len("Modified")
	const maxNameWidth = 40
	const maxModifiedWidth = 12

	for _, entry := range entries {
		nameLen := len(entry.Name)
		if nameLen > maxNameLen && nameLen <= maxNameWidth {
			maxNameLen = nameLen
		}
		typeStr := "file"
		if entry.IsDir {
			typeStr = "folder"
		}
		if len(typeStr) > maxTypeLen {
			maxTypeLen = len(typeStr)
		}
		sizeStr := sc.formatSize(entry.Size)
		if len(sizeStr) > maxSizeLen {
			maxSizeLen = len(sizeStr)
		}
		timeStr := sc.formatTime(entry.Modified)
		if len(timeStr) > maxModifiedLen && len(timeStr) <= maxModifiedWidth {
			maxModifiedLen = len(timeStr)
		}
	}

	// Cap the widths
	if maxNameLen > maxNameWidth {
		maxNameLen = maxNameWidth
	}
	if maxModifiedLen > maxModifiedWidth {
		maxModifiedLen = maxModifiedWidth
	}

	fmt.Println()
	// Print top border
	fmt.Printf("┌%s┬%s┬%s┬%s┐\n", strings.Repeat("─", maxNameLen+2), strings.Repeat("─", maxTypeLen+2), strings.Repeat("─", maxSizeLen+2), strings.Repeat("─", maxModifiedLen+2))
	// Print header
	fmt.Printf("│ %-*s │ %-*s │ %-*s │ %-*s │\n", maxNameLen, "Name", maxTypeLen, "Type", maxSizeLen, "Size", maxModifiedLen, "Modified")
	// Print separator
	fmt.Printf("├%s┼%s┼%s┼%s┤\n", strings.Repeat("─", maxNameLen+2), strings.Repeat("─", maxTypeLen+2), strings.Repeat("─", maxSizeLen+2), strings.Repeat("─", maxModifiedLen+2))

	// Count files and folders
	fileCount := 0
	folderCount := 0

	// Print entries
	for _, entry := range entries {
		typeStr := "file"
		sizeStr := sc.formatSize(entry.Size)

		if entry.IsDir {
			typeStr = "folder"
			sizeStr = "-"
			folderCount++
		} else {
			fileCount++
		}

		timeStr := sc.formatTime(entry.Modified)
		name := sc.truncate(entry.Name, maxNameLen)
		time := sc.truncate(timeStr, maxModifiedLen)
		fmt.Printf("│ %-*s │ %-*s │ %-*s │ %-*s │\n", maxNameLen, name, maxTypeLen, typeStr, maxSizeLen, sizeStr, maxModifiedLen, time)
	}

	// Print bottom border
	fmt.Printf("└%s┴%s┴%s┴%s┘\n", strings.Repeat("─", maxNameLen+2), strings.Repeat("─", maxTypeLen+2), strings.Repeat("─", maxSizeLen+2), strings.Repeat("─", maxModifiedLen+2))

	// Print summary
	fmt.Printf("\nTotal: %d files, %d folders, %s\n\n", fileCount, folderCount, sc.formatSize(totalSize))

	return nil
}

func (sc *StorageCommand) Get(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("storage get", flag.ExitOnError)
	path := fs.String("path", "", "File path to get")
	output := fs.String("output", "", "Output file")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")
	force := fs.Bool("force", false, "Overwrite output file if exists")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *path == "" || *output == "" || *email == "" || *password == "" {
		fmt.Println("Error: --path, --output, --email, and --password are required")
		fs.Usage()
		return fmt.Errorf("missing required flags")
	}

	if err := sc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Prevent directory traversal
	if strings.Contains(*path, "..") {
		return fmt.Errorf("invalid path: contains '..'")
	}

	// Check if source file exists
	fullPath := filepath.Join(sc.config.StoragePath(), *path)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", *path)
		}
		return fmt.Errorf("failed to access file: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", *path)
	}

	// Check if output file already exists
	if _, err := os.Stat(*output); err == nil && !*force {
		return fmt.Errorf("output file already exists: %s (use --force to overwrite)", *output)
	}

	// Read source file
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Write to output file
	if err := os.WriteFile(*output, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fileSize := sc.formatSize(int64(len(fileData)))
	fmt.Printf("✓ File downloaded successfully (%s) to %s\n", fileSize, *output)
	return nil
}

func (sc *StorageCommand) Delete(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("storage delete", flag.ExitOnError)
	path := fs.String("path", "", "Path to delete")
	recursive := fs.Bool("recursive", false, "Delete recursively")
	force := fs.Bool("force", false, "Skip confirmation")
	email := fs.String("email", "", "Admin email")
	password := fs.String("password", "", "Admin password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *path == "" || *email == "" || *password == "" {
		fmt.Println("Error: --path, --email, and --password are required")
		fs.Usage()
		return fmt.Errorf("missing required flags")
	}

	if err := sc.authenticateAdmin(ctx, *email, *password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Prevent directory traversal
	if strings.Contains(*path, "..") {
		return fmt.Errorf("invalid path: contains '..'")
	}

	fullPath := filepath.Join(sc.config.StoragePath(), *path)

	// Check if path exists
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path not found: %s", *path)
		}
		return fmt.Errorf("failed to access path: %w", err)
	}

	// Handle directory deletion
	if fileInfo.IsDir() {
		if !*recursive {
			return fmt.Errorf("path is a directory, use --recursive to delete")
		}

		// Calculate size and count before deletion
		var totalSize int64
		var fileCount int
		_ = filepath.WalkDir(fullPath, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, _ := d.Info()
			totalSize += info.Size()
			fileCount++
			return nil
		})

		// Confirmation prompt
		if !*force {
			fmt.Printf("Are you sure you want to delete this directory and all its contents? (%d files, %s) (y/n): ", fileCount, sc.formatSize(totalSize))
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}
			if response != "y" && response != "Y" {
				fmt.Println("Delete cancelled")
				return nil
			}
		}

		// Delete directory
		if err := os.RemoveAll(fullPath); err != nil {
			return fmt.Errorf("failed to delete directory: %w", err)
		}

		fmt.Printf("✓ Directory deleted successfully (%d files, %s freed)\n", fileCount, sc.formatSize(totalSize))
		return nil
	}

	// Handle file deletion
	fileSize := fileInfo.Size()

	// Confirmation prompt
	if !*force {
		fmt.Printf("Are you sure you want to delete this file? (%s) (y/n): ", sc.formatSize(fileSize))
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		if response != "y" && response != "Y" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	fmt.Printf("✓ File deleted successfully (%s freed)\n", sc.formatSize(fileSize))
	return nil
}

// Helper functions

func (sc *StorageCommand) authenticateAdmin(ctx context.Context, email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password are required")
	}

	var storedPassword string
	err := sc.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
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

func (sc *StorageCommand) formatSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%d %s", int64(size), units[unitIndex])
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

func (sc *StorageCommand) formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}

	return t.Format("2006-01-02")
}

func (sc *StorageCommand) truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
