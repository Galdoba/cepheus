package world

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

type WorldDTO struct {
	Coordinates  [2]int   `json:"coordinates"`
	Name         string   `json:"name,omitempty"`
	MainworldUWP string   `json:"mainworld_uwp,omitempty"`
	TradeCodes   []string `json:"trade_codes,omitempty"`
}

func (w *World) ToDTO() WorldDTO {
	dto := WorldDTO{
		Coordinates:  [2]int{w.coordinates.X(), w.coordinates.Y()},
		MainworldUWP: string(w.UWP()),
		TradeCodes:   classifications.Export(w.tradeCodes...),
	}
	return dto
}

func FromDTO(id string, dto WorldDTO) *World {
	w := World{}
	w.id = id
	w.coordinates = coordinates.NewGlobal(dto.Coordinates[0], dto.Coordinates[1])
	w.mainworldUWP, _ = uwp.New(dto.MainworldUWP)
	w.tradeCodes = classifications.Import(dto.TradeCodes...)
	return &w
}

func (dto WorldDTO) Key() string {
	return fmt.Sprintf("{%v,%v}", dto.Coordinates[0], dto.Coordinates[1])
}
