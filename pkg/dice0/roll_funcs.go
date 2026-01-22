package dice

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

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
	lastDice := len(dp.dices) - 1
	rollOpts.numberedDiceMods[lastDice] = rollOpts.numberedDiceMods[LAST_DICE]
	for i, d := range dp.dices {
		if val, ok := rollOpts.treatMap[d.result]; ok {
			d.result = val
		}
		sum += d.result + rollOpts.perDieMod + rollOpts.numberedDiceMods[i]
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

func (dp *Dicepool) Check(code, goal string, rollOpts ...RollOption) bool {
	sum := dp.Sum(code, rollOpts...)
	goalValues := parseGoal(goal)
	return slices.Contains(goalValues, sum)
}

func parseGoal(goal string) []int {
	output := []int{}
	parts := strings.Split(goal, ",")
	for _, part := range parts {
		//pure tn
		val, err := strconv.Atoi(part)
		if err == nil {
			output = append(output, val)
			continue
		}
		//tn or more
		if strings.HasSuffix(part, "+") {
			partSegment := strings.TrimSuffix(part, "+")
			val, err := strconv.Atoi(partSegment)
			if err == nil {
				for i := val; i <= 30; i++ {
					output = append(output, i)
				}
				continue
			}
		}
		//tn or less
		if strings.HasSuffix(part, "-") {
			partSegment := strings.TrimSuffix(part, "-")
			val, err := strconv.Atoi(partSegment)
			if err == nil {
				for i := val; i >= -30; i-- {
					output = append(output, i)
				}
				continue
			}
		}
		//tn from x to y
		segments := strings.Split(part, "...")
		if len(segments) == 2 {
			x, errX := strconv.Atoi(segments[0])
			y, errY := strconv.Atoi(segments[1])
			if errX == nil && errY == nil {
				for i := x; i <= y; i++ {
					output = append(output, i)
				}
				continue
			}
		}
		fmt.Printf("failed to parse goal '%v'\n", goal)

	}

	return uniqueInts(output)
}

func uniqueInts(sl []int) []int {
	intMap := make(map[int]int)
	for _, i := range sl {
		intMap[i]++
	}
	sl = []int{}
	for k := range intMap {
		sl = append(sl, k)
	}
	slices.Sort(sl)
	return sl
}

func (dp *Dicepool) Sum1D(options ...RollOption) int {
	return dp.Sum("1d6", options...)
}

func RandomIndex[T any](objects []T) int {
	max := fmt.Sprintf("%v", len(objects))
	return NewDicepool().Sum(fmt.Sprintf("1d%v", max)) - 1
}

func (dp *Dicepool) SkillCheck(mods ...int) int {
	r := dp.Sum("2d6")
	for _, mod := range mods {
		r += mod
	}
	return r
}

func FastRandom(code string) int {
	return NewDicepool().Sum(code)
}

func FromSliceRandom[T any](sl []T) int {
	if len(sl) == 0 {
		panic("nothing to choose from")
	}
	return NewDicepool().Sum(fmt.Sprintf("1d%v", len(sl))) - 1
}

func (dp *Dicepool) Roll(code string) int {
	return dp.Sum(code)
}
