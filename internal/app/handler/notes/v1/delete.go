package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

func (h *NoteHandler) DeleteByID(ctx context.Context, req *pb.NoteIDRequest) (*pb.Note, error) {
	id, err := uuid.Parse(req.GetUuid())
	if err != nil || id == uuid.Nil {
		return nil, mapError(notes.ErrInvalidUUID)
	}
	note, err := h.noteUsecase.DeleteByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return toDTOResponse(note), nil
}
