package server

import (
	"net"
	"net/http"

	"github.com/shft1/grpc-notes/observability/logger"
)

type httpServer struct {
	log    logger.Logger
	server *http.Server
}

func NewHTTPServer(log logger.Logger, router http.Handler, host, port string) *httpServer {
	srv := &http.Server{Addr: net.JoinHostPort(host, port), Handler: router}
	return &httpServer{
		log:    log,
		server: srv,
	}
}
