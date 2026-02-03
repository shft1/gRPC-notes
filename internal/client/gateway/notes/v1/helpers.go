package v1

import (
	"context"
	"errors"

	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isContextError(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Canceled, codes.DeadlineExceeded:
			return true
		}
	}
	return false
}

func errorSendHandling(log logger.Logger, err error) error {
	if isContextError(err) {
		log.Info("context done", logger.NewField("reason", err))
		return nil
	}
	log.Error("failed to send message", logger.NewField("error", err))
	// return &statusrpc.Status{Code: int32(codes.Internal), Message: "internal error"}
	return status.Error(codes.Internal, "internal error")
}

func errorRecieveHandling(log logger.Logger, err error) error {
	if isContextError(err) {
		log.Info("context done", logger.NewField("reason", err))
		return nil
	}
	log.Error("failed to recieve message", logger.NewField("error", err))
	// return &statusrpc.Status{Code: int32(codes.Internal), Message: "internal error"}
	return status.Error(codes.Internal, "internal error")
}
