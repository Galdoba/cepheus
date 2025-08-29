package star

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/gametable"
)

type Star struct {
	Type         string   `json:"type,omitempty"`
	SubType      *int     `json:"subtype,omitempty"`
	Class        string   `json:"class,omitempty"`
	Mass         float64  `json:"solar mass,omitempty"`
	Diameter     float64  `json:"solar diameter,omitempty"`
	Luminocity   float64  `json:"luminocity,omitempty"`
	Designation  *string  `json:"designation,omitempty"`
	OrbitN       *float64 `json:"orbit#,omitempty"`
	Eccentricity *float64 `json:"eccentricity,omitempty"`
}

func Generate(dp *dice.Dicepool, knownData ...KnownStarData) (Star, error) {
	st := Star{}
	for _, add := range knownData {
		add(&st)
	}
	if st.Type == "" {
		st.Type = ""
	}
	return st, nil
}

type KnownStarData func(*Star)

func KnownType(sType string) KnownStarData {
	return func(s *Star) {
		s.Type = sType
	}
}

func StarTypeDetermination(dp *dice.Dicepool) (string, error) {
	giantsTable, err := gametable.NewTable("Unusual", "2d6",
		gametable.NewRollResult("8-", "III", nil),
		gametable.NewRollResult("9..10", "II", nil),
		gametable.NewRollResult("11", "Ib", nil),
		gametable.NewRollResult("12+", "Ia", nil),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create giantsTable: %v", err)
	}
	// peculiarTable, err := gametable.NewTable("Unusual", "2d6",
	// 	gametable.NewRollResult("2-", "Black Hole", nil),
	// 	gametable.NewRollResult("3", "Pulsar", nil),
	// 	gametable.NewRollResult("4", "Neutron Star", nil),
	// 	gametable.NewRollResult("5..6", "Nebula", nil),
	// 	gametable.NewRollResult("7..9", "Protostar", nil),
	// 	gametable.NewRollResult("10", "Star Cluster", nil),
	// 	gametable.NewRollResult("11+", "Anomaly", nil),
	// )
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create unusualTable: %v", err)
	// }
	// unusualTable, err := gametable.NewTable("Unusual", "2d6",
	// 	gametable.NewRollResult("2-", "Peculiar", peculiarTable),
	// 	gametable.NewRollResult("3-", "VI", nil),
	// 	gametable.NewRollResult("4", "IV", nil),
	// 	gametable.NewRollResult("5..7", "BD", nil),
	// 	gametable.NewRollResult("8..10", "D", nil),
	// 	gametable.NewRollResult("11", "III", nil),
	// 	gametable.NewRollResult("12+", "Giants", giantsTable),
	// )
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create unusualTable: %v", err)
	// }
	specialTable, err := gametable.NewTable("Special", "2d6",
		gametable.NewRollResult("5-", "VI", nil),
		gametable.NewRollResult("6..8", "IV", nil),
		gametable.NewRollResult("9..10", "III", nil),
		gametable.NewRollResult("11+", "Giants", giantsTable),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create specialTable: %v", err)
	}
	hotTable, err := gametable.NewTable("Hot", "2d6",
		gametable.NewRollResult("9-", "A", nil),
		gametable.NewRollResult("10..11", "B", nil),
		gametable.NewRollResult("12+", "O", nil),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create hotTable: %v", err)
	}
	typeTable, err := gametable.NewTable("Type", "2d6",
		gametable.NewRollResult("2-", "Unusual", specialTable),
		gametable.NewRollResult("3..6", "M", nil),
		gametable.NewRollResult("7..8", "K", nil),
		gametable.NewRollResult("9..10", "G", nil),
		gametable.NewRollResult("11", "F", nil),
		gametable.NewRollResult("12+", "Hot", hotTable),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create typeTable: %v", err)
	}
	starType, err := typeTable.Roll(dp)
	if err != nil {
		return "", fmt.Errorf("table roll failed: %v", err)
	}
	done := false
	sType := fmt.Sprintf("%v", dp.Sum("1d10")-1)
	class := "V"
	for !done {
		switch starType {
		case "O", "B", "A", "F", "G", "K", "M":
			done = true
		case "Ia", "Ib", "II", "III", "IV", "VI":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", fmt.Errorf("table roll failed: %v", err)
			}
		}
	}
	return fmt.Sprintf("%v|%v|%v", starType, sType, class), nil
}
