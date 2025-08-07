package dice

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

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
