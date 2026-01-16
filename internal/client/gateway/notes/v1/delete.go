package v1

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"github.com/shft1/grpc-notes/shared"
	"google.golang.org/grpc/status"
)

func (gw *noteGateway) DeleteByID(ctx context.Context, id uuid.UUID) (*notes.Note, error) {
	protoID := &pb.NoteIDRequest{Id: id.String()}

	if err := protovalidate.Validate(protoID); err != nil {
		return nil, mapError(gw.log, notes.ErrInvalidUUID)
	}
	token, ok := ctx.Value(shared.AuthKey).(string)
	if ok {
		md := gw.createMetadata(token)
		ctx = gw.putMetadataToContext(ctx, md)
	}

	out, err := gw.client.DeleteByID(ctx, protoID)
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
