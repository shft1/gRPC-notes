package notes

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (uc *noteUsecase) GetByID(ctx context.Context, id uuid.UUID) (*notes.Note, error) {
	note, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (uc *noteUsecase) GetMulti(ctx context.Context) ([]*notes.Note, error) {
	notes, err := uc.repo.GetMulti(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
