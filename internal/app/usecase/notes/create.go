package notes

import (
	"context"

	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (uc *noteUsecase) Create(ctx context.Context, in *notes.NoteCreate) (*notes.Note, error) {
	note, err := uc.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return note, nil
}
