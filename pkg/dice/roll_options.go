package dice

import "fmt"

func (dp *Dicepool) Sum(code string, options ...RollOption) int {
	sum := 0
	maxEdge := 0
	locked := false
	if dp.locked || code == "" {
		locked = true
	}
	switch locked {
	case true:
	case false:
		rr := parseDiceCode(code)
		dp.dices = make(map[int]dice)
		for die := rr.diceNum; die > 0; die-- {
			dp.dices[die] = dice{edges: rr.diceEdges, result: dp.roller.rand.Intn(rr.diceEdges) + 1}
			if rr.diceEdges > maxEdge {
				maxEdge = rr.diceEdges
			}
		}
	}
	rollOpts := defaultRollOptions()
	for _, modify := range options {
		modify(&rollOpts)
	}
	for _, d := range dp.dices {
		if val, ok := rollOpts.treatMap[d.result]; ok {
			d.result = val
		}
		sum += d.result + rollOpts.perDieMod
	}
	sum += rollOpts.totalDieMod
	sum = bound(sum, rollOpts.lowerSumLimit, rollOpts.upperSumLimit)
	return sum
}

func (dp *Dicepool) D66() string {
	return fmt.Sprintf("%v%v", dp.Sum("d"), dp.Sum("D"))
}

func (dp *Dicepool) Flux() int {
	return dp.Sum("1d6") - dp.Sum("1d6")
}

func (dp *Dicepool) FluxGood() int {
	d1 := dp.Sum("1d6")
	d2 := dp.Sum("1d6")
	return max(d1, d2) - min(d1, d2)
}

func (dp *Dicepool) FluxBad() int {
	d1 := dp.Sum("1d6")
	d2 := dp.Sum("1d6")
	return min(d1, d2) - max(d1, d2)
}

//OPTIONS

type rollOptions struct {
	treatMap      map[int]int
	upperSumLimit int
	lowerSumLimit int
	perDieMod     int
	totalDieMod   int
}

func defaultRollOptions() rollOptions {
	ro := rollOptions{}
	ro.treatMap = make(map[int]int)
	ro.upperSumLimit = 1000
	ro.lowerSumLimit = -1000
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

//utils
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
