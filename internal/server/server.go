package server

import (
	"context"
	"net/http"

	"github.com/Feokrat/user-balance-api/internal/config"
)

type HTTPserver struct {
	httpServer *http.Server
}

func NewHTTPserver(cfg *config.Config, handler http.Handler) *HTTPserver {
	return &HTTPserver{
		httpServer: &http.Server{
			Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler: handler,
		},
	}
}

func (s *HTTPserver) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPserver) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
