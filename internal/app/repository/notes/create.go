package notes

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (nr *noteRepository) Create(_ context.Context, n *notes.NoteCreate) (*notes.Note, error) {
	nr.mu.Lock()
	defer nr.mu.Unlock()

	id := nr.genUUID()
	createdAt := time.Now()
	noteRow := domainToRow(id, createdAt, n)
	nr.noteDB[id] = noteRow
	return rowToDomain(noteRow), nil
}

func (nr *noteRepository) genUUID() uuid.UUID {
	id := uuid.New()
	_, ok := nr.noteDB[id]
	for ok {
		id = uuid.New()
		_, ok = nr.noteDB[id]
	}
	return id
}
