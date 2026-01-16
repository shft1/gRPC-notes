package v1

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
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
		details := make([]string, 0, len(st.Details()))
		for _, detail := range st.Details() {
			if msg, ok := detail.(*pb.ErrorDetails); ok {
				details = append(details, msg.Detail)
			}
		}
		return fmt.Errorf("%w: %v", notes.ErrInvalidData, strings.Join(details, "; "))
	case codes.Unauthenticated:
		log.Info(st.Message())
		return fmt.Errorf("%w: %v", notes.ErrUnauthenticated, st.Message())
	case codes.Internal:
		return notes.ErrNoteInternal
	default:
		log.Error("unknown error code from note service", logger.NewField("code", st.Code()))
		return notes.ErrNoteInternal
	}
}

func mapError(log logger.Logger, err error) error {
	switch {
	case errors.Is(err, notes.ErrInvalidUUID):
		return err
	case errors.Is(err, notes.ErrNoteResponse):
		return err
	case errors.Is(err, notes.ErrInvalidData):
		return err
	default:
		log.Error("unknown error from client service", logger.NewField("error", err))
		return notes.ErrClientInternal
	}
}
