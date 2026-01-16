package v1

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/shared"
)

func (nh *NoteHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if token := r.Header.Get("Authorization"); token != "" {
		ctx = context.WithValue(ctx, shared.AuthKey, token)
	}
	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil || id == uuid.Nil {
		writeResponse(nh.log, w, 0, nil, notes.ErrInvalidUUID)
		return
	}
	note, err := nh.noteGW.DeleteByID(ctx, id)
	if err != nil {
		writeResponse(nh.log, w, 0, nil, err)
		return
	}
	writeResponse(nh.log, w, http.StatusOK, toDTOResponse(note), nil)
}
