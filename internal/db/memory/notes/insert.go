package notes

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain"
)

func (nt *NoteTable) Insert(_ context.Context, n *domain.NoteCreate) (*domain.Note, error) {
	uuid := uuid.New()
}
