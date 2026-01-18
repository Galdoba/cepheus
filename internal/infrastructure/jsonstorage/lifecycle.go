package jsonstorage

import (
	"fmt"
	"os"
	"time"
)

// DestroyStorage permanently deletes the storage file from disk.
// Returns an error if storage is closed, path is not set, or file removal fails.
func (s *Storage[T]) DestroyStorage() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}
	if s.path == "" {
		return fmt.Errorf("storage path is not set")
	}

	if err := os.Remove(s.path); err != nil {
		return fmt.Errorf("failed to remove file: %v", err)
	}

	// Reset all fields
	s.path = ""
	s.Created = time.Time{}
	s.LastUpdated = time.Time{}
	s.Entries = nil
	s.deleted = nil
	s.changed = nil
	s.closed = true

	return nil
}

// Close marks the storage as closed, preventing further operations.
// Returns an error if storage is already closed.
func (s *Storage[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrStorageClosed
	}

	s.closed = true
	return nil
}

// Commit synchronizes in-memory changes with the file on disk.
// It updates only changed entries, adds new ones, and removes deleted ones.
// Evicted entries (removed via Evict) are kept in the file.
func (s *Storage[T]) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.commitUnsafe()
}

// CommitAndClose commits any pending changes and then closes the storage.
// Returns an error if either commit or close operation fails.
func (s *Storage[T]) CommitAndClose() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.commitUnsafe(); err != nil {
		return err
	}
	if s.closed {
		return ErrStorageClosed
	}
	s.closed = true
	return nil
}

// ReConnect reopens a closed storage connection and reloads data from disk.
// Returns an error if storage is not closed, path is not set, or loading fails.
func (s *Storage[T]) ReConnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		return fmt.Errorf("storage is not closed")
	}
	if s.path == "" {
		return fmt.Errorf("path to storage file is not set")
	}

	s.closed = false
	s.resetTracking()
	return s.loadUnsafe()
}

// Load reloads the storage data from the JSON file on disk.
// This method is thread-safe and replaces all in-memory entries.
// Note: This will clear any uncommitted changes.
func (s *Storage[T]) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}

	return s.loadUnsafe()
}

// loadUnsafe performs the load operation without locking. Caller must hold the lock.
func (s *Storage[T]) loadUnsafe() error {
	fd, err := readFileData[T](s.path)
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	s.Created = fd.Created
	s.LastUpdated = fd.LastUpdated
	s.Entries = fd.Entries
	s.resetTracking()

	return nil
}
