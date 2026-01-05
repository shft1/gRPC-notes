package notes

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (r *noteRepository) Create(_ context.Context, n *notes.NoteCreate) (*notes.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.genUUID()
	createdAt := time.Now()
	noteRow := domainToRow(id, createdAt, n)
	r.noteDB[id] = noteRow
	return rowToDomain(noteRow), nil
}

func (r *noteRepository) genUUID() uuid.UUID {
	id := uuid.New()
	_, ok := r.noteDB[id]
	for ok {
		id = uuid.New()
		_, ok = r.noteDB[id]
	}
	return id
}
