package v1

import (
	"context"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"google.golang.org/grpc/status"
)

func (gw *noteGateway) Create(ctx context.Context, item *notes.NoteCreate) (*notes.Note, error) {
	out, err := gw.client.Create(ctx, toDTORequest(item))
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
