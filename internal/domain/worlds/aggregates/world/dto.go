package world

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
)

type WorldDTO struct {
	Coordinates [2]int         `json:"coordinates"`
	Name        string         `json:"name,omitempty"`
	Imported    t5ss.WorldData `json:"imported"`
}

func (w *World) ToDTO() WorldDTO {
	dto := WorldDTO{
		Coordinates: [2]int{w.coordinates.X(), w.coordinates.Y()},
		Name:        "",
		Imported:    w.imported,
	}
	return dto
}

func FromDTO(id string, dto WorldDTO) *World {
	w := World{}
	w.id = id
	w.coordinates = coordinates.NewGlobal(dto.Coordinates[0], dto.Coordinates[1])
	w.imported = dto.Imported
	return &w
}

func (dto WorldDTO) DatabaseKey() string {
	return fmt.Sprintf("{%v,%v}", dto.Coordinates[0], dto.Coordinates[1])
}
