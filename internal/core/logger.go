package core

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

func InitLogger(level string, format string, logPath string) error {
	var logLevel slog.Level
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var output io.Writer = os.Stdout
	if logPath != "" {
		fl, err := NewFileLogger(logPath)
		if err != nil {
			return err
		}
		SetGlobalFileLogger(fl)
		output = fl.file
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: logLevel}

	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}
