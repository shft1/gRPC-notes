package v1

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (nh *noteHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
