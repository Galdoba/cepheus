package tttable

import (
	"fmt"
	"maps"
	"strconv"
)

// Constants for bounds
const (
	MinRollBound         = -10000
	MaxRollBound         = 10000
	AndLess              = -10001
	AndMore              = 10001
	Flat         ModType = "flat"
	Cumulative   ModType = "cumulative"
	Max          ModType = "max"
	Min          ModType = "min"
)

type ModType string

// Predefined errors
var (
	ErrNoArguments            = fmt.Errorf("no arguments provided")
	ErrMinMaxEqual            = fmt.Errorf("minimum and maximum are equal")
	ErrBothBoundsExceeded     = fmt.Errorf("range exceeds both lower and upper bounds")
	ErrOutOfBounds            = fmt.Errorf("entire range is out of bounds")
	ErrTableNotFound          = fmt.Errorf("table not found")
	ErrRollerNotSet           = fmt.Errorf("roller not set")
	ErrRollerExpressionNotSet = fmt.Errorf("roller expression not set")
	ErrCascadeTooDeep         = fmt.Errorf("cascade too deep (possible loop)")
)

// Table represents a random event table
type Table struct {
	Name           string                `json:"name" toml:"name" yaml:"name"`
	DiceExpression string                `json:"dice_expression" toml:"dice_expression" yaml:"dice_expression" ` // Default dice expression for rolls
	Rows           map[string]TableEntry `json:"rows" toml:"rows" yaml:"rows" `                                  // Key-Value pairs for table events
	ModsFlat       map[string]int        `json:"mods_flat,omitempty" toml:"mods_flat,omitempty" yaml:"mods_flat,omitempty"`
	ModsCumulative map[string]int        `json:"mods_cumulative,omitempty" toml:"mods_cumulative,omitempty" yaml:"mods_cumulative,omitempty"`
	ModsMax        map[string]int        `json:"mods_max,omitempty" toml:"mods_max,omitempty" yaml:"mods_max,omitempty"`
	ModsMin        map[string]int        `json:"mods_min,omitempty" toml:"mods_min,omitempty" yaml:"mods_min,omitempty"`
	modsToApply    []string              `json:"-" toml:"-" yaml:"-"` // Modifiers to apply automatically
	parsed         map[string]*RangeKey  `json:"-" toml:"-" yaml:"-"` // Cache for parsed keys
	path           string                `json:"-" toml:"-" yaml:"-"` // Path to file
}

func (t *Table) GetName() string {
	return t.Name
}

func (t *Table) Roll(r Roller, mods ...string) (string, string, error) {
	index, value, err := t.roll(r, mods...)
	key := strconv.Itoa(index)
	return key, value, err
}

func (t *Table) Find(key string) (string, error) {
	n, err := strconv.Atoi(key)
	if err != nil {
		return "", fmt.Errorf("key is not a number: %v", key)
	}
	return t.FindByRoll(n)
}

func (t *Table) GetAll() map[string]TableEntry {
	events := make(map[string]TableEntry)
	maps.Copy(events, t.Rows)
	return events
}

// TableOption is a functional option for configuring a Table
type TableOption func(*Table) error

// WithRow adds a row to the table
// func WithIndexEntry(row TableEntry) TableOption {
// 	return func(t *Table) error {
// 		return t.AddRow(row.Key, row)
// 	}
// }

// WithRows adds multiple rows to the table
func WithIndexEntries(rows ...TableEntry) TableOption {
	return func(t *Table) error {
		if len(t.Rows) > 0 {
			return fmt.Errorf("duplicated option: entries")
		}
		for _, row := range rows {
			if _, present := t.Rows[row.Key]; present {
				return fmt.Errorf("duplicated entriy key provided: %v", row.Key)
			} else {
				t.Rows[row.Key] = row
			}
		}
		return nil
	}
}

// WithDiceExpression sets the dice expression for table rolls
// Example: "2d6", "1d20+5", "3d10-2"
// Mandatory for Table
func WithDiceExpression(dexpr string) TableOption {
	return func(t *Table) error {
		if t.DiceExpression != "" {
			return fmt.Errorf("duplicated option: dice expression")
		}
		t.DiceExpression = dexpr
		return nil
	}
}

// WithMods sets the modifier values available for this table
// Keys are modifier names, values are integer bonuses/penalties
func WithIndexMods(mtype ModType, mods map[string]int) TableOption {
	return func(t *Table) error {
		err := fmt.Errorf("duplicated option")
		switch mtype {
		case Flat:
			if t.ModsFlat != nil {
				return fmt.Errorf("%v: flat mods", err)
			}
			t.ModsFlat = make(map[string]int)
			t.ModsFlat = mods
		case Cumulative:
			if t.ModsFlat != nil {
				return fmt.Errorf("%v: cumulative mods", err)
			}
			t.ModsCumulative = make(map[string]int)
			t.ModsCumulative = mods
		case Max:
			if t.ModsFlat != nil {
				return fmt.Errorf("%v: max mods", err)
			}
			t.ModsMax = make(map[string]int)
			t.ModsMax = mods
		case Min:
			if t.ModsFlat != nil {
				return fmt.Errorf("%v: min mods", err)
			}
			t.ModsMin = make(map[string]int)
			t.ModsMin = mods
		}
		return nil
	}
}

// WithModsToApply specifies which modifiers should be applied automatically
// during rolls. Modifiers are summed together before applying to dice roll.
func WithIndexModsToApply(mods ...string) TableOption {
	return func(t *Table) error {
		if len(mods) == 0 {
			return fmt.Errorf("no auto mods provided")
		}
		if len(t.modsToApply) != 0 {
			return fmt.Errorf("duplicated option: mods")
		}
		t.modsToApply = mods
		return nil
	}
}

// NewTable creates a new table with given name and options
func NewTable(name string, opts ...TableOption) (*Table, error) {
	if name == "" {
		return nil, fmt.Errorf("empty name provided")
	}
	t := &Table{
		Name:   name,
		Rows:   make(map[string]TableEntry),
		parsed: make(map[string]*RangeKey),
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}

	// Validate the table
	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("table validation failed: %w", err)
	}

	return t, nil
}

// RemoveRow removes a row from the table
func (t *Table) RemoveRow(key string) {
	delete(t.Rows, key)
	delete(t.parsed, key)
}

// Roll performs a dice roll and returns the corresponding result
// If mods are provided they will substitute own table mods for this roll
func (t *Table) roll(roller Roller, mods ...string) (int, string, error) {
	if roller == nil {
		return AndLess, "", ErrRollerNotSet
	}
	if t.DiceExpression == "" {
		return AndLess, "", ErrRollerExpressionNotSet
	}

	// Calculate total modifier
	dm := 0
	switch len(mods) {
	case 0:
		// Use automatic mods if no specific mods provided
		dm = t.combineMods(t.modsToApply...)
	default:
		// Use provided mods
		dm = t.combineMods(mods...)
	}

	// Build dice expression with modifier
	expr := t.DiceExpression
	if dm < 0 {
		expr += fmt.Sprintf("%v", dm)
	} else if dm > 0 {
		expr += fmt.Sprintf("+%v", dm)
	}
	// If dm == 0, expression remains unchanged

	result, err := roller.RollSafe(expr)
	if err != nil {
		return AndLess, "", fmt.Errorf("roll failed: %w", err)
	}
	outcome, err := t.FindByRoll(result)
	if err != nil {
		return AndLess, "", fmt.Errorf("failed to produce outcome for result %v: %v", result, err)
	}

	return result, outcome, nil
}

func (t *Table) combineMods(mods ...string) int {
	flat := make(map[string]bool)
	cumulative := []string{}
	maxs := []string{}
	mins := []string{}
	for _, mod := range mods {
		if _, ok := t.ModsFlat[mod]; ok {
			flat[mod] = true
		}
		if _, ok := t.ModsCumulative[mod]; ok {
			cumulative = append(cumulative, mod)
		}
		if _, ok := t.ModsMin[mod]; ok {
			maxs = append(maxs, mod)
		}
		if _, ok := t.ModsMax[mod]; ok {
			mins = append(mins, mod)
		}
	}
	dm := 0
	for mod := range flat {
		if val, ok := t.ModsFlat[mod]; ok {
			dm += val
		}
	}
	for _, mod := range cumulative {
		if val, ok := t.ModsCumulative[mod]; ok {
			dm += val
		}
	}
	maxSlice := []int{}
	for _, mod := range maxs {
		if val, ok := t.ModsMax[mod]; ok {
			maxSlice = append(maxSlice, val)
		}
	}
	maxVal, detectedMax := maxOf(maxSlice)
	if detectedMax {
		dm += maxVal
	}

	minSlice := []int{}
	for _, mod := range mins {
		if val, ok := t.ModsMin[mod]; ok {
			minSlice = append(minSlice, val)
		}
	}
	minVal, detectedMin := minOf(minSlice)
	if detectedMin {
		dm += minVal
	}
	return dm

}

func minOf(slice []int) (int, bool) {
	if len(slice) < 1 {
		return 0, false
	}
	val := slice[0]
	for _, next := range slice {
		val = min(val, next)
	}
	return val, true
}

func maxOf(slice []int) (int, bool) {
	if len(slice) < 1 {
		return 0, false
	}
	val := slice[0]
	for _, next := range slice {
		val = max(val, next)
	}
	return val, true
}

// FindByRoll finds an event by roll result
func (t *Table) FindByRoll(roll int) (string, error) {
	for key, value := range t.Rows {
		if t.matchKey(key, roll) {
			return value.Value, nil
		}
	}
	return "", fmt.Errorf("no event found for roll %d", roll)
}

// Validate validates the table structure and keys
func (t *Table) Validate() error {
	// Check if table has any rows
	if len(t.Rows) == 0 {
		return fmt.Errorf("table has no rows")
	}

	// Validate each key format
	for key := range t.Rows {
		if err := t.validateKey(key); err != nil {
			return fmt.Errorf("invalid key %q: %w", key, err)
		}
	}

	// Check for range overlaps
	return t.checkOverlaps()
}

// GetKeys returns all keys in the table
func (t *Table) GetKeys() []string {
	keys := make([]string, 0, len(t.Rows))
	for k := range t.Rows {
		keys = append(keys, k)
	}
	return keys
}

// matchKey checks if a roll matches a given key
func (t *Table) matchKey(key string, roll int) bool {
	rng, err := t.parseKey(key)
	if err != nil {
		return false
	}

	// Check minimum boundary
	minCheck := false
	if rng.MinInclusive {
		minCheck = roll >= rng.Min
	} else {
		minCheck = roll > rng.Min
	}

	// Check maximum boundary
	maxCheck := false
	if rng.MaxInclusive {
		maxCheck = roll <= rng.Max
	} else {
		maxCheck = roll < rng.Max
	}

	return minCheck && maxCheck
}

// parseKey parses a key and caches the result
func (t *Table) parseKey(key string) (*RangeKey, error) {
	// Check cache first
	if cached, ok := t.parsed[key]; ok {
		return cached, nil
	}

	rng, err := ParseKey(key)
	if err != nil {
		return nil, err
	}

	// Validate bounds
	if rng.Min < MinRollBound {
		return nil, fmt.Errorf("minimum value %d is below allowed bound %d", rng.Min, MinRollBound)
	}
	if rng.Max > MaxRollBound {
		return nil, fmt.Errorf("maximum value %d is above allowed bound %d", rng.Max, MaxRollBound)
	}

	t.parsed[key] = rng
	return rng, nil
}

// validateKey validates a single key
func (t *Table) validateKey(key string) error {
	_, err := ParseKey(key)
	return err
}

// checkOverlaps checks for overlapping ranges in the table
func (t *Table) checkOverlaps() error {
	var ranges []*RangeKey

	// Parse all keys
	for key := range t.Rows {
		rng, err := ParseKey(key)
		if err != nil {
			return err
		}
		ranges = append(ranges, rng)
	}

	// Check for overlaps
	for i, r1 := range ranges {
		for j, r2 := range ranges {
			if i >= j {
				continue
			}

			if rangesOverlap(r1, r2) {
				return fmt.Errorf("overlapping ranges detected: %s and %s", r1.Original, r2.Original)
			}
		}
	}

	return nil
}
