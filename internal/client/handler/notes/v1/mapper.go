package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
)

func mapError(log logger.Logger, err error) (int, string) {
	switch {
	case errors.Is(err, notes.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, notes.ErrIsDeleted):
		return http.StatusGone, err.Error()
	case errors.Is(err, notes.ErrInvalidUUID):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, notes.ErrInvalidData):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, notes.ErrNoteInternal) ||
		errors.Is(err, notes.ErrNoteResponse) ||
		errors.Is(err, notes.ErrClientInternal):
		return http.StatusInternalServerError, err.Error()
	default:
		log.Error("unknown error", logger.NewField("error", err))
		return http.StatusInternalServerError, notes.ErrServiceInternal.Error()
	}
}

func writeResponse(log logger.Logger, w http.ResponseWriter, st int, data any, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		var msg string
		st, msg = mapError(log, err)
		data = map[string]string{"error": msg}
	}
	w.WriteHeader(st)

	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("failed to parse response data to client", logger.NewField("error", err))
	}
}
