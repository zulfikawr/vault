package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := core.GetRequestID(r.Context())
				slog.Error("Panic recovered",
					"request_id", requestID,
					"error", err,
					"stack", string(debug.Stack()),
				)

				errors.SendError(w, errors.NewError(
					http.StatusInternalServerError,
					"INTERNAL_SERVER_ERROR",
					"An unexpected error occurred",
				))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
