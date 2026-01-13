package server

import (
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shft1/grpc-notes/internal/client/config"
	"github.com/shft1/grpc-notes/observability/logger"
)

type httpServer struct {
	log    logger.Logger
	server *http.Server
}

func NewHTTPServer(log logger.Logger, router chi.Router, cfg *config.ClientEnv) *httpServer {
	srv := &http.Server{Addr: net.JoinHostPort(cfg.Host, cfg.Port), Handler: router}
	return &httpServer{
		log:    log,
		server: srv,
	}
}
