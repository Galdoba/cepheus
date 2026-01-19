package jsonstorage

import (
	"fmt"
)

// Create adds a new entry with the specified key.
// Returns an error if the key is empty, already exists, or storage is closed.
func (s *Storage[T]) Create(key string, entry T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}
	if key == "" {
		return ErrEmptyKey
	}
	if _, ok := s.Entries[key]; ok {
		return fmt.Errorf("%w: %v", ErrKeyExists, key)
	}

	s.Entries[key] = entry
	s.markChanged(key)
	return nil
}

// Read retrieves an entry by its key.
// Returns the entry if found, otherwise returns a zero value and an error.
func (s *Storage[T]) Read(key string) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var noValue T
	if s.closed {
		return noValue, ErrStorageClosed
	}
	if key == "" {
		return noValue, ErrEmptyKey
	}

	if val, ok := s.Entries[key]; ok {
		return val, nil
	}
	return noValue, ErrKeyNotFound
}

// Update modifies an existing entry with the specified key.
// Returns an error if the key is empty, doesn't exist, or storage is closed.
func (s *Storage[T]) Update(key string, newEntry T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}
	if key == "" {
		return ErrEmptyKey
	}
	if _, ok := s.Entries[key]; !ok {
		return fmt.Errorf("%w: %v", ErrKeyNotFound, key)
	}

	s.Entries[key] = newEntry
	s.markChanged(key)
	return nil
}

// Delete removes an entry with the specified key.
// Deleted entries will be removed from the file on next commit.
// Returns an error if the key is empty, doesn't exist, or storage is closed.
func (s *Storage[T]) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}
	if key == "" {
		return ErrEmptyKey
	}
	if _, ok := s.Entries[key]; !ok {
		return fmt.Errorf("%w: %v", ErrKeyNotFound, key)
	}

	delete(s.Entries, key)
	s.markDeleted(key)
	return nil
}

func (s *Storage[T]) ReadAll() map[string]T {
	return s.Entries
}

// AllKeys returns a slice containing all keys in the storage.
// The order of keys is not guaranteed.
func (s *Storage[T]) AllKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := []string{}
	for key := range s.Entries {
		keys = append(keys, key)
	}
	return keys
}

// AllEntries returns a slice containing all values in the storage.
// The order of values is not guaranteed.
func (s *Storage[T]) AllEntries() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries := []T{}
	for _, entry := range s.Entries {
		entries = append(entries, entry)
	}
	return entries
}

// Len returns the number of entries currently stored.
func (s *Storage[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Entries)
}

// Path returns the file path where the JSON data is stored.
func (s *Storage[T]) Path() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.path
}

// Evict removes all entries except those specified in the exceptions list.
// Evicted entries are removed from memory only and will not be changed on commit.
// Returns a slice of keys that were evicted.
// Returns an error if storage is closed.
func (s *Storage[T]) Evict(exceptions ...string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, ErrStorageClosed
	}

	exceptionSet := make(map[string]bool, len(exceptions))
	for _, e := range exceptions {
		exceptionSet[e] = true
	}

	evicted := []string{}
	for key := range s.Entries {
		if !exceptionSet[key] {
			delete(s.Entries, key)
			delete(s.changed, key) // Remove from changed tracking
			evicted = append(evicted, key)
		}
	}

	return evicted, nil
}
