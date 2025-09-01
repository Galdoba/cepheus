package stellarhex

import (
	"fmt"

	"github.com/Galdoba/t5/pkg/dice"
)

func generateDensity(dp *dice.Dicepool) int {
	done := false
	density := Density_Standard
	for !done {
		flux := dp.Flux()
		if flux == 0 {
			done = true
		}
		if flux > 0 {
			density++
		}
		if flux < 0 {
			density--
		}
	}
	return boundInt(density, Density_Void, Density_Core)
}

func boundInt(i, min, max int) int {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}

func (sh *StellarHex) defineStarLikeObjectsPresence(dp *dice.Dicepool) error {
	if sh.StarSystem != "?" {
		return nil
	}
	code, goal := checkGoal(sh.Density)
	if code == "" {
		return fmt.Errorf("failed to process hex density: %v", sh.Density)
	}
	//0 normal stars
	switch getPresence(dp, code, goal) {
	case "+":
		sh.StarSystem = "+"
		return nil
	case "-":
		sh.StarSystem = "-"
		sh.PSR = "-"
	}
	//1 Black hole/Neutron star
	if dp.Check("3d6", "18") {
		switch getPresence(dp, code, goal) {
		case "+":
			switch dp.Check("2d6", "11+") {
			case true:
				sh.BH = "+"
				sh.NS = "-"
				return nil
			case false:
				sh.BH = "-"
				sh.NS = "+"
				return nil
			}
		case "-":
			sh.BH = "-"
			sh.NS = "-"
		}
	}
	//2 White Dwarf
	switch dp.Check("1d6", "6") {
	case false:
		sh.D = "-"
		return nil
	case true:
		sh.D = getPresence(dp, code, goal)
	}
	if sh.D == "+" {
		return nil
	}
	sh.BD = getPresence(dp, code, goal)
	return nil
}

func (sh *StellarHex) defineRoguePlanetsPresence(dp *dice.Dicepool) error {
	code, goal := checkGoal(sh.Density)
	if code == "" {
		return fmt.Errorf("failed to process hex density: %v", sh.Density)
	}
	sh.LGG = getPresence(dp, code, goal)
	for i := 0; i < 2; i++ {
		sh.MGG[i] = getPresence(dp, code, goal)
	}
	for i := 0; i < 3; i++ {
		sh.SGG[i] = getPresence(dp, code, goal)
	}
	for i := 0; i < 4; i++ {
		sh.RogueWorldBig[i] = getPresence(dp, code, goal)
	}
	for i := 0; i < 10; i++ {
		sh.RogueWorldSmall[i] = getPresence(dp, code, goal)
	}
	return nil
}

func (sh *StellarHex) confirmAbsenseOfTheRest() {
	if sh.StarSystem == "?" {
		sh.StarSystem = "-"
	}
	if sh.BH == "?" {
		sh.BH = "-"
	}
	if sh.NS == "?" {
		sh.NS = "-"
	}
	if sh.D == "?" {
		sh.D = "-"
	}
	if sh.BD == "?" {
		sh.BD = "-"
	}
	if sh.LGG == "?" {
		sh.LGG = "-"
	}
	for i := range sh.MGG {
		if sh.MGG[i] == "?" {
			sh.MGG[i] = "-"
		}
	}
	for i := range sh.SGG {
		if sh.SGG[i] == "?" {
			sh.SGG[i] = "-"
		}
	}
	for i := range sh.RogueWorldBig {
		if sh.RogueWorldBig[i] == "?" {
			sh.RogueWorldBig[i] = "-"
		}
	}
	for i := range sh.RogueWorldSmall {
		if sh.RogueWorldSmall[i] == "?" {
			sh.RogueWorldSmall[i] = "-"
		}
	}
}

func getPresence(dp *dice.Dicepool, code, goal string) string {
	verdict := "?"
	switch dp.Check(code, goal) {
	case true:
		verdict = "+"
	case false:
		verdict = "-"
	}
	return verdict
}

func checkGoal(density int) (string, string) {
	code := ""
	goal := ""
	switch density {
	case Density_Void:
		code = "3d6"
		goal = "3-"
	case Density_Rift:
		code = "2d6"
		goal = "2-"
	case Density_Sparse:
		code = "1d6"
		goal = "1-"
	case Density_Scattered:
		code = "1d6"
		goal = "2-"
	case Density_Standard:
		code = "1d6"
		goal = "3-"
	case Density_Dense:
		code = "1d6"
		goal = "4-"
	case Density_Cluster:
		code = "1d6"
		goal = "5-"
	case Density_Core:
		code = "2d6"
		goal = "11-"
	}
	return code, goal
}
