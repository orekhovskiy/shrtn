package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	host = "localhost"
	port = "8080"
)

type Server struct {
	srv *http.Server
}

func NewServer() *Server {
	server := &Server{
		&http.Server{
			Addr: fmt.Sprintf("%s:%s", host, port),
		},
	}

	return server
}

func (s *Server) RegisterRoutes(r *Router) {
	s.srv.Handler = r.router
}

func (s *Server) Start() error {
	if s.srv.Handler == nil {
		return fmt.Errorf("no routes have registered")
	}

	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(context.Background())
}
