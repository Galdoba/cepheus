package tttable

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
)

// D66Table represents a special type of Table with keys are strings (concateneated rolls of 2d6)
// It has mandatory keys 36 keys "11"-"66"
// Each dice modified separatly and it's value is locked between "0" and "9" (with 100 keys max: "00"-"99")
type D66Table struct {
	Name           string               `json:"name" toml:"name" yaml:"name"`
	DiceExpression string               `json:"dice_expression" toml:"dice_expression" yaml:"dice_expression" ` // Default dice expression for rolls
	Rows           map[string]string    `json:"rows" toml:"rows" yaml:"rows" `                                  // Key-Value pairs for table events
	ModsFirst      map[string]int       `json:"mods,omitempty" toml:"mods,omitempty" yaml:"mods,omitempty"`     // Available modifiers for this table
	ModsSecond     map[string]int       `json:"mods,omitempty" toml:"mods,omitempty" yaml:"mods,omitempty"`     // Available modifiers for this table
	modsToApply    []string             `json:"-" toml:"-" yaml:"-"`                                            // Modifiers to apply automatically
	parsed         map[string]*RangeKey `json:"-" toml:"-" yaml:"-"`                                            // Cache for parsed keys
	path           string               `json:"-" toml:"-" yaml:"-"`                                            // Path to file
}

func (t *D66Table) GetName() string {
	return t.Name
}

func (t *D66Table) Roll(r Roller, mods ...string) (string, string, error) {
	key, value, err := t.roll(r, mods...)
	return key, value, err
}

func (t *D66Table) Find(key string) (string, error) {
	return t.FindByCode(key)
}

func (t *D66Table) GetAll() map[string]string {
	events := make(map[string]string)
	maps.Copy(events, t.Rows)
	return events
}

// D66TableOption is a functional option for configuring a D66Table
type D66TableOption func(*D66Table) error

// WithRow adds a row to the table
func WithRow(row Row) D66TableOption {
	return func(t *D66Table) error {
		return t.AddRow(row.Key, row.Value)
	}
}

// WithRows adds multiple rows to the table
func WithRows(rows ...Row) D66TableOption {
	return func(t *D66Table) error {
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
// Mandatory for D66Table
func WithDiceExpression(dexpr string) D66TableOption {
	return func(t *D66Table) error {
		if t.DiceExpression != "" {
			return fmt.Errorf("duplicated option: dice expression")
		}
		t.DiceExpression = dexpr
		return nil
	}
}

// WithMods sets the modifier values available for this table
// Keys are modifier names, values are integer bonuses/penalties
func WithMods(mods map[string]int) D66TableOption {
	return func(t *D66Table) error {
		if len(t.Mods) != 0 {
			return fmt.Errorf("duplicated option: mods")
		}
		t.Mods = mods
		return nil
	}
}

// WithModsToApply specifies which modifiers should be applied automatically
// during rolls. Modifiers are summed together before applying to dice roll.
func WithModsToApply(mods ...string) D66TableOption {
	return func(t *D66Table) error {
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
func NewD66Table(name string, opts ...D66TableOption) (*D66Table, error) {
	if name == "" {
		return nil, fmt.Errorf("empty name provided")
	}
	t := &D66Table{
		Name:       name,
		Rows:       make(map[string]string),
		ModsFirst:  make(map[string]int),
		ModsSecond: make(map[string]int),
		parsed:     make(map[string]*RangeKey),
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

func (t *D66Table) roll(roller Roller, mods ...string) (string, string, error) {
	if roller == nil {
		return "", "", ErrRollerNotSet
	}

	// Calculate total modifier
	dm1 := 0
	dm2 := 0
	modsUsed := []string{}
	switch len(mods) {
	case 0:
		// Use automatic mods if no specific mods provided
		for _, mod := range t.modsToApply {
			if v, ok := t.ModsFirst[mod]; ok {
				dm1 += v
				modsUsed = append(modsUsed, mod)
			}
		}
		for _, mod := range t.modsToApply {
			if v, ok := t.ModsSecond[mod]; ok {
				dm2 += v
				modsUsed = append(modsUsed, mod)
			}
		}
	default:
		// Use provided mods
		for _, mod := range mods {
			if v, ok := t.ModsFirst[mod]; ok {
				dm1 += v
				modsUsed = append(modsUsed, mod)
			}
		}
		for _, mod := range mods {
			if v, ok := t.ModsSecond[mod]; ok {
				dm2 += v
				modsUsed = append(modsUsed, mod)
			}
		}
	}

	// Build dice expression with modifier
	expr := fmt.Sprintf("2d6cm1:%vcm2:%v", dm1, dm2)

	result, err := roller.ConcatRollSafe(expr)
	if err != nil {
		return "", "", fmt.Errorf("roll failed: %w", err)
	}
	outcome, err := t.FindByCode(result)
	if err != nil {
		return "", "", fmt.Errorf("failed to produce outcome for result %v: %v", result, err)
	}

	return result, outcome, nil
}

// FindByRoll finds an event by roll result
func (t *D66Table) FindByCode(roll string) (string, error) {
	for key, value := range t.Rows {
		if t.matchCode(key, roll) {
			return value, nil
		}
	}
	return "", fmt.Errorf("no event found for roll %d", roll)
}

// Validate validates the table structure and keys
func (t *D66Table) Validate() error {
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
func (t *D66Table) GetKeys() []string {
	keys := make([]string, 0, len(t.Rows))
	for k := range t.Rows {
		keys = append(keys, k)
	}
	return keys
}

// matchKey checks if a roll matches a given key
func (t *D66Table) matchCode(key string, roll string) bool {
	rng, err := t.parseKey(key)
	if err != nil {
		return false
	}
	for _, v := range rng {
		if roll == v {
			return true
		}
	}
	return false
}

func toCode(n int) string {
	if n < 0 || n > 99 {
		return ""
	}
	s := ""
	if n < 10 {
		s += "0"
	}
	s += strconv.Itoa(n)
	return s
}

func validCode(code string) bool {
	n, err := strconv.Atoi(code)
	if err != nil {
		return false
	}
	return code == toCode(n)
}

func (t *D66Table) parseKey(key string) ([]string, error) {
	if validCode(key) {
		return []string{key}, nil
	}
	codes := strings.Split(key, "-")
	if len(codes) != 2 {
		return nil, fmt.Errorf("can't parse '%v'", key)
	}
	if !validCode(codes[0]) || !validCode(codes[1]) {
		return nil, fmt.Errorf("can't parse '%v'", key)
	}
	min, _ := strconv.Atoi(codes[0])
	max, _ := strconv.Atoi(codes[1])

	if min >= max {
		return nil, fmt.Errorf("can't parse '%v': invalid range 'min-max'", key)
	}
	keys := []string{}
	for i := min; i <= max; i++ {
		keys = append(keys, toCode(i))
	}

	return keys, nil
}

// validateKey validates a single key
func (t *D66Table) validateKey(key string) error {
	_, err := ParseKey(key)
	return err
}

// checkOverlaps checks for overlapping ranges in the table
func (t *D66Table) checkOverlaps() error {
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
