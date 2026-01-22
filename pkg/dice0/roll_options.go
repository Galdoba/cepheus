package dice

import (
	"math"
)

const (
	FIRST_DICE = 0
	LAST_DICE  = math.MaxInt
)

//OPTIONS

type rollOptions struct {
	treatMap         map[int]int
	numberedDiceMods map[int]int
	upperSumLimit    int
	lowerSumLimit    int
	perDieMod        int
	totalDieMod      int
	firstDiceMod     int
	secondDiceMod    int
}

func defaultRollOptions() rollOptions {
	ro := rollOptions{}
	ro.treatMap = make(map[int]int)
	ro.numberedDiceMods = make(map[int]int)
	ro.upperSumLimit = math.MaxInt
	ro.lowerSumLimit = math.MinInt
	return ro
}

type RollOption func(*rollOptions)

func MaxLimit(max int) RollOption {
	return func(ro *rollOptions) {
		ro.upperSumLimit = max
	}
}

func MinLimit(max int) RollOption {
	return func(ro *rollOptions) {
		ro.lowerSumLimit = max
	}
}

func SingleDiceMod(diceNumber, mod int) RollOption {
	return func(ro *rollOptions) {
		ro.numberedDiceMods[diceNumber] = mod
	}
}

func FirstDiceMod(mod int) RollOption {
	return func(ro *rollOptions) {
		ro.firstDiceMod = mod
	}
}

func SecondDiceMod(mod int) RollOption {
	return func(ro *rollOptions) {
		ro.secondDiceMod = mod
	}
}

func ForEveryDice(i int) RollOption {
	return func(ro *rollOptions) {
		ro.perDieMod = i
	}
}

func TreatAs(result, value int) RollOption {
	return func(ro *rollOptions) {
		ro.treatMap[result] = value
	}
}

func DM(dms ...int) RollOption {
	return func(ro *rollOptions) {
		for _, dm := range dms {
			ro.totalDieMod += dm
		}
	}
}

func DM_conditional(dmMap map[string]int, validConditions ...string) RollOption {
	return func(ro *rollOptions) {
		for _, key := range validConditions {
			if dm, ok := dmMap[key]; ok {
				ro.totalDieMod += dm
			}
		}
	}
}

// utils
func bound(i, min, max int) int {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
