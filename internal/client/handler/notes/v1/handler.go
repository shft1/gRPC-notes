package v1

import "github.com/shft1/grpc-notes/observability/logger"

type NoteHandler struct {
	log    logger.Logger
	noteGW gateway
}

func NewNoteHandler(log logger.Logger, noteGW gateway) *NoteHandler {
	return &NoteHandler{
		log:    log,
		noteGW: noteGW,
	}
}
