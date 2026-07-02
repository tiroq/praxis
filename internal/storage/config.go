package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// Backend names for storage configuration.
const (
	BackendMemory = "memory"
	BackendSQLite = "sqlite"
)

// Config specifies the storage backend and its parameters.
type Config struct {
	// Backend is the storage implementation to use: "memory" or "sqlite"
	Backend string

	// SQLitePath is the file path for SQLite database.
	// Only used when Backend is "sqlite".
	// Empty value defaults to "build/praxis.db" in local development.
	SQLitePath string
}

// DefaultConfig returns the default storage configuration for local development.
// Uses SQLite with a local build/ directory database file.
func DefaultConfig() Config {
	return Config{
		Backend:    BackendSQLite,
		SQLitePath: "build/praxis.db",
	}
}

// ConfigFromEnv constructs a Config by reading environment variables.
// Supported environment variables:
//   - PRAXIS_STORAGE_BACKEND: "memory" or "sqlite" (default: "sqlite")
//   - PRAXIS_SQLITE_PATH: path to SQLite database file (default: "build/praxis.db")
//
// If environment variables are not set, returns DefaultConfig().
func ConfigFromEnv() Config {
	backend := os.Getenv("PRAXIS_STORAGE_BACKEND")
	if backend == "" {
		backend = BackendSQLite
	}

	sqlitePath := os.Getenv("PRAXIS_SQLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "build/praxis.db"
	}

	return Config{
		Backend:    backend,
		SQLitePath: sqlitePath,
	}
}

// Validate checks that the Config has valid values.
func (c Config) Validate() error {
	if c.Backend != BackendMemory && c.Backend != BackendSQLite {
		return ErrUnsupportedBackend{Backend: c.Backend}
	}

	if c.Backend == BackendSQLite && c.SQLitePath == "" {
		return fmt.Errorf("SQLitePath must not be empty when backend is sqlite")
	}

	return nil
}

// EnsureSQLiteDir creates the directory for the SQLite database file if it doesn't exist.
// This is called internally by Open when using SQLite backend.
func (c Config) ensureSQLiteDir() error {
	if c.Backend != BackendSQLite {
		return nil
	}

	// Skip directory creation for special SQLite paths
	if c.SQLitePath == ":memory:" {
		return nil
	}

	dir := filepath.Dir(c.SQLitePath)
	if dir == "." || dir == "" {
		return nil
	}

	return os.MkdirAll(dir, 0o755)
}
