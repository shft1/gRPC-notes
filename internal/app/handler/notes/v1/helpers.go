package v1

import (
	"context"
	"errors"

	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *NoteHandler) isContextError(err error) bool {
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

func (h *NoteHandler) errorSendHandling(err error) error {
	if h.isContextError(err) {
		h.log.Info("context done", logger.NewField("reason", err))
		return nil
	}
	h.log.Error("failed to send message", logger.NewField("error", err))
	return status.Error(codes.Internal, "internal error")
}

func (h *NoteHandler) errorRecieveHandling(err error) error {
	if h.isContextError(err) {
		h.log.Info("context done", logger.NewField("reason", err))
		return nil
	}
	h.log.Error("failed to recieve message", logger.NewField("error", err))
	return status.Error(codes.Internal, "internal error")
}

func (h *NoteHandler) errorConsumeHandling(err error) error {
	if h.isContextError(err) {
		h.log.Info("context done", logger.NewField("reason", err))
		return nil
	}
	h.log.Error("failed to consume event bus", logger.NewField("error", err))
	return status.Error(codes.Internal, "internal error")
}
