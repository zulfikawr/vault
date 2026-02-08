package auth

import (
	"context"
	"net/http"

	"github.com/zulfikawr/vault/internal/core"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", core.NewError(http.StatusInternalServerError, "AUTH_HASH_FAILED", "Failed to hash password").WithDetails(map[string]any{"error": err.Error()})
	}
	return string(bytes), nil
}

func ComparePasswords(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
