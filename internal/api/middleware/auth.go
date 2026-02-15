package middleware

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
)

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
