package middleware

import (
	"context"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc"
)

func NewLoggerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log.Info("start request", logger.NewField("method", info.FullMethod))

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			log.Warn("failed request",
				logger.NewField("method", info.FullMethod),
				logger.NewField("duration", duration.Milliseconds()),
				logger.NewField("error", err))
		} else {
			log.Info("successfull request",
				logger.NewField("method", info.FullMethod),
				logger.NewField("duration", duration.Milliseconds()))
		}
		return resp, err
	}
}
