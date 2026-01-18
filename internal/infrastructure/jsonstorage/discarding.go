package jsonstorage

import (
	"fmt"
)

// SelectiveLoad loads specific keys from the file into memory.
// For each key provided:
// - If key exists in file: loads fresh copy into memory (overwrites any changes)
// - If key doesn't exist in file: removes from memory (if present)
// Returns an error if storage is closed or if file reading fails.
func (s *Storage[T]) SelectiveLoad(keys ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.selectiveLoadUnsafe(keys...)
}

// selectiveLoadUnsafe performs SelectiveLoad without locking. Caller must hold the lock.
func (s *Storage[T]) selectiveLoadUnsafe(keys ...string) error {
	if s.closed {
		return ErrStorageClosed
	}

	fd, err := s.safeReadFileData()
	if err != nil {
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Process each requested key
	for _, key := range keys {
		if value, existsInFile := fd.Entries[key]; existsInFile {
			s.Entries[key] = value
		} else {
			delete(s.Entries, key)
		}
		s.clearChangeTracking(key)
	}

	return nil
}

// Discard resets changes for specific keys by reloading them from the file.
// For each key provided:
// - If key has pending changes (created, updated, or deleted): reload from file (or remove if not in file)
// - If key has no pending changes: do nothing
// Returns an error if storage is closed or if file reading fails.
func (s *Storage[T]) Discard(keys ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}

	keysToDiscard := selectKeysToDiscard(s.changed, s.deleted, keys...)
	if len(keysToDiscard) == 0 {
		return nil
	}

	if err := s.selectiveLoadUnsafe(keysToDiscard...); err != nil {
		return fmt.Errorf("failed to discard: %w", err)
	}

	return nil
}

// DiscardAll resets all pending changes by reloading from the file.
// For each key currently in memory:
// - If key has pending changes (created, updated, or deleted): reload from file (or remove if not in file)
// - If key has no pending changes: do nothing (keep current in-memory value)
// Returns an error if storage is closed or if file reading fails.
func (s *Storage[T]) DiscardAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}

	pendingKeys := s.getPendingKeys()
	if len(pendingKeys) == 0 {
		return nil
	}

	return s.selectiveLoadUnsafe(pendingKeys...)
}

// selectKeysToDiscard returns keys that have pending changes from the given list
func selectKeysToDiscard(changed, deleted map[string]bool, keys ...string) []string {
	pendingKeys := make(map[string]bool)

	for _, key := range keys {
		if changed[key] || deleted[key] {
			pendingKeys[key] = true
		}
	}

	if len(pendingKeys) == 0 {
		return nil
	}

	keysToLoad := make([]string, 0, len(pendingKeys))
	for key := range pendingKeys {
		keysToLoad = append(keysToLoad, key)
	}
	return keysToLoad
}
