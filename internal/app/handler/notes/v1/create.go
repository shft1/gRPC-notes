package v1

import (
	"context"
	"fmt"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

func (h *NoteHandler) Create(ctx context.Context, req *pb.NoteCreateRequest) (*pb.Note, error) {
	h.normalizeCreateReq(req)

	if err := protovalidate.Validate(req); err != nil {
		return nil, mapError(h.log, fmt.Errorf("%w: %s", notes.ErrInvalidData, err.Error()))
	}
	note, err := h.noteUsecase.Create(ctx, toDomainCreate(req))
	if err != nil {
		return nil, mapError(h.log, err)
	}
	notePB := toDTOResponse(note)
	if err := protovalidate.Validate(notePB); err != nil {
		return nil, mapError(h.log, notes.ErrNoteResponse)
	}
	return notePB, nil
}

func (h *NoteHandler) normalizeCreateReq(req *pb.NoteCreateRequest) {
	req.Title = strings.TrimSpace(req.Title)
	req.Desc = strings.TrimSpace(req.Desc)
}
