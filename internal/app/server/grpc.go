package server

import (
	"net"

	"github.com/shft1/grpc-notes/internal/app/config"
	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc"
)

type grpcServer struct {
	srv *grpc.Server
}

func NewServer(log logger.Logger, cfg *config.AppEnv) (*grpcServer, net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+cfg.PortGRPC)
	if err != nil {
		log.Error("failed to listen tcp", logger.NewField("port", cfg.PortGRPC))
		return nil, nil, err
	}
	return &grpcServer{srv: grpc.NewServer()}, lis, nil

}
