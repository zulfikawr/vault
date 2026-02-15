package service

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/auth"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/rules"
)

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
