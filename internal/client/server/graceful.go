package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
)

func (s *httpServer) StartGracefully(ctx context.Context) {
	s.server.RegisterOnShutdown(func() {
		s.log.Warn("Shutting down client-service...")
	})

	go func() {
		if err := s.server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			s.log.Info("http-server correctly stoppted")
		} else {
			s.log.Error("http-server failed to start or crashed", logger.NewField("error", err))
		}
	}()

	s.log.Info("http-server started")
	<-ctx.Done()
	s.log.Info("stopping http-server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.log.Error("failed to shutdown http-server", logger.NewField("error", err))
		return
	}
	s.log.Info("client-service successfully stoppted")
}
