package jsonstorage

import (
	"errors"
	"sync"
	"time"
)

// Common errors returned by storage operations
var (
	ErrStorageClosed = errors.New("storage is closed")
	ErrEmptyKey      = errors.New("empty key provided")
	ErrKeyExists     = errors.New("key already exists")
	ErrKeyNotFound   = errors.New("key not found")
)

// storage is a generic thread-safe JSON file-based storage structure.
// It stores entries of type T in a map, where keys are strings.
type storage[T any] struct {
	path        string          `json:"-"`            // File path where the JSON data is stored (not serialized)
	Created     time.Time       `json:"created"`      // Time storage was created
	LastUpdated time.Time       `json:"last_updated"` // Time storage was last modified
	mu          sync.Mutex      `json:"-"`            // Mutex to ensure thread-safe access (not serialized)
	Entries     map[string]T    `json:"entries"`      // Map storing all key-value pairs
	closed      bool            `json:"-"`            // Flag indicating if storage is closed (not serialized)
	deleted     map[string]bool `json:"-"`            // Keys that were deleted since last commit
	changed     map[string]bool `json:"-"`            // Keys that were changed since last commit
}

// NewStorage creates a new empty JSON storage file at the specified path.
// Returns an error if a file already exists at the path or if file creation fails.
func NewStorage[T any](path string) (*storage[T], error) {
	return newStorage[T](path)
}

// OpenStorage loads an existing JSON storage file from the specified path.
// Returns an error if the file cannot be read or contains invalid JSON.
func OpenStorage[T any](path string) (*storage[T], error) {
	return openStorage[T](path)
}
