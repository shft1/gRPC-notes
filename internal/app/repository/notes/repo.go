package notes

import (
	"sync"

	"github.com/google/uuid"
)

type repository struct {
	mu     sync.RWMutex
	noteDB map[uuid.UUID]*noteRow
}

func NewNoteRepository() *repository {
	return &repository{noteDB: make(map[uuid.UUID]*noteRow)}
}
