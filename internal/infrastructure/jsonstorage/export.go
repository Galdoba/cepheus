package jsonstorage

import (
	"encoding/json"
	"fmt"
	"io"
)

// Export writes a formatted JSON representation of the storage to the provided writer.
// This is useful for backups or external inspection without saving to the original file.
func (s *storage[T]) Export(w io.Writer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}

	fd := &fileData[T]{
		Created:     s.Created,
		LastUpdated: s.LastUpdated,
		Entries:     s.Entries,
	}

	data, err := json.MarshalIndent(fd, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if _, err = w.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}
	return nil
}