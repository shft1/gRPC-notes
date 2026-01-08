package v1

import "github.com/shft1/grpc-notes/observability/logger"

type noteHandler struct {
	log    logger.Logger
	noteGW gateway
}

func NewNoteHandler(log logger.Logger, noteGW gateway) *noteHandler {
	return &noteHandler{
		log:    log,
		noteGW: noteGW,
	}
}
