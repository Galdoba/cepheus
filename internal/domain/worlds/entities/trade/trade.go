package trade

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/astrogation"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func Exists(source, destination coordinates.Global) (bool, error) {
	fmt.Println(source, destination, "exists?")
	if source.DatabaseKey() == destination.DatabaseKey() {
		return false, nil
	}
	as, err := astrogation.New()
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
	if sourceWorld.Zone == "R" || destinationWorld.Zone == "R" {
		return false, nil
	}
	if !as.TradePathExist(source.ToCube(), destination.ToCube()) {
		return false, err
	}
	stc := classifications.Classify(uwp.UWP(sourceWorld.UWP))
	dtc := classifications.Classify(uwp.UWP(destinationWorld.UWP))
	fmt.Println("maybe...")
	for range 2 {
		for _, tc1 := range stc {
			switch tc1 {
			case classifications.In, classifications.Ht:
				for _, tc2 := range dtc {
					fmt.Println(tc1, tc2)
					switch tc2 {
					case classifications.As, classifications.De, classifications.Ic, classifications.Ni:
						return true, nil
					case classifications.Red:
						return false, nil
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
