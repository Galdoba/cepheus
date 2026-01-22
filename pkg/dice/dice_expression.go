package dice

import (
	"fmt"
	"strings"
)

// DiceExpression is a string type that represents a dice expression.
// It provides methods to parse the expression into executable directives.
type DiceExpression string

// These constants represent different categories of mathematical operations.
const (
	Additive       ModType = "additive"       // Adds/subtracts from the final sum
	Multiplicative ModType = "multiplicative" // Multiplies the final sum
	Deletive       ModType = "deletive"       // Divides the final sum (integer division)
	SumMininum     ModType = "minimum"        // Sets a minimum bound for the final sum
	SumMaximum     ModType = "maximum"        // Sets a maximum bound for the final sum
	DropLow        ModType = "drop low"       // Drops the lowest n dice from the pool
	DropHigh       ModType = "drop high"      // Drops the highest n dice from the pool
	Individual     ModType = "individual"     // Modifies each individual die value
)

// ModType defines the type of modifier that can be applied to dice rolls.
type ModType string

// SumDirectives contains all the parsed instructions for a sum-based dice roll.
// It represents the complete set of operations to perform on a dice pool.
type SumDirectives struct {
	Num     int             // Number of dice to roll
	Faces   int             // Number of faces on each die
	SumMods map[ModType]int // Mathematical modifiers to apply to the sum
	Replace map[int]int     // Value replacement mapping (original -> replacement)
	ReRoll  map[int]bool    // Values that should be rerolled
}

// newSumDirectives creates and initializes a new SumDirectives struct with default values.
// The default multiplicative and deletive factors are set to 1 (no effect).
func newSumDirectives() SumDirectives {
	sd := SumDirectives{}
	sd.SumMods = make(map[ModType]int)
	sd.SumMods[Multiplicative] = 1 // Default: multiply by 1 (no change)
	sd.SumMods[Deletive] = 1       // Default: divide by 1 (no change)
	sd.Replace = make(map[int]int)
	sd.ReRoll = make(map[int]bool)
	return sd
}

// ConcatDirectives contains all the parsed instructions for a concatenated dice roll.
// Each die in the pool contributes a single digit to the final string result.
type ConcatDirectives struct {
	Faces []int // Number of faces for each die in the concatenation
	Mods  []int // Individual modifiers for each die position
}

// ParseRoll parses a dice expression string into SumDirectives for a sum-based roll.
// The expression is case-insensitive and trimmed of whitespace before parsing.
func (de DiceExpression) ParseRoll() (SumDirectives, error) {
	s := string(de)
	s = strings.TrimSpace(strings.ToLower(s))
	sumD, err := parseSumString(s)
	if err != nil {
		return SumDirectives{}, fmt.Errorf("failed to parse roll expression: %v", err)
	}
	return sumD, nil
}

// ParseConcatRoll parses a dice expression string into ConcatDirectives for a concatenated roll.
// The expression is case-insensitive and trimmed of whitespace before parsing.
func (de DiceExpression) ParseConcatRoll() (ConcatDirectives, error) {
	s := string(de)
	s = strings.TrimSpace(strings.ToLower(s))
	conD, err := parseConcatString(s)
	if err != nil {
		return ConcatDirectives{}, fmt.Errorf("failed to parse concat dice expression: %v", err)
	}
	return conD, nil
}
