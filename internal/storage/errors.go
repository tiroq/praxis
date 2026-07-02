package storage

import "fmt"

// ErrUnsupportedBackend is returned when an unknown storage backend is requested.
type ErrUnsupportedBackend struct {
	Backend string
}

func (e ErrUnsupportedBackend) Error() string {
	return fmt.Sprintf("unsupported storage backend: %q (supported: memory, sqlite)", e.Backend)
}
