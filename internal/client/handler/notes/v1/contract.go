package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

type gateway interface {
	Create(ctx context.Context, item *notes.NoteCreate) (*notes.Note, error)
	GetByID(ctx context.Context, id uuid.UUID) (*notes.Note, error)
	GetMulti(ctx context.Context) ([]*notes.Note, error)
	DeleteByID(ctx context.Context, id uuid.UUID) (*notes.Note, error)
	SubscribeToEvents(ctx context.Context, errChan chan<- error)
	UploadMetrics(ctx context.Context) (int64, error)
}
