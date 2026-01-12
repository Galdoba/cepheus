package jsonstorage

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"os"
// 	"sync"
// 	"time"
// )

// // Common errors returned by storage operations
// var (
// 	ErrStorageClosed = errors.New("storage is closed")
// 	ErrEmptyKey      = errors.New("empty key provided")
// 	ErrKeyExists     = errors.New("key already exists")
// 	ErrKeyNotFound   = errors.New("key not found")
// )

// // fileData represents the structure stored in the JSON file
// type fileData[T any] struct {
// 	Created     time.Time    `json:"created"`
// 	LastUpdated time.Time    `json:"last_updated"`
// 	Entries     map[string]T `json:"entries"`
// }

// // storage is a generic thread-safe JSON file-based storage structure.
// // It stores entries of type T in a map, where keys are strings.
// type storage[T any] struct {
// 	path        string          `json:"-"`            // File path where the JSON data is stored (not serialized)
// 	Created     time.Time       `json:"created"`      // Time storage was created
// 	LastUpdated time.Time       `json:"last_updated"` // Time storage was last modified
// 	mu          sync.Mutex      `json:"-"`            // Mutex to ensure thread-safe access (not serialized)
// 	Entries     map[string]T    `json:"entries"`      // Map storing all key-value pairs
// 	closed      bool            `json:"-"`            // Flag indicating if storage is closed (not serialized)
// 	deleted     map[string]bool `json:"-"`            // Keys that were deleted since last commit
// 	changed     map[string]bool `json:"-"`            // Keys that were changed since last commit
// }

// // NewStorage creates a new empty JSON storage file at the specified path.
// // Returns an error if a file already exists at the path or if file creation fails.
// func NewStorage[T any](path string) (*storage[T], error) {
// 	// Try to create the file atomically, failing if it exists.
// 	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
// 	if err != nil {
// 		if os.IsExist(err) {
// 			return nil, fmt.Errorf("path %v exists", path)
// 		}
// 		return nil, fmt.Errorf("failed to create file: %v", err)
// 	}
// 	defer file.Close()

// 	s := &storage[T]{
// 		path:        path,
// 		Created:     time.Now(),
// 		LastUpdated: time.Now(),
// 		Entries:     make(map[string]T),
// 		deleted:     make(map[string]bool),
// 		changed:     make(map[string]bool),
// 	}

// 	// Write initial empty state
// 	fd := &fileData[T]{
// 		Created:     s.Created,
// 		LastUpdated: s.LastUpdated,
// 		Entries:     s.Entries,
// 	}

// 	if err := writeFileData(file, fd); err != nil {
// 		os.Remove(path)
// 		return nil, fmt.Errorf("failed to save initial state: %w", err)
// 	}

// 	return s, nil
// }

// // OpenStorage loads an existing JSON storage file from the specified path.
// // Returns an error if the file cannot be read or contains invalid JSON.
// func OpenStorage[T any](path string) (*storage[T], error) {
// 	fd, err := readFileData[T](path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open storage: %w", err)
// 	}

// 	s := &storage[T]{
// 		path:        path,
// 		Created:     fd.Created,
// 		LastUpdated: fd.LastUpdated,
// 		Entries:     fd.Entries,
// 		closed:      false,
// 		deleted:     make(map[string]bool),
// 		changed:     make(map[string]bool),
// 	}
// 	return s, nil
// }

// // DestroyStorage permanently deletes the storage file from disk.
// // Returns an error if storage is closed, path is not set, or file removal fails.
// func (s *storage[T]) DestroyStorage() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}
// 	if s.path == "" {
// 		return fmt.Errorf("storage path is not set")
// 	}

// 	if err := os.Remove(s.path); err != nil {
// 		return fmt.Errorf("failed to remove file: %v", err)
// 	}

// 	// Reset all fields
// 	s.path = ""
// 	s.Created = time.Time{}
// 	s.LastUpdated = time.Time{}
// 	s.Entries = nil
// 	s.deleted = nil
// 	s.changed = nil
// 	s.closed = true

// 	return nil
// }

// // Close marks the storage as closed, preventing further operations.
// // Returns an error if storage is already closed.
// func (s *storage[T]) Close() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	s.closed = true
// 	return nil
// }

// // Commit synchronizes in-memory changes with the file on disk.
// // It updates only changed entries, adds new ones, and removes deleted ones.
// // Evicted entries (removed via Evict) are kept in the file.
// func (s *storage[T]) Commit() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	return s.commitUnsafe()
// }

// // commitUnsafe performs the commit operation without locking. Caller must hold the lock.
// func (s *storage[T]) commitUnsafe() error {
// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	fd, err := s.safeReadFileData()
// 	if err != nil {
// 		return fmt.Errorf("failed to read file data: %w", err)
// 	}

// 	// Apply deletions from memory to file
// 	for key := range s.deleted {
// 		delete(fd.Entries, key)
// 	}

// 	// Apply changes (creates and updates) from memory to file
// 	for key := range s.changed {
// 		if entry, ok := s.Entries[key]; ok {
// 			fd.Entries[key] = entry
// 		}
// 	}

// 	// Update timestamps if changes were made
// 	if len(s.deleted) > 0 || len(s.changed) > 0 {
// 		s.LastUpdated = time.Now()
// 		fd.LastUpdated = s.LastUpdated
// 	}

// 	if err := writeFileDataToPath(s.path, fd); err != nil {
// 		return fmt.Errorf("failed to write file: %w", err)
// 	}

// 	// Reset change tracking
// 	s.deleted = make(map[string]bool)
// 	s.changed = make(map[string]bool)

// 	return nil
// }

// // CommitAndClose commits any pending changes and then closes the storage.
// // Returns an error if either commit or close operation fails.
// func (s *storage[T]) CommitAndClose() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if err := s.commitUnsafe(); err != nil {
// 		return err
// 	}
// 	if s.closed {
// 		return ErrStorageClosed
// 	}
// 	s.closed = true
// 	return nil
// }

// // ReConnect reopens a closed storage connection and reloads data from disk.
// // Returns an error if storage is not closed, path is not set, or loading fails.
// func (s *storage[T]) ReConnect() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if !s.closed {
// 		return fmt.Errorf("storage is not closed")
// 	}
// 	if s.path == "" {
// 		return fmt.Errorf("path to storage file is not set")
// 	}

// 	s.closed = false
// 	s.deleted = make(map[string]bool)
// 	s.changed = make(map[string]bool)
// 	return s.loadUnsafe()
// }

// // Evict removes all entries except those specified in the exceptions list.
// // Evicted entries are removed from memory only and will not be changed on commit.
// // Returns a slice of keys that were evicted.
// // Returns an error if storage is closed.
// func (s *storage[T]) Evict(exceptions ...string) ([]string, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return nil, ErrStorageClosed
// 	}

// 	exceptionSet := make(map[string]bool, len(exceptions))
// 	for _, e := range exceptions {
// 		exceptionSet[e] = true
// 	}

// 	evicted := []string{}
// 	for key := range s.Entries {
// 		if !exceptionSet[key] {
// 			delete(s.Entries, key)
// 			delete(s.changed, key) // Remove from changed tracking
// 			evicted = append(evicted, key)
// 		}
// 	}

// 	return evicted, nil
// }

// // Create adds a new entry with the specified key.
// // Returns an error if the key is empty, already exists, or storage is closed.
// func (s *storage[T]) Create(key string, entry T) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}
// 	if key == "" {
// 		return ErrEmptyKey
// 	}
// 	if _, ok := s.Entries[key]; ok {
// 		return fmt.Errorf("%w: %v", ErrKeyExists, key)
// 	}

// 	s.Entries[key] = entry
// 	s.markChanged(key)
// 	return nil
// }

// // Read retrieves an entry by its key.
// // Returns the entry if found, otherwise returns a zero value and an error.
// func (s *storage[T]) Read(key string) (T, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	var noValue T
// 	if s.closed {
// 		return noValue, ErrStorageClosed
// 	}
// 	if key == "" {
// 		return noValue, ErrEmptyKey
// 	}

// 	if val, ok := s.Entries[key]; ok {
// 		return val, nil
// 	}
// 	return noValue, fmt.Errorf("%w: %v", ErrKeyNotFound, key)
// }

// // Update modifies an existing entry with the specified key.
// // Returns an error if the key is empty, doesn't exist, or storage is closed.
// func (s *storage[T]) Update(key string, newEntry T) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}
// 	if key == "" {
// 		return ErrEmptyKey
// 	}
// 	if _, ok := s.Entries[key]; !ok {
// 		return fmt.Errorf("%w: %v", ErrKeyNotFound, key)
// 	}

// 	s.Entries[key] = newEntry
// 	s.markChanged(key)
// 	return nil
// }

// // Delete removes an entry with the specified key.
// // Deleted entries will be removed from the file on next commit.
// // Returns an error if the key is empty, doesn't exist, or storage is closed.
// func (s *storage[T]) Delete(key string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}
// 	if key == "" {
// 		return ErrEmptyKey
// 	}
// 	if _, ok := s.Entries[key]; !ok {
// 		return fmt.Errorf("%w: %v", ErrKeyNotFound, key)
// 	}

// 	delete(s.Entries, key)
// 	s.markDeleted(key)
// 	return nil
// }

// // AllKeys returns a slice containing all keys in the storage.
// // The order of keys is not guaranteed.
// func (s *storage[T]) AllKeys() []string {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	keys := []string{}
// 	for key := range s.Entries {
// 		keys = append(keys, key)
// 	}
// 	return keys
// }

// // AllEntries returns a slice containing all values in the storage.
// // The order of values is not guaranteed.
// func (s *storage[T]) AllEntries() []T {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	entries := []T{}
// 	for _, entry := range s.Entries {
// 		entries = append(entries, entry)
// 	}
// 	return entries
// }

// // Len returns the number of entries currently stored.
// func (s *storage[T]) Len() int {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	return len(s.Entries)
// }

// // Path returns the file path where the JSON data is stored.
// func (s *storage[T]) Path() string {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	return s.path
// }

// // Export writes a formatted JSON representation of the storage to the provided writer.
// // This is useful for backups or external inspection without saving to the original file.
// func (s *storage[T]) Export(w io.Writer) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	fd := &fileData[T]{
// 		Created:     s.Created,
// 		LastUpdated: s.LastUpdated,
// 		Entries:     s.Entries,
// 	}

// 	data, err := json.MarshalIndent(fd, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal data: %v", err)
// 	}

// 	if _, err = w.Write(data); err != nil {
// 		return fmt.Errorf("failed to write data: %v", err)
// 	}
// 	return nil
// }

// // SelectiveLoad loads specific keys from the file into memory.
// // For each key provided:
// // - If key exists in file: loads fresh copy into memory (overwrites any changes)
// // - If key doesn't exist in file: removes from memory (if present)
// // Returns an error if storage is closed or if file reading fails.
// func (s *storage[T]) SelectiveLoad(keys ...string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	return s.selectiveLoadUnsafe(keys...)
// }

// // selectiveLoadUnsafe performs SelectiveLoad without locking. Caller must hold the lock.
// func (s *storage[T]) selectiveLoadUnsafe(keys ...string) error {
// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	fd, err := s.safeReadFileData()
// 	if err != nil {
// 		return fmt.Errorf("failed to read file data: %w", err)
// 	}

// 	// Process each requested key
// 	for _, key := range keys {
// 		if value, existsInFile := fd.Entries[key]; existsInFile {
// 			s.Entries[key] = value
// 		} else {
// 			delete(s.Entries, key)
// 		}
// 		s.clearChangeTracking(key)
// 	}

// 	return nil
// }

// // Discard resets changes for specific keys by reloading them from the file.
// // For each key provided:
// // - If key has pending changes (created, updated, or deleted): reload from file (or remove if not in file)
// // - If key has no pending changes: do nothing
// // Returns an error if storage is closed or if file reading fails.
// func (s *storage[T]) Discard(keys ...string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	keysToDiscard := selectKeysToDiscard(s.changed, s.deleted, keys...)
// 	if len(keysToDiscard) == 0 {
// 		return nil
// 	}

// 	if err := s.selectiveLoadUnsafe(keysToDiscard...); err != nil {
// 		return fmt.Errorf("failed to discard: %w", err)
// 	}

// 	return nil
// }

// // DiscardAll resets all pending changes by reloading from the file.
// // For each key currently in memory:
// // - If key has pending changes (created, updated, or deleted): reload from file (or remove if not in file)
// // - If key has no pending changes: do nothing (keep current in-memory value)
// // Returns an error if storage is closed or if file reading fails.
// func (s *storage[T]) DiscardAll() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	pendingKeys := s.getPendingKeys()
// 	if len(pendingKeys) == 0 {
// 		return nil
// 	}

// 	return s.selectiveLoadUnsafe(pendingKeys...)
// }

// // Load reloads the storage data from the JSON file on disk.
// // This method is thread-safe and replaces all in-memory entries.
// // Note: This will clear any uncommitted changes.
// func (s *storage[T]) Load() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return ErrStorageClosed
// 	}

// 	return s.loadUnsafe()
// }

// // loadUnsafe performs the load operation without locking. Caller must hold the lock.
// func (s *storage[T]) loadUnsafe() error {
// 	fd, err := readFileData[T](s.path)
// 	if err != nil {
// 		return fmt.Errorf("failed to load data: %w", err)
// 	}

// 	s.Created = fd.Created
// 	s.LastUpdated = fd.LastUpdated
// 	s.Entries = fd.Entries
// 	s.deleted = make(map[string]bool)
// 	s.changed = make(map[string]bool)

// 	return nil
// }

// // markChanged marks a key as changed and updates timestamp
// func (s *storage[T]) markChanged(key string) {
// 	s.changed[key] = true
// 	delete(s.deleted, key)
// 	s.LastUpdated = time.Now()
// }

// // markDeleted marks a key as deleted and updates timestamp
// func (s *storage[T]) markDeleted(key string) {
// 	s.deleted[key] = true
// 	delete(s.changed, key)
// 	s.LastUpdated = time.Now()
// }

// // clearChangeTracking clears all change tracking for a key
// func (s *storage[T]) clearChangeTracking(key string) {
// 	delete(s.changed, key)
// 	delete(s.deleted, key)
// }

// // getPendingKeys returns a slice of keys that have pending changes
// func (s *storage[T]) getPendingKeys() []string {
// 	pendingKeys := make([]string, 0, len(s.changed)+len(s.deleted))

// 	for key := range s.changed {
// 		pendingKeys = append(pendingKeys, key)
// 	}

// 	for key := range s.deleted {
// 		if !s.changed[key] {
// 			pendingKeys = append(pendingKeys, key)
// 		}
// 	}

// 	return pendingKeys
// }

// // safeReadFileData reads file data safely, returning empty data if file doesn't exist
// func (s *storage[T]) safeReadFileData() (*fileData[T], error) {
// 	if _, err := os.Stat(s.path); err != nil {
// 		// File doesn't exist, return empty data with current timestamps
// 		return &fileData[T]{
// 			Created:     s.Created,
// 			LastUpdated: s.LastUpdated,
// 			Entries:     make(map[string]T),
// 		}, nil
// 	}

// 	return readFileData[T](s.path)
// }

// // Helper functions

// // readFileData reads and parses the storage file
// func readFileData[T any](path string) (*fileData[T], error) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read file: %v", err)
// 	}

// 	var fd fileData[T]
// 	if err := json.Unmarshal(data, &fd); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
// 	}

// 	return &fd, nil
// }

// // writeFileData writes data to an io.Writer
// func writeFileData[T any](w io.Writer, fd *fileData[T]) error {
// 	data, err := json.Marshal(fd)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal data: %v", err)
// 	}

// 	_, err = w.Write(data)
// 	if err != nil {
// 		return fmt.Errorf("failed to write data: %v", err)
// 	}

// 	return nil
// }

// // writeFileDataToPath writes data to a file at the specified path
// func writeFileDataToPath[T any](path string, fd *fileData[T]) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file: %v", err)
// 	}
// 	defer file.Close()

// 	return writeFileData(file, fd)
// }

// // selectKeysToDiscard returns keys that have pending changes from the given list
// func selectKeysToDiscard(changed, deleted map[string]bool, keys ...string) []string {
// 	pendingKeys := make(map[string]bool)

// 	for _, key := range keys {
// 		if changed[key] || deleted[key] {
// 			pendingKeys[key] = true
// 		}
// 	}

// 	if len(pendingKeys) == 0 {
// 		return nil
// 	}

// 	keysToLoad := make([]string, 0, len(pendingKeys))
// 	for key := range pendingKeys {
// 		keysToLoad = append(keysToLoad, key)
// 	}
// 	return keysToLoad
// }
