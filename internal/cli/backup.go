package cli

import (
	"archive/tar"
	"compress/gzip"
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
	output := cmd.String("output", "", "Output backup file path (default: vault_backup_TIMESTAMP.tar.gz)")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	// Generate default filename if not provided
	if *output == "" {
		*output = fmt.Sprintf("vault_backup_%s.tar.gz", time.Now().Format("20060102_150405"))
	}

	// Create backup file
	backupFile, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backupFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(backupFile)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Add database file
	if err := bc.addFileToTar(tarWriter, bc.config.DBPath, "vault.db"); err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}

	// Add config.json if exists
	if _, err := os.Stat("config.json"); err == nil {
		if err := bc.addFileToTar(tarWriter, "config.json", "config.json"); err != nil {
			return fmt.Errorf("failed to backup config: %w", err)
		}
	}

	// Add storage directory if exists
	storageDir := filepath.Join(bc.config.DataDir, "storage")
	if _, err := os.Stat(storageDir); err == nil {
		if err := bc.addDirToTar(tarWriter, storageDir, "storage"); err != nil {
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
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
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
	backupFile, err := os.Open(*input)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer backupFile.Close()

	// Create gzip reader
	gzipReader, err := gzip.NewReader(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		// Determine target path
		var targetPath string
		switch header.Name {
		case "vault.db":
			targetPath = bc.config.DBPath
		case "config.json":
			targetPath = "config.json"
		default:
			// Storage files
			targetPath = filepath.Join(bc.config.DataDir, header.Name)
		}

		// Create directories if needed
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Extract file
		if header.Typeflag == tar.TypeReg {
			file, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return fmt.Errorf("failed to extract file: %w", err)
			}
			file.Close()
		}
	}

	fmt.Printf("✓ Backup restored successfully\n")
	fmt.Printf("  From: %s\n", *input)

	return nil
}

func (bc *BackupCommand) addFileToTar(tw *tar.Writer, filePath, tarPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: tarPath,
		Size: info.Size(),
		Mode: 0644,
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}

func (bc *BackupCommand) addDirToTar(tw *tar.Writer, dirPath, tarPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		tarFilePath := filepath.Join(tarPath, rel)

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header := &tar.Header{
			Name: tarFilePath,
			Size: info.Size(),
			Mode: 0644,
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		_, err = io.Copy(tw, file)
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
