package notes

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

type noteRepository interface {
	Create(context.Context, *notes.NoteCreate) (*notes.Note, error)
	GetByID(context.Context, uuid.UUID) (*notes.Note, error)
	GetMulti(context.Context) ([]*notes.Note, error)
	DeleteByID(context.Context, uuid.UUID) (*notes.Note, error)
}
