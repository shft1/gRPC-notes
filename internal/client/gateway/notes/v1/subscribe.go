package v1

import (
	"context"
	"io"
	"reflect"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc/status"
)

func (gw *noteGateway) SubscribeToEvents(ctx context.Context, errChan chan<- error) {
	eventCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	st, err := gw.client.SubscribeToEvents(eventCtx, &pb.Empty{})

	if err != nil {
		gw.log.Error("failed to init server-stream", logger.NewField("error", err))
		if stat, ok := status.FromError(err); ok {
			err = mapErrorRPC(gw.log, stat)
		} else {
			err = mapError(gw.log, err)
		}
		errChan <- err
		return
	}

	errChan <- nil
	gw.log.Info("successfully create server-stream event")

	for {
		if eventCtx.Err() != nil {
			gw.log.Info("context canceled", logger.NewField("reason", eventCtx.Err()))
			break
		}
		event, err := st.Recv()
		if err != nil {
			if err == io.EOF {
				gw.log.Warn("server closed server-stream")
				break
			}
			gw.log.Error("failed to recieve event", logger.NewField("error", err))
			break
		}
		if reflect.TypeOf(event.Event) != reflect.TypeOf(&pb.EventResponse_Health{}) {
			gw.log.Info("[NOTIFICATION]: new note", logger.NewField("message", event))
		}
	}
	gw.log.Info("subscription stoped")
}
