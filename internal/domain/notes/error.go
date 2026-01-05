package notes

import "errors"

var (
	ErrNotFound     = errors.New("note not found")
	ErrIsDeleted    = errors.New("note has been deleted")
	ErrEmptyTitle   = errors.New("required note fields aren't filled")
	ErrTooLongTitle = errors.New("note title is too long")
	ErrTooLongDesc  = errors.New("note description is too long")
	ErrInvalidUUID  = errors.New("invalid note uuid")
	ErrInternal     = errors.New("internal error in note service")
)
