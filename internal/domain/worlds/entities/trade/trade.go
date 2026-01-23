package trade

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/astrogation"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func Exists(sourceWorld, destinationWorld t5ss.WorldData) (bool, error) {
	if sourceWorld.Zone == "R" || destinationWorld.Zone == "R" {
		return false, nil
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
