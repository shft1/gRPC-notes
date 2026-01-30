package notes

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID
	Title     string
	Desc      string
	IsDel     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteCreate struct {
	Title string
	Desc  string
}

type NoteEvent struct {
	ID    uuid.UUID
	Title string
}
