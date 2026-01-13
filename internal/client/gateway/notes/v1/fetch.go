package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (gw *noteGateway) GetByID(ctx context.Context, id uuid.UUID) (*notes.Note, error) {
	out, err := gw.client.GetByID(ctx, &pb.NoteIDRequest{Id: id.String()})
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

func (gw *noteGateway) GetMulti(ctx context.Context) ([]*notes.Note, error) {
	out, err := gw.client.GetMulti(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, mapErrorRPC(gw.log, st)
		} else {
			return nil, mapError(gw.log, err)
		}
	}
	notes, err := toDomainListResponse(out.Notes)
	if err != nil {
		return nil, mapError(gw.log, err)
	}
	return notes, nil
}
