package core

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	Port      int    `json:"port"`
	DBPath    string `json:"db_path"`
	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"`
	JWTSecret string `json:"jwt_secret"`
	JWTExpiry int    `json:"jwt_expiry"` // in hours
}

func LoadConfig() *Config {
	// Defaults
	cfg := &Config{
		Port:      8090,
		DBPath:    "./vault_data/vault.db",
		LogLevel:  "INFO",
		LogFormat: "text",
		JWTSecret: "change-me-please-use-a-strong-secret",
		JWTExpiry: 72,
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
	if jwtSecret := os.Getenv("VAULT_JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}
	if jwtExpiry := os.Getenv("VAULT_JWT_EXPIRY"); jwtExpiry != "" {
		if expiry, err := strconv.Atoi(jwtExpiry); err == nil {
			cfg.JWTExpiry = expiry
		}
	}

	return cfg
}
