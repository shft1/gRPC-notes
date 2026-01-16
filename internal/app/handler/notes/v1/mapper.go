package v1

import (
	"errors"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func mapError(log logger.Logger, err error) error {
	switch {
	case errors.Is(err, notes.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, notes.ErrIsDeleted):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, notes.ErrInvalidData) || errors.Is(err, notes.ErrEmptyTitle) ||
		errors.Is(err, notes.ErrTooLongTitle) || errors.Is(err, notes.ErrTooLongDesc) ||
		errors.Is(err, notes.ErrInvalidUUID):

		detail := &pb.ErrorDetails{Detail: err.Error()}
		st := withDetails(codes.InvalidArgument, notes.ErrInvalidData.Error(), detail)
		return st.Err()

	default:
		log.Error("unknown service error", logger.NewField("error", err))
		return status.Error(codes.Internal, notes.ErrNoteInternal.Error())
	}
}

func withDetails(code codes.Code, msg string, details ...protoadapt.MessageV1) *status.Status {
	st := status.New(code, msg)
	var stErr error
	if st, stErr = st.WithDetails(details...); stErr != nil {
		return status.New(codes.Internal, "failed to put details into error")
	}
	return st
}
