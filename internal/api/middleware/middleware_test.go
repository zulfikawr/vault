package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zulfikawr/vault/internal/core"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := core.GetRequestID(r.Context())
		if id == "" {
			t.Error("request id not found in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	respID := w.Header().Get("X-Request-ID")
	if respID == "" {
		t.Error("X-Request-ID header missing from response")
	}
}

func TestChain(t *testing.T) {
	var calls []string
	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "m1")
			next.ServeHTTP(w, r)
		})
	}
	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "m2")
			next.ServeHTTP(w, r)
		})
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "final")
	})

	chained := Chain(finalHandler, m1, m2)
	chained.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	if len(calls) != 3 || calls[0] != "m1" || calls[1] != "m2" || calls[2] != "final" {
		t.Errorf("unexpected chain order: %v", calls)
	}
}
