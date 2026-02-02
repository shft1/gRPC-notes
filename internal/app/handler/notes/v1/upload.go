package v1

import (
	"io"

	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc"
)

func (h *NoteHandler) UploadMetrics(st grpc.ClientStreamingServer[pb.MetricRequest, pb.SummaryResponse]) error {
	var metricsTotal int64

	for {
		metric, err := st.Recv()
		if err != nil {
			if err == io.EOF {
				h.log.Info("client stop streaming")
				return st.SendAndClose(&pb.SummaryResponse{Summary: metricsTotal})
			}
			return errorRecieveHandling(h.log, err)
		}
		metricsTotal += metric.Number
	}
}
