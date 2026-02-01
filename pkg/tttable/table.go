package tttable

import (
	"fmt"
	"maps"
	"strconv"
)

// Constants for bounds
const (
	MinRollBound = -10000
	MaxRollBound = 10000
	AndLess      = -10001
	AndMore      = 10001
)

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
	Name           string               `json:"name" toml:"name" yaml:"name"`
	DiceExpression string               `json:"dice_expression" toml:"dice_expression" yaml:"dice_expression" ` // Default dice expression for rolls
	Rows           map[string]string    `json:"rows" toml:"rows" yaml:"rows" `                                  // Key-Value pairs for table events
	Mods           map[string]int       `json:"mods,omitempty" toml:"mods,omitempty" yaml:"mods,omitempty"`     // Available modifiers for this table
	modsToApply    []string             `json:"-" toml:"-" yaml:"-"`                                            // Modifiers to apply automatically
	parsed         map[string]*RangeKey `json:"-" toml:"-" yaml:"-"`                                            // Cache for parsed keys
	path           string               `json:"-" toml:"-" yaml:"-"`                                            // Path to file
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

func (t *Table) GetAll() map[string]string {
	events := make(map[string]string)
	maps.Copy(events, t.Rows)
	return events
}

// TableOption is a functional option for configuring a Table
type TableOption func(*Table) error

// WithRow adds a row to the table
func WithRow(row Row) TableOption {
	return func(t *Table) error {
		return t.AddRow(row.Key, row.Value)
	}
}

// WithRows adds multiple rows to the table
func WithRows(rows ...Row) TableOption {
	return func(t *Table) error {
		for _, row := range rows {
			if err := t.AddRow(row.Key, row.Value); err != nil {
				return err
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
func WithMods(mods map[string]int) TableOption {
	return func(t *Table) error {
		if len(t.Mods) != 0 {
			return fmt.Errorf("duplicated option: mods")
		}
		t.Mods = mods
		return nil
	}
}

// WithModsToApply specifies which modifiers should be applied automatically
// during rolls. Modifiers are summed together before applying to dice roll.
func WithModsToApply(mods ...string) TableOption {
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
		Rows:   make(map[string]string),
		Mods:   make(map[string]int),
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

// AddRow adds a new row to the table
func (t *Table) AddRow(key, value string) error {
	if err := t.validateKey(key); err != nil {
		return err
	}

	// Check for duplicates
	if _, exists := t.Rows[key]; exists {
		return fmt.Errorf("duplicate key: %s", key)
	}

	t.Rows[key] = value
	delete(t.parsed, key) // Clear cache for this key
	return nil
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
	modsUsed := []string{}
	switch len(mods) {
	case 0:
		// Use automatic mods if no specific mods provided
		for _, mod := range t.modsToApply {
			if v, ok := t.Mods[mod]; ok {
				dm += v
				modsUsed = append(modsUsed, mod)
			}
		}
	default:
		// Use provided mods
		for _, mod := range mods {
			if v, ok := t.Mods[mod]; ok {
				dm += v
				modsUsed = append(modsUsed, mod)
			}
		}
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

// FindByRoll finds an event by roll result
func (t *Table) FindByRoll(roll int) (string, error) {
	for key, value := range t.Rows {
		if t.matchKey(key, roll) {
			return value, nil
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
