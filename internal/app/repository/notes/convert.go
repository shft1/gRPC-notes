package notes

import (
	"time"

	"github.com/google/uuid"
	"github.com/shft1/grpc-notes/internal/domain/notes"
)

func domainToRow(id uuid.UUID, createdAt time.Time, n *notes.NoteCreate) *noteRow {
	return &noteRow{
		UUID:      id,
		Title:     n.Title,
		Desc:      n.Desc,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func rowToDomain(row *noteRow) *notes.Note {
	return &notes.Note{
		UUID:      row.UUID,
		Title:     row.Title,
		Desc:      row.Desc,
		IsDel:     row.IsDel,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
