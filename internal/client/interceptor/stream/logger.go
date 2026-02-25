package stream

import (
	"context"

	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc"
)

func NewLoggerInterceptor(log logger.Logger) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Debug("request gRPC stream", logger.NewField("method", method))

		stream, err := streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			return stream, err
		}
		wrappedStream := &wrappedClientStream{
			log:          log,
			ClientStream: stream,
		}
		return wrappedStream, nil
	}
}

type wrappedClientStream struct {
	log logger.Logger
	grpc.ClientStream
}

func (cs *wrappedClientStream) SendMsg(m any) error {
	cs.log.Debug("client send:", logger.NewField("message", m))
	return cs.ClientStream.SendMsg(m)
}

func (cs *wrappedClientStream) RecvMsg(m any) error {
	err := cs.ClientStream.RecvMsg(m)
	cs.log.Debug("client recieve:", logger.NewField("message", m))
	return err
}
