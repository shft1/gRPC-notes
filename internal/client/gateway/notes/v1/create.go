package v1

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/shared"
	"google.golang.org/grpc/status"
)

func (gw *noteGateway) Create(ctx context.Context, item *notes.NoteCreate) (*notes.Note, error) {
	pbCreate := toDTORequest(item)

	if err := protovalidate.Validate(pbCreate); err != nil {
		return nil, mapError(gw.log, notes.ErrInvalidData)
	}
	token, ok := ctx.Value(shared.AuthKey).(string)
	if ok {
		md := gw.createMetadata(token)
		ctx = gw.putMetadataToContext(ctx, md)
	}

	out, err := gw.client.Create(ctx, toDTORequest(item))
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, mapErrorRPC(gw.log, st)
		} else {
			return nil, mapError(gw.log, err)
		}
	}
	if err := protovalidate.Validate(out); err != nil {
		return nil, mapError(gw.log, notes.ErrNoteResponse)
	}
	return toDomainResponse(out), nil
}
