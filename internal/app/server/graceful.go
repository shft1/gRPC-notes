package server

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"
)

func (gs *grpcServer) StartGracefully(ctx context.Context, lis net.Listener) {
	gs.wg.Add(1)
	go func() {
		defer gs.wg.Done()
		if err := gs.srv.Serve(lis); errors.Is(err, grpc.ErrServerStopped) {
			gs.log.Warn("request to a closed grpc-server")
		} else {
			gs.log.Info("grpc-server has been stopped")
		}
	}()

	gs.log.Info("grpc-server started")
	<-ctx.Done()
	gs.log.Info("stopping grpc-server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})

	gs.wg.Add(1)
	go func() {
		defer gs.wg.Done()
		gs.srv.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		gs.log.Info("server stopped gracefully")
	case <-shutdownCtx.Done():
		gs.srv.Stop()
		gs.log.Warn("server stopped forcibly")
	}
	gs.wg.Wait()

	lis.Close()
	gs.log.Info("tcp connection closed")
}
