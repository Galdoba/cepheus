package trade

import (
	"fmt"
	"slices"

	"github.com/Galdoba/cepheus/internal/domain/support/entities/paths"
	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/astrogation"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/tradegoods"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func Exists(source, destination coordinates.Global) (bool, error) {
	if source.DatabaseKey() == destination.DatabaseKey() {
		return false, nil
	}
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
	if sourceWorld.Zone == "R" || destinationWorld.Zone == "R" {
		return false, nil
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

func CalculateImport(localCrd, distantCrd coordinates.Global) ([]string, error) {
	if coordinates.Equal(localCrd, distantCrd) {
		return nil, fmt.Errorf("self-trade not implemented")
	}
	as, err := astrogation.New(paths.DefaultExternalDB_File())
	if err != nil {
		return nil, err
	}

	local, err := as.World(localCrd.DatabaseKey())
	if err != nil {
		return nil, err
	}
	distant, err := as.World(distantCrd.DatabaseKey())
	if err != nil {
		return nil, err
	}
	importing := []string{}
	localTC := classifications.Classify(uwp.UWP(local.UWP))
	distantGoods := tradegoods.Available(classifications.Classify(uwp.UWP(distant.UWP))...)
	for _, code := range distantGoods {
		good, _ := tradegoods.New(code)
		best := classifications.Classification("")
		for _, tc := range localTC {
			if dm, ok := good.SaleDM[tc]; ok {
				if dm > good.SaleDM[best] {
					best = tc
				}
			}
		}
		if best != "" {
			importing = append(importing, code)
		}

	}

	return clearSlice(importing), nil
}

func clearSlice(sl []string) []string {
	newSl := []string{}
	for _, s := range sl {
		if slices.Contains(newSl, s) {
			continue
		}
		newSl = append(newSl, s)
	}
	return sl
}
