package api

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/rules"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", id)

		ctx := core.WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(ww, r)

		requestID := ww.Header().Get("X-Request-ID")

		slog.Info("Request",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.status,
			"duration", time.Since(start),
		)
	})
}

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

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
				tokenStr = tokenStr[7:]
			} else {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := auth.ValidateToken(r.Context(), tokenStr, secret)
			if err != nil {
				errors.SendError(w, err)
				return
			}

			ctx := core.WithAuth(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authClaims := core.GetAuth(r.Context())
		claims, ok := authClaims.(*auth.Claims)
		if !ok || claims == nil {
			errors.SendError(w, errors.NewError(http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required"))
			return
		}

		// Simple rule: any user in 'users' collection is admin for this prototype
		if claims.Collection != "users" {
			errors.SendError(w, errors.NewError(http.StatusForbidden, "FORBIDDEN", "Admin access required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetEvaluationContext(r *http.Request, recordData map[string]any) rules.EvaluationContext {
	authClaims := core.GetAuth(r.Context())

	evalCtx := rules.EvaluationContext{
		Auth:    make(map[string]any),
		Data:    make(map[string]any),
		Record:  recordData,
		IsAdmin: false,
	}

	// If we have JWT claims, populate @request.auth
	if claims, ok := authClaims.(*auth.Claims); ok {
		evalCtx.Auth["id"] = claims.RecordID
		evalCtx.Auth["collection"] = claims.Collection

		// For now, any user in the 'users' collection is an admin if they have a certain flag
		// but let's just say any verified user in 'users' is admin for this prototype
		// or better: check a specific 'is_admin' field.
		if claims.Collection == "users" {
			evalCtx.IsAdmin = true // Placeholder: Real logic would check user record
		}
	}

	return evalCtx
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
