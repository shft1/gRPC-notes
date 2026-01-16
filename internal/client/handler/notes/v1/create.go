package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/shared"
)

func (nh *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if token := r.Header.Get("Authorization"); token != "" {
		ctx = context.WithValue(ctx, shared.AuthKey, token)
	}
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

func (nh *NoteHandler) normalizeCreate(noteReq *noteCreateRequest) {
	noteReq.Title = strings.TrimSpace(noteReq.Title)
	noteReq.Desc = strings.TrimSpace(noteReq.Desc)
}

func (nh *NoteHandler) validateCreate(noteReq *noteCreateRequest) error {
	lenTitle := utf8.RuneCountInString(noteReq.Title)
	lenDesc := utf8.RuneCountInString(noteReq.Desc)
	switch {
	case lenTitle == 0:
		return fmt.Errorf("%w: %w", notes.ErrInvalidData, notes.ErrEmptyTitle)
	case lenTitle >= 50:
		return fmt.Errorf("%w: %w", notes.ErrInvalidData, notes.ErrTooLongTitle)
	case lenDesc >= 255:
		return fmt.Errorf("%w: %w", notes.ErrInvalidData, notes.ErrTooLongDesc)
	default:
		return nil
	}
}
