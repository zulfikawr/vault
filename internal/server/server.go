package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	port   int
	server *http.Server
}

func NewServer(port int, handler http.Handler) *Server {
	return &Server{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	fmt.Printf("\n✓ Vault server started successfully\n")
	fmt.Printf("  Web UI:  http://localhost:%d/_/\n", s.port)
	fmt.Printf("  API:     http://localhost:%d/api\n\n", s.port)
	slog.Info("Starting server", "port", s.port)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}
