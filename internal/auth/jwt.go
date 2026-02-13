package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zulfikawr/vault/internal/core"
	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
)

type Claims struct {
	RecordID   string `json:"record_id"`
	Collection string `json:"collection"`
	RequestID  string `json:"request_id"`
	jwt.RegisteredClaims
}

func GenerateToken(ctx context.Context, record *models.Record, secret string, expiryHours int) (string, error) {
	requestID := core.GetRequestID(ctx)

	claims := Claims{
		RecordID:   record.ID,
		Collection: record.Collection,
		RequestID:  requestID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   record.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.NewError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to sign token").WithDetails(map[string]any{"error": err.Error()})
	}

	return tokenString, nil
}

func ValidateToken(ctx context.Context, tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.NewError(http.StatusUnauthorized, "INVALID_TOKEN", "Token validation failed").WithDetails(map[string]any{"error": err.Error()})
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewError(http.StatusUnauthorized, "INVALID_TOKEN_CLAIMS", "Invalid token claims")
}
