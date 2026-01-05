package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *NoteHandler) GetByID(ctx context.Context, req *pb.NoteIDRequest) (*pb.Note, error) {
	id, err := uuid.Parse(req.GetUuid())
	if err != nil || id == uuid.Nil {
		return nil, mapError(notes.ErrInvalidUUID)
	}
	note, err := h.noteUsecase.GetByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return toDTOResponse(note), nil
}

func (h *NoteHandler) GetMulti(ctx context.Context, _ *emptypb.Empty) (*pb.NoteList, error) {
	notes, err := h.noteUsecase.GetMulti(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return toDTOListResponse(notes), nil
}
