package notes

import (
	"time"

	"github.com/google/uuid"
)

type noteRow struct {
	UUID      uuid.UUID
	Title     string
	Desc      string
	IsDel     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
