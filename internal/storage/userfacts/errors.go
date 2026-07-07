package userfacts

import (
	"errors"
	"fmt"
)

// ErrStoreClosed is returned when an operation is attempted after the store is closed.
var ErrStoreClosed = errors.New("user facts store is closed")

// ErrFactNotFound is returned when a candidate fact ID does not exist.
type ErrFactNotFound struct {
	ID string
}

func (e ErrFactNotFound) Error() string {
	return fmt.Sprintf("candidate fact with ID %q was not found", e.ID)
}

// ErrDuplicateFact is returned when appending a fact with a duplicate ID.
type ErrDuplicateFact struct {
	ID string
}

func (e ErrDuplicateFact) Error() string {
	return fmt.Sprintf("candidate fact with ID %q already exists", e.ID)
}

// ErrMissingField is returned when a required field is missing.
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return fmt.Sprintf("required field %q is missing or empty", string(e))
}

// ErrInvalidConfidence is returned when a fact confidence is outside [0,1].
type ErrInvalidConfidence float64

func (e ErrInvalidConfidence) Error() string {
	return fmt.Sprintf("confidence %.4f is outside [0,1]", float64(e))
}

// ErrInvalidValidationState is returned when an unknown validation state is used.
type ErrInvalidValidationState struct {
	State string
}

func (e ErrInvalidValidationState) Error() string {
	return fmt.Sprintf("validation state %q is invalid", e.State)
}

// ErrInvalidValidationTransition is returned for illegal validation lifecycle moves.
type ErrInvalidValidationTransition struct {
	From ValidationState
	To   ValidationState
}

func (e ErrInvalidValidationTransition) Error() string {
	return fmt.Sprintf("cannot transition validation state from %q to %q", e.From, e.To)
}
