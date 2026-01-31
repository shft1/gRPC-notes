package v1

import (
	"context"

	"github.com/shft1/grpc-notes/observability/logger"
)

type NoteHandler struct {
	sysCtx context.Context
	log    logger.Logger
	noteGW gateway
}

func NewNoteHandler(ctx context.Context, log logger.Logger, noteGW gateway) *NoteHandler {
	return &NoteHandler{
		sysCtx: ctx,
		log:    log,
		noteGW: noteGW,
	}
}
