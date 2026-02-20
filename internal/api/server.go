package api

import (
	"context"
	"net/http"

	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/config"
	handler "github.com/MyNameIsWhaaat/shortener/internal/httpapi"
)

type Server struct {
	httpServer *http.Server
	router     *Router
	config     *config.Config
}

func NewServer(cfg *config.Config, h *handler.Handler) *Server {
	router := NewRouter(h)

	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router.GetHandler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		router:     router,
		config:     cfg,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
