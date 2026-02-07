package v1

import (
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc"
)

func (h *NoteHandler) SubscribeToEvents(req *pb.Empty, st grpc.ServerStreamingServer[pb.EventResponse]) error {
	ctx := st.Context()

	err := st.Send(&pb.EventResponse{
		Event: &pb.EventResponse_Health{Health: &pb.Health{Message: "event subscription is successful"}},
	})
	if err != nil {
		return errorSendHandling(h.log, err)
	}

	for {
		event, err := h.bus.Consume(ctx)
		if err != nil {
			return h.errorConsumeHandling(h.log, err)
		}

		err = st.Send(&pb.EventResponse{
			Event: &pb.EventResponse_Note{
				Note: &pb.NoteEvent{Id: event.ID.String(), Title: event.Title},
			},
		})
		if err != nil {
			return errorSendHandling(h.log, err)
		}
	}
}
