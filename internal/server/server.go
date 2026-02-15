package server

import (
	"context"
	"fmt"
	"net"
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
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Vault server started successfully\n")
	fmt.Printf("  Web UI:  http://localhost:%d/\n", s.port)
	fmt.Printf("  API:     http://localhost:%d/api\n\n", s.port)

	if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
