package notes

import (
	"context"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (nr *noteRepository) GetByID(_ context.Context, id uuid.UUID) (*notes.Note, error) {
	nr.mu.RLock()
	defer nr.mu.RUnlock()

	row, ok := nr.noteDB[id]
	switch {
	case !ok:
		return nil, notes.ErrNotFound
	case row.IsDel:
		return nil, notes.ErrIsDeleted
	}
	return rowToDomain(row), nil
}

func (nr *noteRepository) GetMulti(_ context.Context) ([]*notes.Note, error) {
	nr.mu.RLock()
	defer nr.mu.RUnlock()

	list := make([]*notes.Note, 0, len(nr.noteDB))
	for _, row := range nr.noteDB {
		if !row.IsDel {
			list = append(list, rowToDomain(row))
		}
	}
	return list, nil
}
