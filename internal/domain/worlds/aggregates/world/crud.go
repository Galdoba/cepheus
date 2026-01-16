package world

import (
	"errors"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/domain/support/entities/paths"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
)

var dbPath = paths.WorldsStoragePath()

func (w *World) Create() error {
	storage, err := openStorage()
	if err != nil {
		return fmt.Errorf("failed to open storage: %v", err)
	}
	if err := storage.Create(w.id, w.ToDTO()); err != nil {
		return fmt.Errorf("failed to create world entry (%v): %v", w.id, err)
	}
	if err := storage.CommitAndClose(); err != nil {
		return fmt.Errorf("failed to close storage: %v", err)
	}
	return nil
}

func Read(key string) (*World, error) {
	storage, err := openStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %v", err)
	}
	dto, err := storage.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to world entry (%v): %v", key, err)
	}
	if err := storage.CommitAndClose(); err != nil {
		return nil, fmt.Errorf("failed to close storage: %v", err)
	}
	return FromDTO(key, dto), nil
}

func (w *World) Update() error {
	storage, err := openStorage()
	if err != nil {
		return fmt.Errorf("failed to open storage: %v", err)
	}
	if err := storage.Update(w.id, w.ToDTO()); err != nil {
		return fmt.Errorf("failed to update world entry (%v): %v", w.id, err)
	}
	if err := storage.CommitAndClose(); err != nil {
		return fmt.Errorf("failed to close storage: %v", err)
	}
	return nil
}

func (w *World) Delete() error {
	storage, err := openStorage()
	if err != nil {
		return fmt.Errorf("failed to open storage: %v", err)
	}
	if err := storage.Delete(w.id); err != nil {
		return fmt.Errorf("failed to delete world entry (%v): %v", w.id, err)
	}
	if err := storage.CommitAndClose(); err != nil {
		return fmt.Errorf("failed to close storage: %v", err)
	}
	return nil
}

type storage interface {
	Create(string, WorldDTO) error
	Read(string) (WorldDTO, error)
	Update(string, WorldDTO) error
	Delete(string) error
	CommitAndClose() error
}

func openStorage() (storage, error) {
	js, err := jsonstorage.OpenStorage[WorldDTO](dbPath)
	if err != nil {
		switch errors.Is(err, os.ErrNotExist) {
		case true:
			fmt.Printf("world storage does not exits!\ncreate new: ")
			js, err = jsonstorage.NewStorage[WorldDTO](dbPath)
			if err != nil {
				fmt.Println("failed!")
				fmt.Println("aborting program...")
				return nil, err
			}
		case false:
			return nil, fmt.Errorf("failed to create new storage: %v", err)
		}
	}
	return js, nil
}
