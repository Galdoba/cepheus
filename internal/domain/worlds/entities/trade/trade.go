package trade

import (
	"github.com/Galdoba/cepheus/internal/domain/support/entities/paths"
	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/astrogation"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func Exists(source, destination coordinates.Global) (bool, error) {
	as, err := astrogation.New(paths.DefaultExternalDB_File())
	if err != nil {
		return false, err
	}
	sourceWorld, err := as.World(source.DatabaseKey())
	if err != nil {
		return false, err
	}
	destinationWorld, err := as.World(destination.DatabaseKey())
	if err != nil {
		return false, err
	}
	if !as.TradePathExist(source.ToCube(), destination.ToCube()) {
		return false, err
	}
	stc := classifications.Classify(uwp.UWP(sourceWorld.UWP))
	dtc := classifications.Classify(uwp.UWP(destinationWorld.UWP))
	for range 2 {
		for _, tc1 := range stc {
			switch tc1 {
			case classifications.In, classifications.Ht:
				for _, tc2 := range dtc {
					switch tc2 {
					case classifications.As, classifications.De, classifications.Ic, classifications.Ni:
						return true, nil
					}
				}
			case classifications.Hi, classifications.Ri:
				for _, tc2 := range dtc {
					switch tc2 {
					case classifications.Ag, classifications.Ga, classifications.Wa:
						return true, nil
					}
				}
			}
		}
		stc, dtc = dtc, stc
	}

	return false, nil
}
