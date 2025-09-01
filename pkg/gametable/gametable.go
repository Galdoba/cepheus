package gametable

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Galdoba/cepheus/pkg/dice"
)

const (
	lowResultBorder  = -30
	highResultBorder = 30
	maxDepth         = 10
)

// GameTable represents game rolling table with dice result(s) as key, and strings as result.
type GameTable struct {
	Name     string        `json:"table name"`
	RollCode string        `json:"roll"`
	Options  []*RollResult `json:"options"`
	mods     int
}

// RollResult result of GameTable Roll.
// Index - is coded dice value
// Result - will be returned if value is matched. It will be ignored if NextTable is Provided.
// NextTable - will be automaticly rolled if provided.
type RollResult struct {
	Index     string     `json:"index"`
	Result    string     `json:"result"`
	NextTable *GameTable `json:"next table,omitempty"`
}

// NewTable creates new game table.
func NewTable(name, code string, options ...*RollResult) (*GameTable, error) {
	tab := GameTable{
		Name:     name,
		RollCode: code,
		Options:  options,
	}
	return &tab, assertIndexes(tab.Options)
}

func (tab *GameTable) WithMod(mod int) *GameTable {
	tab.mods = mod
	return tab
}

func assertIndexes(rrs []*RollResult) error {
	indexAssetions := make(map[int]int)
	allIndexes := []int{}
	for _, option := range rrs {
		indexes, err := parseIndex(option.Index)
		if err != nil {
			return fmt.Errorf("index assertion: %v", err)
		}
		allIndexes = append(allIndexes, indexes...)
	}
	slices.Sort(allIndexes)
	if len(allIndexes) == 0 {
		return fmt.Errorf("table is empty")
	}
	for _, index := range allIndexes {
		indexAssetions[index]++
		if index < lowResultBorder || index > highResultBorder {
			return fmt.Errorf("out of bounds index provided: %v", index)
		}
	}
	for i := allIndexes[0]; i <= allIndexes[len(allIndexes)-1]; i++ {
		if indexAssetions[i] != 1 {
			return fmt.Errorf("failed to assert index %v: met %v times", i, indexAssetions[i])
		}
	}
	return nil
}

// NewRollResult created new GameTable Option struct.
func NewRollResult(index string, result string, next *GameTable) *RollResult {
	rr := RollResult{
		Index:     index,
		Result:    result,
		NextTable: next,
	}
	return &rr
}

// /matchRollValue checks if dice roll result matches result
func (rr *RollResult) matchRollValue(rv int) bool {
	valid, _ := parseIndex(rr.Index) //ignore error as it would be caught at assertIndexes function
	for _, v := range valid {
		if v == rv {
			return true
		}
	}
	return false
}

// parseIndex parses index to slice of values to compare dice roll result.
func parseIndex(index string) ([]int, error) {
	if index == "" {
		return nil, fmt.Errorf("index was not provided")
	}
	result := []int{}
	simple, err := strconv.Atoi(index)
	if err == nil {
		return append(result, simple), nil
	}
	if strings.HasSuffix(index, "+") {
		indexTrimmed := strings.TrimSuffix(index, "+")
		accending, err := strconv.Atoi(indexTrimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid accending index provided: %v", index)
		}
		for i := accending; i < highResultBorder; i++ {
			result = append(result, i)
		}
		return result, nil
	}
	if strings.HasSuffix(index, "-") {
		indexTrimmed := strings.TrimSuffix(index, "-")
		decending, err := strconv.Atoi(indexTrimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid decending index provided: %v", index)
		}
		for i := decending; i > lowResultBorder; i-- {
			result = append(result, i)
		}
		return result, nil
	}
	if strings.Contains(index, "..") {
		bounds := strings.Split(index, "..")
		low, err := strconv.Atoi(bounds[0])
		if err != nil {
			return nil, fmt.Errorf("invalid slice index: %v", index)
		}
		high, err := strconv.Atoi(bounds[1])
		if err != nil {
			return nil, fmt.Errorf("invalid slice index: %v", index)
		}
		if low > high {
			return nil, fmt.Errorf("invalid range: %v", index)
		}
		if low < lowResultBorder || high > highResultBorder {
			return nil, fmt.Errorf("index out of bounds: %v", index)
		}
		for i := low; i <= high; i++ {
			result = append(result, i)
		}
		return result, nil
	}
	return nil, fmt.Errorf("invalid index format: %v", index)
}

// Roll is a cascading roll function, that will return last result from chain of tables.
func (gt *GameTable) Roll(dp *dice.Dicepool) (string, error) {
	return gt.roll(dp, 0)
}

func (gt *GameTable) roll(dp *dice.Dicepool, depth int) (string, error) {
	if depth > maxDepth {
		return "", fmt.Errorf("maximum depth reached: %v", depth)
	}
	roll := dp.Sum(gt.RollCode)
	for _, opt := range gt.Options {
		if !opt.matchRollValue(roll) {
			continue
		}
		if opt.NextTable != nil {
			return opt.NextTable.roll(dp, depth+1)
		}
		gt.mods = 0
		return opt.Result, nil
	}

	return "", fmt.Errorf("table %v: index %v not found", gt.Name, roll)
}
