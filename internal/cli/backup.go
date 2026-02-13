package cli

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupCommand struct {
	config *Config
}

type Config struct {
	DataDir string
	DBPath  string
}

func NewBackupCommand(dataDir, dbPath string) *BackupCommand {
	return &BackupCommand{
		config: &Config{
			DataDir: dataDir,
			DBPath:  dbPath,
		},
	}
}

func (bc *BackupCommand) Run(args []string) error {
	if len(args) < 1 {
		bc.printUsage()
		return fmt.Errorf("no backup subcommand provided")
	}

	subcommand := args[0]

	switch subcommand {
	case "create":
		return bc.Create(args[1:])
	case "list":
		return bc.List(args[1:])
	case "restore":
		return bc.Restore(args[1:])
	default:
		bc.printUsage()
		return fmt.Errorf("unknown backup subcommand: %s", subcommand)
	}
}

func (bc *BackupCommand) Create(args []string) error {
	cmd := flag.NewFlagSet("backup create", flag.ContinueOnError)
	output := cmd.String("output", "", "Output backup file path (default: vault_backup_TIMESTAMP.zip)")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	// Generate default filename if not provided
	if *output == "" {
		*output = fmt.Sprintf("vault_backup_%s.zip", time.Now().Format("20060102_150405"))
	}

	// Create backup file
	backupFile, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backupFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(backupFile)
	defer zipWriter.Close()

	// Add database file
	if err := bc.addFileToZip(zipWriter, bc.config.DBPath, "vault.db"); err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}

	// Add config.json if exists
	if _, err := os.Stat("config.json"); err == nil {
		if err := bc.addFileToZip(zipWriter, "config.json", "config.json"); err != nil {
			return fmt.Errorf("failed to backup config: %w", err)
		}
	}

	// Add storage directory if exists
	storageDir := filepath.Join(bc.config.DataDir, "storage")
	if _, err := os.Stat(storageDir); err == nil {
		if err := bc.addDirToZip(zipWriter, storageDir, "storage"); err != nil {
			return fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	fmt.Printf("✓ Backup created successfully\n")
	fmt.Printf("  File: %s\n", *output)
	fmt.Printf("  Size: %s\n", formatFileSize(*output))

	return nil
}

func (bc *BackupCommand) List(args []string) error {
	cmd := flag.NewFlagSet("backup list", flag.ContinueOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	// List backup files in current directory
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	var backups []os.FileInfo
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".zip" && filepath.Base(entry.Name())[:6] == "vault_" {
			info, err := entry.Info()
			if err == nil {
				backups = append(backups, info)
			}
		}
	}

	if len(backups) == 0 {
		fmt.Println("No backup files found")
		return nil
	}

	fmt.Printf("Total backups: %d\n\n", len(backups))
	fmt.Printf("%-40s %-15s %-20s\n", "Filename", "Size", "Modified")
	fmt.Println(fmt.Sprintf("%s", "-----------------------------------------------------------"))

	for _, backup := range backups {
		fmt.Printf("%-40s %-15s %-20s\n",
			backup.Name(),
			formatSize(backup.Size()),
			backup.ModTime().Format("2006-01-02 15:04:05"),
		)
	}

	return nil
}

func (bc *BackupCommand) Restore(args []string) error {
	cmd := flag.NewFlagSet("backup restore", flag.ContinueOnError)
	input := cmd.String("input", "", "Backup file to restore")
	force := cmd.Bool("force", false, "Skip confirmation prompt")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *input == "" {
		return fmt.Errorf("input backup file is required")
	}

	// Check if backup file exists
	if _, err := os.Stat(*input); err != nil {
		return fmt.Errorf("backup file not found: %s", *input)
	}

	// Confirm restoration
	if !*force {
		fmt.Printf("This will overwrite your current database and storage. Continue? (yes/no): ")
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			fmt.Println("Restore cancelled")
			return nil
		}
	}

	// Open backup file
	reader, err := zip.OpenReader(*input)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer reader.Close()

	// Extract files
	for _, file := range reader.File {
		// Determine target path
		var targetPath string
		switch file.Name {
		case "vault.db":
			targetPath = bc.config.DBPath
		case "config.json":
			targetPath = "config.json"
		default:
			// Storage files
			targetPath = filepath.Join(bc.config.DataDir, file.Name)
		}

		// Create directories if needed
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Extract file
		if !file.FileInfo().IsDir() {
			srcFile, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in backup: %w", err)
			}

			dstFile, err := os.Create(targetPath)
			if err != nil {
				srcFile.Close()
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				srcFile.Close()
				dstFile.Close()
				return fmt.Errorf("failed to extract file: %w", err)
			}

			srcFile.Close()
			dstFile.Close()
		}
	}

	fmt.Printf("✓ Backup restored successfully\n")
	fmt.Printf("  From: %s\n", *input)

	return nil
}

func (bc *BackupCommand) addFileToZip(zw *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipPath

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func (bc *BackupCommand) addDirToZip(zw *zip.Writer, dirPath, zipPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		zipFilePath := filepath.Join(zipPath, rel)

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = zipFilePath

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}

func (bc *BackupCommand) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  vault backup create [--output FILE]")
	fmt.Println("  vault backup list")
	fmt.Println("  vault backup restore --input FILE [--force]")
}

func formatSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

func formatFileSize(filePath string) string {
	info, err := os.Stat(filePath)
	if err != nil {
		return "unknown"
	}
	return formatSize(info.Size())
}
