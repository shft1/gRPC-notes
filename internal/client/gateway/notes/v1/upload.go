package v1

import (
	"context"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

func (gw *noteGateway) UploadMetrics(ctx context.Context) (int64, error) {
	st, err := gw.client.UploadMetrics(ctx)

	if err != nil {
		gw.log.Error("failed to init client-stream", logger.NewField("error", err))
		return 0, notes.ErrServiceInternal
	}

	gw.log.Info("successfully create client-stream metrics")

	for i := 1; i < 100; i++ {
		err := st.Send(&pb.MetricRequest{Number: int64(i)})

		if err != nil {
			if ctx.Err() != nil {
				gw.log.Info("context canceled", logger.NewField("reason", ctx.Err()))
				return 0, notes.ErrServiceInternal
			}
			gw.log.Error("failed to send metric", logger.NewField("error", err))
			return 0, notes.ErrServiceInternal
		}
	}
	metricsTotal, err := st.CloseAndRecv()

	if err != nil {
		gw.log.Error("failed to recieve total metrics", logger.NewField("error", err))
		return 0, notes.ErrServiceInternal
	}
	return metricsTotal.Summary, nil
}
