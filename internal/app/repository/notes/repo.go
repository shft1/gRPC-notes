package notes

import (
	"sync"

	"github.com/google/uuid"
)

type noteRepository struct {
	mu     sync.RWMutex
	noteDB map[uuid.UUID]*noteRow
}

func NewNoteRepository() *noteRepository {
	return &noteRepository{noteDB: make(map[uuid.UUID]*noteRow)}
}
