package core

import (
	"context"
)

type contextKey string

const RequestIDKey contextKey = "request_id"
const AuthRecordKey contextKey = "auth_record"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

func WithAuth(ctx context.Context, record any) context.Context {
	return context.WithValue(ctx, AuthRecordKey, record)
}

func GetAuth(ctx context.Context) any {
	return ctx.Value(AuthRecordKey)
}
