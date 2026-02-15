package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/zulfikawr/vault/internal/errors"
)

func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF check for GET, HEAD, OPTIONS
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Check CSRF token for state-changing operations
		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			token = r.FormValue("csrf_token")
		}

		if token == "" {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "CSRF_TOKEN_MISSING", "CSRF token is required"))
			return
		}

		// In production, validate token against session
		// For now, just check it exists and is non-empty
		if len(token) < 32 {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "CSRF_TOKEN_INVALID", "Invalid CSRF token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
