package v1

import (
	"errors"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(log logger.Logger, err error) error {
	switch {
	case errors.Is(err, notes.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, notes.ErrIsDeleted):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, notes.ErrEmptyTitle) ||
		errors.Is(err, notes.ErrTooLongTitle) ||
		errors.Is(err, notes.ErrTooLongDesc) ||
		errors.Is(err, notes.ErrInvalidUUID):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		log.Error("unknown service error", logger.NewField("error", err))
		return status.Error(codes.Internal, notes.ErrInternal.Error())
	}
}
