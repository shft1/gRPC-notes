package v1

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

func (h *NoteHandler) DeleteByID(ctx context.Context, req *pb.NoteIDRequest) (*pb.Note, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, mapError(h.log, notes.ErrInvalidUUID)
	}
	id, _ := uuid.Parse(req.Id)
	note, err := h.noteUsecase.DeleteByID(ctx, id)
	if err != nil {
		return nil, mapError(h.log, err)
	}
	notePB := toDTOResponse(note)
	if err := protovalidate.Validate(notePB); err != nil {
		return nil, mapError(h.log, notes.ErrNoteResponse)
	}
	return notePB, nil
}
