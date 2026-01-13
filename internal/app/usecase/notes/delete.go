package notes

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (uc *noteUsecase) DeleteByID(ctx context.Context, id uuid.UUID) (*notes.Note, error) {
	note, err := uc.repo.DeleteByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}
