package notes

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (r *repository) DeleteByID(_ context.Context, id uuid.UUID) (*notes.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	row, ok := r.noteDB[id]
	switch {
	case !ok:
		return nil, notes.ErrNotFound
	case row.IsDel:
		return nil, notes.ErrIsDeleted
	}
	row.IsDel = true
	row.UpdatedAt = time.Now()
	return rowToDomain(row), nil
}
