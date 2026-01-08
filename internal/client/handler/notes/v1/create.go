package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func (nh *noteHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var noteReq noteCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&noteReq); err != nil {
		writeResponse(nh.log, w, 0, nil, notes.ErrInvalidData)
		return
	}
	nh.normalizeCreate(&noteReq)
	if err := nh.validateCreate(&noteReq); err != nil {
		writeResponse(nh.log, w, 0, nil, err)
		return
	}
	note, err := nh.noteGW.Create(ctx, toDomainCreate(&noteReq))
	if err != nil {
		writeResponse(nh.log, w, 0, nil, err)
		return
	}
	writeResponse(nh.log, w, http.StatusCreated, toDTOResponse(note), nil)
}

func (nh *noteHandler) normalizeCreate(noteReq *noteCreateRequest) {
	noteReq.Title = strings.TrimSpace(noteReq.Title)
	noteReq.Desc = strings.TrimSpace(noteReq.Desc)
}

func (nh *noteHandler) validateCreate(noteReq *noteCreateRequest) error {
	lenTitle := utf8.RuneCountInString(noteReq.Title)
	lenDesc := utf8.RuneCountInString(noteReq.Desc)
	switch {
	case lenTitle == 0:
		return notes.ErrEmptyTitle
	case lenTitle > 50:
		return notes.ErrTooLongTitle
	case lenDesc > 300:
		return notes.ErrTooLongDesc
	default:
		return nil
	}
}
