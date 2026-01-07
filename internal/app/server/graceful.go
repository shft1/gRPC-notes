package server

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc"
)

func (gs *grpcServer) StartGracefully(ctx context.Context, log logger.Logger, lis net.Listener) {
	go func() {
		if err := gs.srv.Serve(lis); errors.Is(err, grpc.ErrServerStopped) {
			log.Warn("request to a closed grpc-server")
		} else {
			log.Info("grpc-server has been stopped")
		}
	}()
	log.Info("grpc-server is running")
	<-ctx.Done()
	log.Info("stopping server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})

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
