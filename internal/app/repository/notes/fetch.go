package notes

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (r *repository) GetByID(_ context.Context, id uuid.UUID) (*notes.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	row, ok := r.noteDB[id]
	switch {
	case !ok:
		return nil, notes.ErrNotFound
	case row.IsDel:
		return nil, notes.ErrIsDeleted
	}
	return rowToDomain(row), nil
}

func (r *repository) GetMulti(_ context.Context) ([]*notes.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*notes.Note, 0, len(r.noteDB))
	for _, row := range r.noteDB {
		if !row.IsDel {
			list = append(list, rowToDomain(row))
		}
	}
	return list, nil
}
