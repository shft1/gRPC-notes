package notes

import "errors"

var (
	ErrNotFound       = errors.New("note not found")
	ErrIsDeleted      = errors.New("note has been deleted")
	ErrEmptyTitle     = errors.New("required note fields aren't filled")
	ErrTooLongTitle   = errors.New("note title is too long")
	ErrTooLongDesc    = errors.New("note description is too long")
	ErrInvalidUUID    = errors.New("invalid note uuid")
	ErrNoteInternal   = errors.New("internal error in note service")
	ErrClientInternal = errors.New("internal error in client service")
	ErrInvalidData    = errors.New("invalid transmitted data")
	ErrNoteReponse    = errors.New("returned data from note service are not expected")
)
