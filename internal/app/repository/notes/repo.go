package notes

import (
	"sync"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/observability/logger"
)

type noteRepository struct {
	log    logger.Logger
	mu     sync.RWMutex
	noteDB map[uuid.UUID]*noteRow
}

func NewNoteRepository(log logger.Logger) *noteRepository {
	return &noteRepository{
		log:    log,
		noteDB: make(map[uuid.UUID]*noteRow),
	}
}
