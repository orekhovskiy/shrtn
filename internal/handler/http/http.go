package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/orekhovskiy/shrtn/config"
)

type Server struct {
	srv *http.Server
}

func NewServer(opts *config.Config) *Server {
	server := &Server{
		&http.Server{
			Addr: opts.ServerAddress,
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
