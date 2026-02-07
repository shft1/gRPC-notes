package stream

import (
	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc"
)

func NewLoggerInterceptor(log logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Debug("request to grpc-stream", logger.NewField("method", info.FullMethod))

		wrappedStream := &wrappedServerStream{
			log:          log,
			ServerStream: ss,
		}

		err := handler(srv, wrappedStream)

		if err != nil {
			log.Error("error with grpc-stream", logger.NewField("method", info.FullMethod))
		}
		log.Debug("grpc-stream closed", logger.NewField("method", info.FullMethod))
		return err
	}
}

type wrappedServerStream struct {
	log logger.Logger
	grpc.ServerStream
}

func (ss *wrappedServerStream) SendMsg(m any) error {
	ss.log.Debug("server send:", logger.NewField("message", m))
	return ss.ServerStream.SendMsg(m)
}

func (ss *wrappedServerStream) RecvMsg(m any) error {
	err := ss.ServerStream.RecvMsg(m)
	ss.log.Debug("server recieve:", logger.NewField("message", m))
	return err
}
