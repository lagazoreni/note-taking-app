package sqlite

import "errors"

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("record version conflict")
	ErrNameConflict      = errors.New("name already exists")
	ErrWorkspaceBoundary = errors.New("record belongs to another workspace")
	ErrInvalidLink       = errors.New("invalid note-question link")
)
