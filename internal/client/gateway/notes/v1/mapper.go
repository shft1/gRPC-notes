package v1

import (
	"errors"
	"fmt"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErrorRPC(log logger.Logger, st *status.Status) error {
	switch st.Code() {
	case codes.NotFound:
		return notes.ErrNotFound
	case codes.FailedPrecondition:
		return notes.ErrIsDeleted
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %v", notes.ErrInvalidData, st.Message())
	case codes.Internal:
		return notes.ErrNoteInternal
	default:
		log.Error("unknown error code from note service", logger.NewField("code", st.Code()))
		return notes.ErrNoteInternal
	}
}

func mapError(log logger.Logger, err error) error {
	switch {
	case errors.Is(err, notes.ErrNoteReponse):
		return err
	default:
		log.Error("unknown error from client service", logger.NewField("error", err))
		return notes.ErrClientInternal
	}
}
