package eventstore

import "fmt"

// ErrDuplicateEvent is returned when attempting to append an event with a duplicate ID.
type ErrDuplicateEvent struct {
	ID string
}

func (e ErrDuplicateEvent) Error() string {
	return fmt.Sprintf("event with ID %q already exists", e.ID)
}

// ErrEventNotFound is returned when an event cannot be found by ID.
type ErrEventNotFound struct {
	ID string
}

func (e ErrEventNotFound) Error() string {
	return fmt.Sprintf("event with ID %q not found", e.ID)
}

// ErrMissingField is returned when a required field is missing or empty.
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return fmt.Sprintf("required field %q is missing or empty", string(e))
}

// ErrInvalidJSON is returned when JSON validation fails.
type ErrInvalidJSON struct {
	Field string
	Err   error
}

func (e ErrInvalidJSON) Error() string {
	return fmt.Sprintf("invalid JSON in field %q: %v", e.Field, e.Err)
}

func (e ErrInvalidJSON) Unwrap() error {
	return e.Err
}
