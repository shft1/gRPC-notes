package notes

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (nr *noteRepository) DeleteByID(_ context.Context, uuid uuid.UUID) (*notes.Note, error) {
	nr.mu.Lock()
	defer nr.mu.Unlock()

	row, ok := nr.noteDB[uuid]
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
