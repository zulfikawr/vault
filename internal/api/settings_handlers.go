package api

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/core"
)

type SettingsHandler struct {
	config *core.Config
}

func NewSettingsHandler(config *core.Config) *SettingsHandler {
	return &SettingsHandler{config: config}
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]interface{}{
		"port":                   h.config.Port,
		"log_level":              h.config.LogLevel,
		"log_format":             h.config.LogFormat,
		"jwt_expiry":             h.config.JWTExpiry,
		"max_file_upload_size":   h.config.MaxFileUploadSize,
		"cors_origins":           h.config.CORSOrigins,
		"rate_limit_per_min":     h.config.RateLimitPerMin,
		"tls_enabled":            h.config.TLSEnabled,
	}

	SendJSON(w, http.StatusOK, settings, nil)
}
