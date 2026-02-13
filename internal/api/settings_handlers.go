package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
)

type SettingsHandler struct {
	config *core.Config
}

func NewSettingsHandler(config *core.Config) *SettingsHandler {
	return &SettingsHandler{config: config}
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, http.StatusOK, h.config, nil)
}

func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.SendError(w, errors.NewError(http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body"))
		return
	}

	// Update config in memory
	if port, ok := updates["port"].(float64); ok {
		h.config.Port = int(port)
	}
	if logLevel, ok := updates["log_level"].(string); ok {
		h.config.LogLevel = logLevel
	}
	if logFormat, ok := updates["log_format"].(string); ok {
		h.config.LogFormat = logFormat
	}
	if jwtExpiry, ok := updates["jwt_expiry"].(float64); ok {
		h.config.JWTExpiry = int(jwtExpiry)
	}
	if maxUpload, ok := updates["max_file_upload_size"].(float64); ok {
		h.config.MaxFileUploadSize = int64(maxUpload)
	}
	if corsOrigins, ok := updates["cors_origins"].(string); ok {
		h.config.CORSOrigins = corsOrigins
	}
	if rateLimit, ok := updates["rate_limit_per_min"].(float64); ok {
		h.config.RateLimitPerMin = int(rateLimit)
	}
	if tlsEnabled, ok := updates["tls_enabled"].(bool); ok {
		h.config.TLSEnabled = tlsEnabled
	}

	// Save to config.json
	configData, err := json.MarshalIndent(h.config, "", "  ")
	if err != nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "CONFIG_MARSHAL_FAILED", "Failed to marshal config"))
		return
	}
	if err := os.WriteFile("config.json", configData, 0644); err != nil {
		errors.Log(r.Context(), err, "write config file")
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "CONFIG_WRITE_FAILED", "Failed to write config file"))
		return
	}

	SendJSON(w, http.StatusOK, map[string]string{"message": "Settings updated"}, nil)
}
