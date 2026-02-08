package api

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Base routes
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	return mux
}
