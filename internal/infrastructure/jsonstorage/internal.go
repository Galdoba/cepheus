package jsonstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// fileData represents the structure stored in the JSON file
type fileData[T any] struct {
	Created     time.Time    `json:"created"`
	LastUpdated time.Time    `json:"last_updated"`
	Entries     map[string]T `json:"entries"`
}

// newStorage creates a new empty JSON storage file at the specified path.
func newStorage[T any](path string) (*Storage[T], error) {
	// Try to create the file atomically, failing if it exists.
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("path %v exists", path)
		}
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	s := &Storage[T]{
		path:        path,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Entries:     make(map[string]T),
		deleted:     make(map[string]bool),
		changed:     make(map[string]bool),
	}

	// Write initial empty state
	fd := &fileData[T]{
		Created:     s.Created,
		LastUpdated: s.LastUpdated,
		Entries:     s.Entries,
	}

	if err := writeFileData(file, fd); err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("failed to save initial state: %w", err)
	}

	return s, nil
}

// openStorage loads an existing JSON storage file from the specified path.
func openStorage[T any](path string) (*Storage[T], error) {
	fd, err := readFileData[T](path)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	s := &Storage[T]{
		path:        path,
		Created:     fd.Created,
		LastUpdated: fd.LastUpdated,
		Entries:     fd.Entries,
		closed:      false,
		deleted:     make(map[string]bool),
		changed:     make(map[string]bool),
	}
	return s, nil
}

// commitUnsafe performs the commit operation without locking. Caller must hold the lock.
func (s *Storage[T]) commitUnsafe() error {
	if s.closed {
		return ErrStorageClosed
	}

	fd, err := s.safeReadFileData()
	if err != nil {
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Apply deletions from memory to file
	for key := range s.deleted {
		delete(fd.Entries, key)
	}

	// Apply changes (creates and updates) from memory to file
	for key := range s.changed {
		if entry, ok := s.Entries[key]; ok {
			fd.Entries[key] = entry
		}
	}

	// Update timestamps if changes were made
	if len(s.deleted) > 0 || len(s.changed) > 0 {
		s.LastUpdated = time.Now()
		fd.LastUpdated = s.LastUpdated
	}

	if err := atomicWriteFileData(s.path, fd); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	s.resetTracking()

	return nil
}

// markChanged marks a key as changed and updates timestamp
func (s *Storage[T]) markChanged(key string) {
	s.changed[key] = true
	delete(s.deleted, key)
	s.LastUpdated = time.Now()
}

// markDeleted marks a key as deleted and updates timestamp
func (s *Storage[T]) markDeleted(key string) {
	s.deleted[key] = true
	delete(s.changed, key)
	s.LastUpdated = time.Now()
}

// clearChangeTracking clears all change tracking for a key
func (s *Storage[T]) clearChangeTracking(key string) {
	delete(s.changed, key)
	delete(s.deleted, key)
}

// clearChangeTracking clears all change tracking for a key
func (s *Storage[T]) resetTracking() {
	s.deleted = make(map[string]bool)
	s.changed = make(map[string]bool)
}

// getPendingKeys returns a slice of keys that have pending changes
func (s *Storage[T]) getPendingKeys() []string {
	pendingKeys := make([]string, 0, len(s.changed)+len(s.deleted))

	for key := range s.changed {
		pendingKeys = append(pendingKeys, key)
	}

	for key := range s.deleted {
		if !s.changed[key] {
			pendingKeys = append(pendingKeys, key)
		}
	}

	return pendingKeys
}

// safeReadFileData reads file data safely, returning empty data if file doesn't exist
func (s *Storage[T]) safeReadFileData() (*fileData[T], error) {
	if _, err := os.Stat(s.path); err != nil {
		// File doesn't exist, return empty data with current timestamps
		return &fileData[T]{
			Created:     s.Created,
			LastUpdated: s.LastUpdated,
			Entries:     make(map[string]T),
		}, nil
	}

	return readFileData[T](s.path)
}

// readFileData reads and parses the storage file
func readFileData[T any](path string) (*fileData[T], error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var fd fileData[T]
	if err := json.Unmarshal(data, &fd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return &fd, nil
}

// writeFileData writes data to an io.Writer
func writeFileData[T any](w io.Writer, fd *fileData[T]) error {
	data, err := json.Marshal(fd)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}

	return nil
}

// writeFileDataToPath writes data to a file at the specified path
// func writeFileDataToPath[T any](path string, fd *fileData[T]) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file: %v", err)
// 	}
// 	defer file.Close()

// 	return writeFileData(file, fd)
// }

func atomicWriteFileData[T any](path string, fd *fileData[T]) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	dir := filepath.Dir(absPath)
	base := filepath.Base(absPath)

	// Create temp file with random name to avoid collisions
	tempFile, err := os.CreateTemp(dir, "."+base+".*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()

	// Clean up temp file if we fail
	defer func() {
		if err != nil {
			os.Remove(tempPath)
		}
	}()

	// Write data
	if err := writeFileData(tempFile, fd); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	// Sync to ensure data is written to disk
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to sync temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %v", err)
	}

	// Atomic rename (atomic on most filesystems)
	if err := os.Rename(tempPath, absPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %v", err)
	}

	return nil
}
