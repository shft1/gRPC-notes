package notes

import (
	"context"
	"unicode/utf8"

	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (uc *useCase) Create(ctx context.Context, in *notes.NoteCreate) (*notes.Note, error) {
	if err := uc.validateIn(in); err != nil {
		return nil, err
	}
	note, err := uc.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (uc *useCase) validateIn(in *notes.NoteCreate) error {
	lenTitle := utf8.RuneCountInString(in.Title)
	lenDesc := utf8.RuneCountInString(in.Desc)

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
