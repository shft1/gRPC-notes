package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc/status"
)

func (gw *noteGateway) DeleteByID(ctx context.Context, id uuid.UUID) (*notes.Note, error) {
	out, err := gw.client.DeleteByID(ctx, &pb.NoteIDRequest{Uuid: id.String()})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, mapErrorRPC(gw.log, st)
		} else {
			return nil, mapError(gw.log, err)
		}
	}
	note, err := toDomainResponse(out)
	if err != nil {
		return nil, mapError(gw.log, err)
	}
	return note, nil
}
