package api

import (
	"net/http"
	"strconv"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
)

type LogsHandler struct{}

func NewLogsHandler() *LogsHandler {
	return &LogsHandler{}
}

func (h *LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	fileLogger := core.GetFileLogger()
	if fileLogger == nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "INTERNAL_ERROR", "Logger not initialized"))
		return
	}

	logs, err := fileLogger.ReadLogs(limit)
	if err != nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "INTERNAL_ERROR", err.Error()))
		return
	}

	SendJSON(w, http.StatusOK, logs, nil)
}

func (h *LogsHandler) ClearLogs(w http.ResponseWriter, r *http.Request) {
	fileLogger := core.GetFileLogger()
	if fileLogger == nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "INTERNAL_ERROR", "Logger not initialized"))
		return
	}

	if err := fileLogger.Clear(); err != nil {
		errors.SendError(w, errors.NewError(http.StatusInternalServerError, "CLEAR_FAILED", err.Error()))
		return
	}

	SendJSON(w, http.StatusOK, map[string]string{"message": "Logs cleared"}, nil)
}
