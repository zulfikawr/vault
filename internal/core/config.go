package core

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	Port              int    `json:"port"`
	DBPath            string `json:"db_path"`
	DataDir           string `json:"data_dir"`
	LogLevel          string `json:"log_level"`
	LogFormat         string `json:"log_format"`
	JWTSecret         string `json:"jwt_secret"`
	JWTExpiry         int    `json:"jwt_expiry"` // in hours
	MaxFileUploadSize int64  `json:"max_file_upload_size"` // in bytes
	CORSOrigins       string `json:"cors_origins"` // comma-separated
	RateLimitPerMin   int    `json:"rate_limit_per_min"`
	TLSEnabled        bool   `json:"tls_enabled"`
	TLSCertPath       string `json:"tls_cert_path"`
	TLSKeyPath        string `json:"tls_key_path"`
}

func LoadConfig() *Config {
	// Defaults
	cfg := &Config{
		Port:              8090,
		DBPath:            "./vault_data/vault.db",
		DataDir:           "./vault_data",
		LogLevel:          "INFO",
		LogFormat:         "text",
		JWTSecret:         "change-me-please-use-a-strong-secret",
		JWTExpiry:         72,
		MaxFileUploadSize: 10 * 1024 * 1024, // 10MB
		CORSOrigins:       "*",
		RateLimitPerMin:   60,
		TLSEnabled:        false,
		TLSCertPath:       "",
		TLSKeyPath:        "",
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
	if dataDir := os.Getenv("VAULT_DATA_DIR"); dataDir != "" {
		cfg.DataDir = dataDir
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
	if maxUpload := os.Getenv("VAULT_MAX_FILE_UPLOAD_SIZE"); maxUpload != "" {
		if size, err := strconv.ParseInt(maxUpload, 10, 64); err == nil {
			cfg.MaxFileUploadSize = size
		}
	}
	if corsOrigins := os.Getenv("VAULT_CORS_ORIGINS"); corsOrigins != "" {
		cfg.CORSOrigins = corsOrigins
	}
	if rateLimit := os.Getenv("VAULT_RATE_LIMIT_PER_MIN"); rateLimit != "" {
		if limit, err := strconv.Atoi(rateLimit); err == nil {
			cfg.RateLimitPerMin = limit
		}
	}
	if tlsEnabled := os.Getenv("VAULT_TLS_ENABLED"); tlsEnabled != "" {
		cfg.TLSEnabled = tlsEnabled == "true"
	}
	if tlsCert := os.Getenv("VAULT_TLS_CERT_PATH"); tlsCert != "" {
		cfg.TLSCertPath = tlsCert
	}
	if tlsKey := os.Getenv("VAULT_TLS_KEY_PATH"); tlsKey != "" {
		cfg.TLSKeyPath = tlsKey
	}

	return cfg
}
