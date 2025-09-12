package star

import (
	"fmt"
	"strconv"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/gametable"
)

func StarTypeAndClassDetermination(dp *dice.Dicepool) (string, string, error) {
	giantsTable, err := gametable.NewTable("Unusual", "2d6",
		gametable.NewRollResult("8-", "III", nil),
		gametable.NewRollResult("9..10", "II", nil),
		gametable.NewRollResult("11", "Ib", nil),
		gametable.NewRollResult("12+", "Ia", nil),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create giantsTable: %v", err)
	}
	specialTable, err := gametable.NewTable("Special", "2d6",
		gametable.NewRollResult("5-", "VI", nil),
		gametable.NewRollResult("6..8", "IV", nil),
		gametable.NewRollResult("9..10", "III", nil),
		gametable.NewRollResult("11+", "Giants", giantsTable),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create specialTable: %v", err)
	}
	hotTable, err := gametable.NewTable("Hot", "2d6",
		gametable.NewRollResult("9-", "A", nil),
		gametable.NewRollResult("10..11", "B", nil),
		gametable.NewRollResult("12+", "O", nil),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create hotTable: %v", err)
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
		return "", "", fmt.Errorf("failed to create typeTable: %v", err)
	}
	starType, err := typeTable.Roll(dp)
	if err != nil {
		return "", "", fmt.Errorf("table roll failed: %v", err)
	}
	done := false
	class := "V"
	for !done {
		switch starType {
		case "O", "B", "A", "F", "G", "K", "M":
			done = true
		case "Ia", "Ib", "II", "III":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", "", fmt.Errorf("table roll failed: %v", err)
			}
		case "IV":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", "", fmt.Errorf("table roll failed: %v", err)
			}
			switch starType {
			case "M":
				starType = "IV"
			case "O":
				starType = "B"
			}

		case "VI":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", "", fmt.Errorf("table roll failed: %v", err)
			}
			switch starType {
			case "F":
				starType = "G"
			case "A":
				starType = "B"
			}
		}
	}
	return starType, class, nil
}

func StarSubTypeDetermination(dp *dice.Dicepool, st Star) (*int, error) {
	subtypeTable := &gametable.GameTable{}
	err := fmt.Errorf("table not created")
	switch st.Type {
	case "M":
		subtypeTable, err = gametable.NewTable("m-type", "2d6",
			gametable.NewRollResult("2", "8", nil),
			gametable.NewRollResult("3", "6", nil),
			gametable.NewRollResult("4", "5", nil),
			gametable.NewRollResult("5", "4", nil),
			gametable.NewRollResult("6", "0", nil),
			gametable.NewRollResult("7", "2", nil),
			gametable.NewRollResult("8", "1", nil),
			gametable.NewRollResult("9", "3", nil),
			gametable.NewRollResult("10", "5", nil),
			gametable.NewRollResult("11", "7", nil),
			gametable.NewRollResult("12", "9", nil),
		)
	case "O", "B", "A", "F", "G", "K":
		subtypeTable, err = gametable.NewTable("numeric", "2d6",
			gametable.NewRollResult("2", "0", nil),
			gametable.NewRollResult("3", "1", nil),
			gametable.NewRollResult("4", "3", nil),
			gametable.NewRollResult("5", "5", nil),
			gametable.NewRollResult("6", "7", nil),
			gametable.NewRollResult("7", "9", nil),
			gametable.NewRollResult("8", "8", nil),
			gametable.NewRollResult("9", "6", nil),
			gametable.NewRollResult("10", "4", nil),
			gametable.NewRollResult("11", "2", nil),
			gametable.NewRollResult("12", "0", nil),
		)
	default:
		return nil, fmt.Errorf("what shaall we do with %v?", st.Type)
	}
	r, err := subtypeTable.Roll(dp)
	if err != nil {
		return nil, fmt.Errorf("sybtype table roll: %v", err)
	}
	n, _ := strconv.Atoi(r)
	return &n, nil
}
