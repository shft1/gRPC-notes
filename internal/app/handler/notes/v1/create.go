package v1

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

func (h *NoteHandler) Create(ctx context.Context, req *pb.NoteCreateRequest) (*pb.Note, error) {
	req = h.normalizeCreateReq(req)
	if err := h.validateCreateReq(req); err != nil {
		return nil, mapError(err)
	}
	note, err := h.noteUsecase.Create(ctx, toDomainCreate(req))
	if err != nil {
		return nil, mapError(err)
	}
	return toDTOResponse(note), nil
}

func (h *NoteHandler) normalizeCreateReq(req *pb.NoteCreateRequest) *pb.NoteCreateRequest {
	return &pb.NoteCreateRequest{
		Title: strings.TrimSpace(req.Title),
		Desc:  strings.TrimSpace(req.Desc),
	}
}

func (h *NoteHandler) validateCreateReq(req *pb.NoteCreateRequest) error {
	lenTitle := utf8.RuneCountInString(req.Title)
	lenDesc := utf8.RuneCountInString(req.Desc)
	switch {
	case lenTitle == 0:
		return notes.ErrEmptyTitle
	case lenTitle > 25:
		return notes.ErrTooLongTitle
	case lenDesc > 200:
		return notes.ErrTooLongDesc
	default:
		return nil
	}
}
