package core

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	Port     int    `json:"port"`
	DBPath   string `json:"db_path"`
	LogLevel string `json:"log_level"`
	LogFormat string `json:"log_format"`
}

func LoadConfig() *Config {
	// Defaults
	cfg := &Config{
		Port:     8090,
		DBPath:   "./vault_data/vault.db",
		LogLevel: "INFO",
		LogFormat: "text",
	}

	// Load from config.json if exists
	if data, err := os.ReadFile("config.json"); err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	// Environment overrides
	if portStr := os.Getenv("VAULT_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}
	if dbPath := os.Getenv("VAULT_DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}
	if logLevel := os.Getenv("VAULT_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if logFormat := os.Getenv("VAULT_LOG_FORMAT"); logFormat != "" {
		cfg.LogFormat = logFormat
	}

	return cfg
}
