package dice

/*
Package dice provides a comprehensive dice rolling system with two main rolling methods:
- Roll: Rolls dice pool, applies modifiers, and returns the sum of values as an integer.
- ConcatRoll: Rolls dice pool, applies modifiers, and returns concatenated values as a string (code).
  This method only returns values consisting of digits 0-9.
*/

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Roller represents a dice roller with its own random number generator and state.
// It maintains state for the last roll results, allowing for inspection of rolls.
type Roller struct {
	seed            string     // Seed used to initialize the random number generator
	rng             *rand.Rand // Random number generator instance
	lastResult      []int      // Last dice roll results before any modifications
	lastModifiedSum int        // Last sum after applying all modifiers
	lastConcatRoll  string     // Last concatenated roll result
}

// New creates and returns a new Roller instance with an optional seed.
// If no seed is provided, a time-based random seed is used.
func New(seed string) *Roller {
	return newRoller(seed)
}

// newRoller initializes a Roller with a specific seed.
// It converts the seed string to an int64 or uses a random seed if empty.
func newRoller(seed string) *Roller {
	seedInt := randomSeed()
	if seed != "" {
		seedInt = stringToInt64(seed)
	}
	r := Roller{}
	r.rng = rand.New(rand.NewSource(int64(seedInt)))
	return &r
}

// stringToInt64 converts a string seed to an int64 value using a custom hash function.
// This ensures consistent random sequences from the same seed string.
func stringToInt64(seed string) int64 {
	if seed == "" {
		return 0
	}

	var h1, h2 int64 = 0, 0

	for i := 0; i < len(seed); i++ {
		c := int64(seed[i])

		// First polynomial with a large prime factor
		h1 = h1*131 + c + int64(i)*31

		// Second polynomial with a different factor for better mixing
		h2 = h2*257 + c - int64(i)*17

		// Mix the hashes together
		h1 ^= h2
		h2 ^= h1
	}

	// Combine both hashes
	return h1 ^ (h2 << 32) ^ int64(len(seed))
}

// randomSeed generates a time-based random seed using the current time in nanoseconds.
func randomSeed() int64 {
	return time.Now().UnixNano()
}

// baseRoll performs the basic dice roll without any modifiers.
// It rolls 'number' dice with 'faces' sides each and stores the results.
func (r *Roller) baseRoll(number, faces int) []int {
	r.lastResult = []int{}
	for i := 0; i < number; i++ {
		v := r.rng.Intn(faces) + 1 // Generate random value between 1 and faces
		r.lastResult = append(r.lastResult, v)
	}
	return r.lastResult
}

// rollConcat rolls a single die for concatenation purposes.
// Returns a value between 1 and faces (inclusive).
func (r *Roller) rollConcat(faces int) int {
	if faces == 0 {
		return 0
	}
	return r.rng.Intn(faces) + 1
}

// rollSafe performs a dice roll with comprehensive error handling.
// It applies all directives from SumDirectives in the correct order.
func (r *Roller) rollSafe(sd SumDirectives) (int, error) {
	// Step 1: Roll the base dice
	res := r.baseRoll(sd.Num, sd.Faces)

	// Step 2: Apply modifiers to each individual die roll
	modifiedRolls := make([]int, len(res))
	for i, oldVal := range res {
		currentVal := oldVal

		// 2a. Apply reroll directives
		if len(sd.ReRoll) > 0 {
			reroledVal, err := reroll(r, sd.Faces, currentVal, sd.ReRoll)
			if err != nil {
				return 0, err
			}
			currentVal = reroledVal
		}

		// 2b. Apply value replacements
		if newVal, ok := sd.Replace[currentVal]; ok {
			currentVal = newVal
		}

		// 2c. Apply individual modifiers (additive per die)
		if ind, ok := sd.SumMods[Individual]; ok {
			currentVal = currentVal + ind
		}

		modifiedRolls[i] = currentVal
	}

	// Step 3: Apply drop low and drop high directives
	// Create a copy for sorting to preserve original order
	sortedRolls := make([]int, len(modifiedRolls))
	copy(sortedRolls, modifiedRolls)

	// Drop low values (remove lowest n dice)
	if dl, ok := sd.SumMods[DropLow]; ok && dl > 0 {
		// Sort in ascending order
		slices.Sort(sortedRolls)
		if dl >= len(sortedRolls) {
			return 0, fmt.Errorf("all dices dropped")
		}
		// Drop the lowest 'dl' values
		sortedRolls = sortedRolls[dl:]
	}

	// Drop high values (remove highest n dice)
	if dh, ok := sd.SumMods[DropHigh]; ok && dh > 0 {
		if dh >= len(sortedRolls) {
			return 0, fmt.Errorf("all dices dropped")
		}
		// Drop the highest 'dh' values
		sortedRolls = sortedRolls[:len(sortedRolls)-dh]
	}

	// Step 4: Sum the remaining values
	sum := 0
	r.lastResult = sortedRolls
	for _, v := range sortedRolls {
		sum += v
	}

	// Step 5: Apply modifiers to the final sum
	if mod, ok := sd.SumMods[Additive]; ok {
		sum += mod
	}
	if mod, ok := sd.SumMods[Multiplicative]; ok {
		sum = sum * mod
	}
	if mod, ok := sd.SumMods[Deletive]; ok {
		if mod != 0 {
			sum = sum / mod
		}
	}
	if mod, ok := sd.SumMods[SumMininum]; ok {
		sum = max(mod, sum)
	}
	if mod, ok := sd.SumMods[SumMaximum]; ok {
		sum = min(mod, sum)
	}

	r.lastModifiedSum = sum
	return sum, nil
}

// concat performs a concatenated dice roll where each die result becomes a digit in a string.
// This is commonly used for percentile rolls or other multi-digit codes.
func (r *Roller) concat(cd ConcatDirectives) string {
	result := ""
	for i, faces := range cd.Faces {
		// Roll a single die
		rollValue := r.rollConcat(faces)

		// Apply individual modifier for this die position
		if i < len(cd.Mods) {
			rollValue += cd.Mods[i]
		}

		// Clamp the value to 0-9 range (single digit)
		rollValue = max(0, rollValue)
		rollValue = min(9, rollValue)

		// Append to result string
		result += strconv.Itoa(rollValue)
	}
	r.lastConcatRoll = result

	return result
}

// reroll recursively rerolls a die until it gets a value not in the excluded set.
// Prevents infinite loops by limiting to 1000 attempts.
func reroll(r *Roller, faces, val int, excluded map[int]bool) (int, error) {
	count := 0
	done := false
	current := val

	for !done {
		// Roll a single die
		current = r.baseRoll(1, faces)[0]

		// Check if we need to reroll
		needReroll := false
		for ex := range excluded {
			if current == ex {
				needReroll = true
				count++
				// Safety check: prevent infinite loops
				if count > 1000 {
					return 0, fmt.Errorf("impossible roll 1d%v to exclude %v", faces, excluded)
				}
				break
			}
		}

		if !needReroll {
			done = true
		}
	}

	return current, nil
}

// RollSafe is a convenience function that parses a dice expression and rolls it safely.
// Returns an error if parsing fails or if the roll encounters issues.
func RollSafe(expression string, seed ...string) (int, error) {
	var seedStr string
	if len(seed) > 0 {
		seedStr = seed[0]
	}

	rd, err := DiceExpression(expression).ParseRoll()
	if err != nil {
		return 0, err
	}

	r := newRoller(seedStr)
	return r.rollSafe(rd)
}

// ConcatRollSafe is a convenience function that parses a concatenated dice expression
// and rolls it safely. Returns an error if parsing fails.
func ConcatRollSafe(expression string, seed ...string) (string, error) {
	var seedStr string
	if len(seed) > 0 {
		seedStr = seed[0]
	}

	rd, err := DiceExpression(expression).ParseConcatRoll()
	if err != nil {
		return "", err
	}

	r := newRoller(seedStr)
	return r.concat(rd), nil
}

// Roll is a convenience function that parses a dice expression and rolls it.
// Panics if there's an error during parsing or rolling.
func Roll(expression string, seed ...string) int {
	result, err := RollSafe(expression, seed...)
	if err != nil {
		panic(err)
	}
	return result
}

// ConcatRoll is a convenience function that parses a concatenated dice expression
// and rolls it. Panics if there's an error during parsing or rolling.
func ConcatRoll(expression string, seed ...string) string {
	result, err := ConcatRollSafe(expression, seed...)
	if err != nil {
		panic(err)
	}
	return result
}

// Roll uses the Roller instance to parse and roll a dice expression.
// Panics if there's an error during parsing or rolling.
func (r *Roller) Roll(expression string) int {
	rd, err := DiceExpression(expression).ParseRoll()
	if err != nil {
		panic(err)
	}
	s, err := r.rollSafe(rd)
	if err != nil {
		panic(err)
	}
	return s
}

// ConcatRoll uses the Roller instance to parse and roll a concatenated dice expression.
// Panics if there's an error during parsing.
func (r *Roller) ConcatRoll(expression string) string {
	rd, err := DiceExpression(expression).ParseConcatRoll()
	if err != nil {
		panic(err)
	}
	return r.concat(rd)
}

// LastRoll returns the last modified sum from the most recent roll.
func (r *Roller) LastRoll() int {
	return r.lastModifiedSum
}

// LastConcatRoll returns the last concatenated roll result.
func (r *Roller) LastConcatRoll() string {
	return r.lastConcatRoll
}

// Result returns the last dice roll results before any modifications.
func (r *Roller) Result() []int {
	return r.lastResult
}

// Special Rolls
func Flux() int {
	r := newRoller("")
	r1 := r.Roll("1d6")
	r2 := r.Roll("1d6")
	return r1 - r2
}

func (r *Roller) Flux() int {
	r1 := r.Roll("1d6")
	r2 := r.Roll("1d6")
	return r1 - r2
}

func D66(mods ...int) string {
	r := newRoller("")
	return r.D66(mods...)
}

func (r *Roller) D66(mods ...int) string {
	nodes := []string{}
	for i, mod := range mods {
		nodes = append(nodes, fmt.Sprintf("cm%v:%v", i+1, mod))
	}
	return r.ConcatRoll("d66" + strings.Join(nodes, ""))
}
