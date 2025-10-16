package todo

import "errors"

var (
	// ErrNotFound is returned when a todo is not found.
	ErrNotFound = errors.New("todo not found")

	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict is returned when there's a conflict (e.g., duplicate ID).
	ErrConflict = errors.New("todo already exists")

	// ErrInvalidStatus is returned when a status transition is not allowed.
	ErrInvalidStatus = errors.New("invalid status transition")
)
