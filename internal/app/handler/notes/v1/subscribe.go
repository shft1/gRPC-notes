package v1

import (
	"context"
	"errors"

	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *NoteHandler) SubscribeToEvents(req *pb.Empty, st grpc.ServerStreamingServer[pb.EventResponse]) error {
	ctx := st.Context()

	err := st.Send(&pb.EventResponse{
		Event: &pb.EventResponse_Health{Health: &pb.Health{Message: "event subscription is successful"}},
	})
	if err != nil {
		return h.errorSendHandling(err)
	}

	for {
		event, err := h.bus.Consume(ctx)
		if err != nil {
			return h.errorConsumeHandling(err)
		}

		err = st.Send(&pb.EventResponse{
			Event: &pb.EventResponse_Note{
				Note: &pb.NoteEvent{Id: event.ID.String(), Title: event.Title},
			},
		})
		if err != nil {
			return h.errorSendHandling(err)
		}
	}
}

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
		h.log.Info("stream context done", logger.NewField("reason", err))
		return nil
	}
	h.log.Error("failed to send event", logger.NewField("error", err))
	return status.Error(codes.Internal, "internal error")
}

func (h *NoteHandler) errorConsumeHandling(err error) error {
	if h.isContextError(err) {
		h.log.Info("stream context done", logger.NewField("reason", err))
		return nil
	}
	h.log.Error("failed to consume event bus", logger.NewField("error", err))
	return status.Error(codes.Internal, "internal error")
}
