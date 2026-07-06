package conversationstore

import (
	"errors"
	"fmt"
)

// ErrStoreClosed is returned when an operation is attempted after the store has been closed.
var ErrStoreClosed = errors.New("conversation store is closed")

// ErrDuplicateConversation is returned when attempting to create a conversation with a duplicate ID.
type ErrDuplicateConversation struct {
	ID string
}

func (e ErrDuplicateConversation) Error() string {
	return fmt.Sprintf("conversation with ID %q already exists", e.ID)
}

// ErrConversationNotFound is returned when a conversation cannot be found by ID.
type ErrConversationNotFound struct {
	ID string
}

func (e ErrConversationNotFound) Error() string {
	return fmt.Sprintf("conversation with ID %q not found", e.ID)
}

// ErrDuplicateMessage is returned when attempting to append a message with a duplicate ID.
type ErrDuplicateMessage struct {
	ID string
}

func (e ErrDuplicateMessage) Error() string {
	return fmt.Sprintf("message with ID %q already exists", e.ID)
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
