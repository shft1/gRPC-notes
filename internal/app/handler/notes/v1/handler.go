package v1

import (
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
)

type NoteHandler struct {
	pb.UnimplementedNoteAPIServer
	log         logger.Logger
	bus         eventBus
	noteUsecase noteUsecase
}

func NewNoteHandler(log logger.Logger, bus eventBus, noteUsecase noteUsecase) *NoteHandler {
	return &NoteHandler{
		log:         log,
		bus:         bus,
		noteUsecase: noteUsecase,
	}
}
