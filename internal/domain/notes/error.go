package notes

import "errors"

var (
	ErrNotFound  = errors.New("note not found")
	ErrIsDeleted = errors.New("note has been deleted")
)
