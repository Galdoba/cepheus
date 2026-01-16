package world

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

type World struct {
	id           string
	coordinates  coordinates.Global
	name         string
	mainworldUWP uwp.UWP
	tradeCodes   []classifications.Classification
}

func Import(wd t5ss.WorldData) (*World, error) {
	w := World{}
	w.id = fmt.Sprintf("{%v,%v}", wd.WorldX, wd.WorldY)
	w.name = wd.Name
	w.coordinates = coordinates.NewGlobal(wd.WorldX, wd.WorldY)
	mwUWP, err := uwp.New(wd.UWP)
	if err != nil {
		fmt.Println(err)
	}
	w.mainworldUWP = mwUWP
	return &w, nil
}

func (w *World) Coordinates() coordinates.Cube {
	return w.coordinates.ToCube()
}

func (w *World) UWP() uwp.UWP {
	return w.mainworldUWP
}

func (w *World) Hydrate() error {
	// fill world with custom data
	return nil
}
