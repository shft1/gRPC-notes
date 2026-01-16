package server

import (
	"net"
	"sync"

	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc"
)

type grpcServer struct {
	wg  sync.WaitGroup
	log logger.Logger
	srv *grpc.Server
}

func NewServer(log logger.Logger, logInter, authInter grpc.UnaryServerInterceptor, opts ...option) (*grpcServer, net.Listener, error) {
	parameters := setupParameters(opts...)
	lis, err := net.Listen("tcp", net.JoinHostPort("", parameters.port))
	if err != nil {
		log.Error("failed to listen tcp", logger.NewField("port", parameters.port))
		return nil, nil, err
	}
	srv := grpc.NewServer(
		grpc.KeepaliveParams(parameters.ServerParameters),
		grpc.ChainUnaryInterceptor(logInter, authInter),
	)
	log.Debug(
		"server info",
		logger.NewField("port", parameters.port),
		logger.NewField("MaxConnectionIdle", parameters.MaxConnectionIdle),
		logger.NewField("MaxConnectionAge", parameters.MaxConnectionAge),
		logger.NewField("MaxConnectionAgeGrace", parameters.MaxConnectionAgeGrace),
		logger.NewField("Time", parameters.Time),
		logger.NewField("Timeout", parameters.Timeout),
	)
	return &grpcServer{log: log, srv: srv}, lis, nil
}
