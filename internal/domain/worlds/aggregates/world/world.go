package world

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

type World struct {
	id string
	// canonicalData t5ss.WorldData
	coordinates coordinates.Global
	name        string
	remarks     string
	bases       []string
	description string

	mainworldUWP uwp.UWP
	tradeCodes   []classifications.Classification
}

func Import(wd t5ss.WorldData) (*World, error) {
	w := World{}
	w.id = fmt.Sprintf("{%v,%v}", wd.WorldX, wd.WorldY)
	// w.canonicalData = wd
	w.name = wd.Name
	w.coordinates = coordinates.NewGlobal(wd.WorldX, wd.WorldY)
	if w.mainworldUWP == "" {
		mwUWP, err := uwp.New(wd.UWP)
		if err != nil {
			if err != nil {
				switch strings.Contains(err.Error(), "fields undefined") {
				case true:
				case false:
					return nil, fmt.Errorf("failed to import uwp: %v", err)
				}
			}
		}
		w.mainworldUWP = mwUWP
	}
	if w.name == "" {
		w.name = wd.NormalizeName()
	}
	if w.remarks == "" {
		w.remarks = wd.Remarks
	}
	if len(w.bases) == 0 {
		w.bases = wd.ConfirmedBases()
	}
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
