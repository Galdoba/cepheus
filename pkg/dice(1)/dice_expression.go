package dice

import (
	"fmt"
	"strings"
)

type DiceExpression string

const (
	Additive       ModType = "additive"
	Multiplicative ModType = "multiplicative"
	Deletive       ModType = "deletive"
	SumMininum     ModType = "minimum"
	SumMaximum     ModType = "maximum"
	DropLow        ModType = "drop low"
	DropHigh       ModType = "drop high"
	Individual     ModType = "individual"
)

type ModType string

type SumDirectives struct {
	Num     int
	Faces   int
	SumMods map[ModType]int
	Replace map[int]int
	ReRoll  map[int]bool
}

func newSumDirectives() SumDirectives {
	sd := SumDirectives{}
	sd.SumMods = make(map[ModType]int)
	sd.SumMods[Multiplicative] = 1
	sd.SumMods[Deletive] = 1
	sd.Replace = make(map[int]int)
	sd.ReRoll = make(map[int]bool)
	return sd
}

type ConcatDirectives struct {
	Faces []int
	Mods  []int
}

func (de DiceExpression) ParseRoll() (SumDirectives, error) {
	s := string(de)
	s = strings.TrimSpace(strings.ToLower(s))
	sumD, err := parseSumString(s)
	if err != nil {
		return SumDirectives{}, fmt.Errorf("failed to parse roll expression: %v", err)
	}
	return sumD, nil
}
func (de DiceExpression) ParseConcatRoll() (ConcatDirectives, error) {
	s := string(de)
	s = strings.TrimSpace(strings.ToLower(s))
	conD, err := parseConcatString(s)
	if err != nil {
		return ConcatDirectives{}, fmt.Errorf("failed to parse concat dice expression: %v", err)
	}
	return conD, nil
}
