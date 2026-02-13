package core

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileLogger struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
	path   string
}

var fileLogger *FileLogger

func InitFileLogger(dataDir string) error {
	logPath := filepath.Join(dataDir, "vault.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	fileLogger = &FileLogger{
		file:   file,
		writer: bufio.NewWriter(file),
		path:   logPath,
	}

	return nil
}

func GetFileLogger() *FileLogger {
	return fileLogger
}

func (fl *FileLogger) Write(record slog.Record) error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	line := fmt.Sprintf("[%s] %s: %s\n", record.Time.Format(time.RFC3339), record.Level, record.Message)
	_, err := fl.writer.WriteString(line)
	if err != nil {
		return err
	}
	return fl.writer.Flush()
}

func (fl *FileLogger) Close() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	if fl.file != nil {
		fl.writer.Flush()
		return fl.file.Close()
	}
	return nil
}

func (fl *FileLogger) ReadLogs(limit int) ([]string, error) {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	file, err := os.Open(fl.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) > limit {
		lines = lines[len(lines)-limit:]
	}

	return lines, scanner.Err()
}

func (fl *FileLogger) Clear() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fl.writer.Flush()
	fl.file.Close()

	// Remove the file
	os.Remove(fl.path)

	// Recreate the file
	file, err := os.OpenFile(fl.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	fl.file = file
	fl.writer = bufio.NewWriter(file)
	return nil
}
