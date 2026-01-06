package server

import (
	"context"
	"net"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
)

func (gs *grpcServer) StartGracefully(ctx context.Context, log logger.Logger, lis net.Listener) {
	go func() {
		if err := gs.srv.Serve(lis); err != nil {
			log.Warn("grpc-server has been stopped", logger.NewField("error", err))
		}
	}()

	<-ctx.Done()
	log.Info("stopping server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	defer close(done)

	go func() {
		gs.srv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("server stopped gracefully")
	case <-shutdownCtx.Done():
		gs.srv.Stop()
		log.Warn("server stopped forcibly")
	}

	lis.Close()
	log.Info("tcp connection closed")
}
