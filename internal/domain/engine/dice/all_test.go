package dice

import (
	"reflect"
	"testing"
)

// ---------------------------------------------------------------------
// Helpers for deterministic Roller
// ---------------------------------------------------------------------

func newTestRoller(seed string) *Roller {
	return newRoller(seed)
}

// ---------------------------------------------------------------------
// Tests for builders.go
// ---------------------------------------------------------------------

func TestStringToInt64(t *testing.T) {
	tests := []struct {
		name string
		seed string
		want int64 // we don't check exact value, just non-zero for non-empty and zero for empty
	}{
		{"empty", "", 0},
		{"simple", "test", 0}, // we'll check consistency later
		{"same", "test", 0},
		{"different", "abc", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringToInt64(tt.seed)
			if tt.seed == "" && got != 0 {
				t.Errorf("stringToInt64(\"\") = %d, want 0", got)
			}
			if tt.seed != "" && got == 0 {
				t.Errorf("stringToInt64(%q) = 0, want non-zero", tt.seed)
			}
		})
	}
	// Consistency check
	seed := "consistency"
	first := stringToInt64(seed)
	second := stringToInt64(seed)
	if first != second {
		t.Errorf("stringToInt64(%q) not consistent: %d vs %d", seed, first, second)
	}
}

func TestRandomSeed(t *testing.T) {
	// Just ensure it returns something (non-zero almost always)
	seed := randomSeed()
	if seed == 0 {
		t.Errorf("randomSeed() returned 0, unlikely")
	}
}

func TestNewDice(t *testing.T) {
	d := NewDice(6)
	if d.Faces != 6 {
		t.Errorf("NewDice(6).Faces = %d, want 6", d.Faces)
	}
	if d.Codes != nil {
		t.Errorf("NewDice(6).Codes should be nil, got %v", d.Codes)
	}
	if d.Metadata != nil {
		t.Errorf("NewDice(6).Metadata should be nil, got %v", d.Metadata)
	}
}

func TestDieWithCodes(t *testing.T) {
	d := NewDice(20)
	codes := map[int]string{20: "crit", 1: "fumble"}
	d2 := d.WithCodes(codes)
	if d2.Faces != 20 {
		t.Errorf("Faces = %d, want 20", d2.Faces)
	}
	if !reflect.DeepEqual(d2.Codes, codes) {
		t.Errorf("Codes = %v, want %v", d2.Codes, codes)
	}
	// Original should be unchanged
	if d.Codes != nil {
		t.Errorf("Original Die was mutated")
	}
}

func TestDieWithMeta(t *testing.T) {
	d := NewDice(6)
	meta := map[string]string{"color": "red"}
	d2 := d.WithMeta(meta)
	if !reflect.DeepEqual(d2.Metadata, meta) {
		t.Errorf("Metadata = %v, want %v", d2.Metadata, meta)
	}
	if d.Metadata != nil {
		t.Errorf("Original Die was mutated")
	}
}

func TestNewDicepool(t *testing.T) {
	d1 := NewDice(6)
	d2 := NewDice(8)
	dp := NewDicepool(d1, d2)
	if len(dp.Dice) != 2 {
		t.Errorf("Dice count = %d, want 2", len(dp.Dice))
	}
	if dp.Dice[0].Faces != 6 || dp.Dice[1].Faces != 8 {
		t.Errorf("Dice faces: got %d,%d want 6,8", dp.Dice[0].Faces, dp.Dice[1].Faces)
	}
	if len(dp.Modifiers) != 1 || dp.Modifiers[0] != (None{}) {
		t.Errorf("Modifiers should default to [None], got %v", dp.Modifiers)
	}
	if dp.Metadata == nil {
		t.Errorf("Metadata should be initialized non-nil")
	}
}

func TestDicepoolWithMods(t *testing.T) {
	dp := NewDicepool(NewDice(6))
	mods := []Mod{Sum{}, AddToEach{value: 1}}
	dp2 := dp.WithMods(mods...)
	if !reflect.DeepEqual(dp2.Modifiers, mods) {
		t.Errorf("Modifiers = %v, want %v", dp2.Modifiers, mods)
	}
	// Original unchanged
	if len(dp.Modifiers) != 1 || dp.Modifiers[0] != (None{}) {
		t.Errorf("Original Dicepool was mutated")
	}
}

func TestDicepoolWithMeta(t *testing.T) {
	dp := NewDicepool(NewDice(6))
	meta := map[string]string{"name": "test"}
	dp2 := dp.WithMeta(meta)
	if !reflect.DeepEqual(dp2.Metadata, meta) {
		t.Errorf("Metadata = %v, want %v", dp2.Metadata, meta)
	}
	// Original unchanged
	if dp.Metadata["name"] != "" {
		t.Errorf("Original Dicepool was mutated")
	}
}

// ---------------------------------------------------------------------
// Tests for mods.go
// ---------------------------------------------------------------------

func TestNoneMod(t *testing.T) {
	m := None{}
	input := []int{1, 2, 3}
	got := m.Apply(input)
	if !reflect.DeepEqual(got, input) {
		t.Errorf("None.Apply(%v) = %v, want %v", input, got, input)
	}
	if m.Priority() != PriorityNone {
		t.Errorf("None.Priority() = %d, want %d", m.Priority(), PriorityNone)
	}
}

func TestSumMod(t *testing.T) {
	m := Sum{}
	input := []int{1, 2, 3}
	got := m.Apply(input)
	want := []int{6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sum.Apply(%v) = %v, want %v", input, got, want)
	}
	if m.Priority() != PrioritySum {
		t.Errorf("Sum.Priority() = %d, want %d", m.Priority(), PrioritySum)
	}
	// Empty slice
	got = m.Apply([]int{})
	want = []int{0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sum.Apply([]) = %v, want %v", got, want)
	}
}

func TestAddConst(t *testing.T) {
	m := AddConst{value: 5}
	input := []int{1, 2, 3}
	got := m.Apply(input)
	want := []int{6, 7, 8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddConst{5}.Apply(%v) = %v, want %v", input, got, want)
	}
	if m.Priority() != PriorityAddToSum {
		t.Errorf("AddConst.Priority() = %d, want %d", m.Priority(), PriorityAddToSum)
	}
	// Negative
	m2 := AddConst{value: -2}
	got = m2.Apply([]int{10})
	want = []int{8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddConst{-2}.Apply([10]) = %v, want %v", got, want)
	}
}

func TestAddToEach(t *testing.T) {
	m := AddToEach{value: 3}
	input := []int{1, 2, 3}
	got := m.Apply(input)
	want := []int{4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddToEach{3}.Apply(%v) = %v, want %v", input, got, want)
	}
	if m.Priority() != PriorityAddToEach {
		t.Errorf("AddToEach.Priority() = %d, want %d", m.Priority(), PriorityAddToEach)
	}
}

func TestAddIndividual(t *testing.T) {
	tests := []struct {
		name  string
		m     AddIndividual
		input []int
		want  []int
	}{
		{"add at position 0", AddIndividual{position: 0, value: 10}, []int{1, 2, 3}, []int{11, 2, 3}},
		{"add at position 1", AddIndividual{position: 1, value: -5}, []int{1, 2, 3}, []int{1, -3, 3}},
		{"position out of range", AddIndividual{position: 5, value: 1}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty slice", AddIndividual{position: 0, value: 5}, []int{}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Apply(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddIndividual{%d,%d}.Apply(%v) = %v, want %v", tt.m.position, tt.m.value, tt.input, got, tt.want)
			}
		})
	}
	m := AddIndividual{position: 1, value: 2}
	if m.Priority() != PriorityAddIndividual {
		t.Errorf("AddIndividual.Priority() = %d, want %d", m.Priority(), PriorityAddIndividual)
	}
}

func TestDropLowest(t *testing.T) {
	tests := []struct {
		name  string
		m     DropLowest
		input []int
		want  []int
	}{
		{"drop 1 from 3", DropLowest{1}, []int{5, 1, 3}, []int{3, 5}},
		{"drop 2 from 4", DropLowest{2}, []int{4, 1, 3, 2}, []int{3, 4}},
		{"drop more than available", DropLowest{5}, []int{1, 2}, []int{2}}, // len-1 =1, drop min(5,1)=1 -> [2]
		{"drop 0", DropLowest{0}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"negative quantity", DropLowest{-1}, []int{1, 2}, []int{1, 2}}, // min(-1,1) = -1? Actually min(-1,1) = -1, drop <=0 => clone
		{"empty slice", DropLowest{1}, []int{}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Apply(tt.input)
			// Sort the result because Apply sorts the slice (it does sort out, but out is sorted)
			// Actually the implementation sorts out, so we compare sorted.
			// But we want to allow any order? The doc says drop lowest, but then sorts.
			// We'll compare sorted slices.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DropLowest{%d}.Apply(%v) = %v, want %v", tt.m.quantity, tt.input, got, tt.want)
			}
		})
	}
	m := DropLowest{1}
	if m.Priority() != PriorityDropLowest {
		t.Errorf("DropLowest.Priority() = %d, want %d", m.Priority(), PriorityDropLowest)
	}
}

func TestDropHighest(t *testing.T) {
	tests := []struct {
		name  string
		m     DropHighest
		input []int
		want  []int
	}{
		{"drop 1 from 3", DropHighest{1}, []int{5, 1, 3}, []int{1, 3}},
		{"drop 2 from 4", DropHighest{2}, []int{4, 1, 3, 2}, []int{1, 2}},
		{"drop more than available", DropHighest{5}, []int{1, 2}, []int{1}}, // len-1=1 -> drop 1 -> [1]
		{"drop 0", DropHighest{0}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty slice", DropHighest{1}, []int{}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Apply(tt.input)
			// The implementation sorts and then takes prefix.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DropHighest{%d}.Apply(%v) = %v, want %v", tt.m.quantity, tt.input, got, tt.want)
			}
		})
	}
	m := DropHighest{1}
	if m.Priority() != PriorityDropHighest {
		t.Errorf("DropHighest.Priority() = %d, want %d", m.Priority(), PriorityDropHighest)
	}
}

func TestDivide(t *testing.T) {
	m := Divide{value: 2}
	input := []int{1, 2, 3}
	got := m.Apply(input)
	want := []int{0, 1, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Divide{2}.Apply(%v) = %v, want %v", input, got, want)
	}
	if m.Priority() != PriorityDivide {
		t.Errorf("Divide.Priority() = %d, want %d", m.Priority(), PriorityDivide)
	}
	// Division by 1
	m2 := Divide{1}
	got = m2.Apply([]int{5})
	want = []int{5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Divide{1}.Apply([5]) = %v, want %v", got, want)
	}
}

func TestMultiply(t *testing.T) {
	m := Multiply{value: 3}
	input := []int{1, 2, 3}
	got := m.Apply(input)
	want := []int{3, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Multiply{3}.Apply(%v) = %v, want %v", input, got, want)
	}
	if m.Priority() != PriorityMultiply {
		t.Errorf("Multiply.Priority() = %d, want %d", m.Priority(), PriorityMultiply)
	}
	// Zero multiplication
	m2 := Multiply{0}
	got = m2.Apply([]int{5})
	want = []int{0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Multiply{0}.Apply([5]) = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------
// Tests for parse.go
// ---------------------------------------------------------------------

func TestParseSimpleAdditive(t *testing.T) {
	tests := []struct {
		input string
		want  int
		err   bool
	}{
		{"+5", 5, false},
		{"-3", -3, false},
		{"   +12  ", 12, false},
		{"0", 0, false},
		{"", 0, false},
		{"abc", 0, true},
		{"++5", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseSimpleAdditive(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("parseSimpleAdditive(%q) error = %v, wantErr %v", tt.input, err, tt.err)
			}
			if !tt.err && got != tt.want {
				t.Errorf("parseSimpleAdditive(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDicePart(t *testing.T) {
	tests := []struct {
		name         string
		part         string
		wantCount    int
		wantDiceType string
		wantSides    int
		wantLeftover string
		wantErr      bool
	}{
		{"simple d6", "d6", 1, "d", 6, "", false},
		{"2d20", "2d20", 2, "d", 20, "", false},
		{"D8", "D8", 1, "D", 8, "", false},
		{"3DD12", "3DD12", 3, "DD", 12, "", false},
		{"d10+5", "d10+5", 1, "d", 10, "+5", false},
		{"2d6-1", "2d6-1", 2, "d", 6, "-1", false},
		{"invalid count", "ad6", 0, "", 0, "", true},
		{"missing sides", "d", 0, "", 0, "", true},
		{"sides not number", "dX", 0, "", 0, "", true},
		{"sides <2", "d1", 0, "", 0, "", true},
		{"bad leftover", "2d6+abc", 2, "d", 6, "+abc", true}, // parseSimpleAdditive will catch later
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, diceType, sides, leftover, err := parseDicePart(tt.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDicePart(%q) error = %v, wantErr %v", tt.part, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if count != tt.wantCount {
				t.Errorf("count = %d, want %d", count, tt.wantCount)
			}
			if diceType != tt.wantDiceType {
				t.Errorf("diceType = %s, want %s", diceType, tt.wantDiceType)
			}
			if sides != tt.wantSides {
				t.Errorf("sides = %d, want %d", sides, tt.wantSides)
			}
			if leftover != tt.wantLeftover {
				t.Errorf("leftover = %q, want %q", leftover, tt.wantLeftover)
			}
		})
	}
}

func TestParseOneComplexModifier(t *testing.T) {
	tests := []struct {
		input string
		want  Mod
		err   bool
	}{
		{"+5e", AddToEach{5}, false},
		{"-3e", AddToEach{-3}, false},
		{"+10to2", AddIndividual{position: 2, value: 10}, false},
		{"-4to1", AddIndividual{position: 1, value: -4}, false},
		{"dl2", DropLowest{2}, false},
		{"dh1", DropHighest{1}, false},
		{"/3", Divide{3}, false},
		{"x4", Multiply{4}, false},
		{"*5", Multiply{5}, false},
		{"invalid", nil, true},
		{"dl", nil, true},
		{"to5", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseOneComplexModifier(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("parseOneComplexModifier(%q) error = %v, wantErr %v", tt.input, err, tt.err)
				return
			}
			if tt.err {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOneComplexModifier(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseComplexModifiers(t *testing.T) {
	tests := []struct {
		input string
		want  []Mod
		err   bool
	}{
		{"", nil, false},
		{"+5e", []Mod{AddToEach{5}}, false},
		{"dl2:dh1", []Mod{DropLowest{2}, DropHighest{1}}, false},
		{"+5e:/2", []Mod{AddToEach{5}, Divide{2}}, false},
		{"invalid", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseComplexModifiers(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("parseComplexModifiers(%q) error = %v, wantErr %v", tt.input, err, tt.err)
				return
			}
			if tt.err {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseComplexModifiers(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSortModifiers(t *testing.T) {
	mods := []Mod{
		DropLowest{2},       // priority 70
		AddToEach{1},        // 30
		Sum{},               // 100
		None{},              // 0
		AddIndividual{0, 5}, // 20
		Divide{2},           // 120
	}
	sorted := sortModifiers(mods)
	priorities := make([]int, len(sorted))
	for i, m := range sorted {
		priorities[i] = m.Priority()
	}
	// Should be non-decreasing
	for i := 1; i < len(priorities); i++ {
		if priorities[i-1] > priorities[i] {
			t.Errorf("sortModifiers not stable by priority: %v", priorities)
		}
	}
	// First should be None (0)
	if _, ok := sorted[0].(None); !ok {
		t.Errorf("first modifier should be None, got %T", sorted[0])
	}
	// Last should be Divide (120)
	if _, ok := sorted[len(sorted)-1].(Divide); !ok {
		t.Errorf("last modifier should be Divide, got %T", sorted[len(sorted)-1])
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		wantDiceCnt int
		wantFaces   int
		wantMods    []Mod // we'll compare types and values
		wantErr     bool
	}{
		{
			name:        "simple d6",
			expr:        "d6",
			wantDiceCnt: 1,
			wantFaces:   6,
			wantMods:    []Mod{Sum{}},
			wantErr:     false,
		},
		{
			name:        "2d20+5",
			expr:        "2d20+5",
			wantDiceCnt: 2,
			wantFaces:   20,
			wantMods:    []Mod{Sum{}, AddConst{5}}, // order after sorting: Sum(100), AddConst(110) -> Sum, AddConst? Actually priority: Sum 100, AddConst 110, so Sum first then AddConst
			wantErr:     false,
		},
		{
			name:        "d10:dl1",
			expr:        "d10:dl1",
			wantDiceCnt: 1,
			wantFaces:   10,
			wantMods:    []Mod{DropLowest{1}, Sum{}},
			wantErr:     false,
		},
		{
			name:        "3d6:+5e:/2",
			expr:        "3d6:+5e:/2",
			wantDiceCnt: 3,
			wantFaces:   6,
			wantMods:    []Mod{AddToEach{5}, Sum{}, Divide{2}}, // order: AddToEach(30), Sum(100), Divide(120)? Actually Divide 120 > Sum 100, so after Sum? Wait: AddToEach 30, Sum 100, Divide 120 -> sorted: AddToEach, Sum, Divide. But Apply order: AddToEach, then Sum, then Divide. That's fine.
			wantErr:     false,
		},
		{
			name:        "D20", // special type
			expr:        "D20",
			wantDiceCnt: 1,
			wantFaces:   0, // because NewDice(0) for special
			wantMods:    []Mod{Sum{}},
			wantErr:     false,
		},
		{
			name:    "empty",
			expr:    "",
			wantErr: true,
		},
		{
			name:    "invalid dice type",
			expr:    "x6",
			wantErr: true,
		},
		{
			name:    "invalid modifier",
			expr:    "d6:invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp, err := ParseExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpression(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(dp.Dice) != tt.wantDiceCnt {
				t.Errorf("Dice count = %d, want %d", len(dp.Dice), tt.wantDiceCnt)
			}
			if tt.wantFaces != 0 && dp.Dice[0].Faces != tt.wantFaces {
				t.Errorf("Faces = %d, want %d", dp.Dice[0].Faces, tt.wantFaces)
			}
			// Compare modifiers: we need to ignore the None that is added by NewDicepool? Actually NewDicepool adds None as default, but WithMods replaces. So dp.Modifiers should exactly match the list after sorting.
			// However, ParseExpression adds Sum and any others, then sorts. We'll compare lengths and types/values.
			if len(dp.Modifiers) != len(tt.wantMods) {
				t.Errorf("Modifiers count = %d, want %d: %v", len(dp.Modifiers), len(tt.wantMods), dp.Modifiers)
				return
			}
			for i, m := range dp.Modifiers {
				if !reflect.DeepEqual(m, tt.wantMods[i]) {
					t.Errorf("Modifier[%d] = %v, want %v", i, m, tt.wantMods[i])
				}
			}
		})
	}
}

// ---------------------------------------------------------------------
// Tests for roll.go (requires deterministic RNG)
// ---------------------------------------------------------------------

func TestRollerRollDice(t *testing.T) {
	// Use fixed seed for deterministic results
	roller := newTestRoller("fixed")
	die := NewDice(6)
	// Since RNG is deterministic, we can check that first roll is known
	// Actually we don't know the exact sequence, but we can test bounds.
	for i := 0; i < 100; i++ {
		val := roller.rollDice(die)
		if val < 1 || val > 6 {
			t.Errorf("rollDice returned %d, want between 1 and 6", val)
		}
	}
}

func TestRollerRollPool(t *testing.T) {
	roller := newTestRoller("fixed")
	// Build a dicepool with known dice and mods
	dice := []Die{NewDice(6), NewDice(6)}
	mods := []Mod{AddToEach{1}, Sum{}, Multiply{2}} // add 1 to each, sum, then multiply by 2
	dp := NewDicepool(dice...).WithMods(mods...)

	result := roller.rollPool(dp)
	// Check Raw length
	if len(result.Raw) != 2 {
		t.Errorf("Raw length = %d, want 2", len(result.Raw))
	}
	// Each raw roll between 1 and 6
	for _, v := range result.Raw {
		if v < 1 || v > 6 {
			t.Errorf("Raw roll %d out of bounds", v)
		}
	}
	// Final should have 1 element (after Sum)
	if len(result.Final) != 1 {
		t.Errorf("Final length = %d, want 1", len(result.Final))
	}
	// Compute expected: raw1+1 + raw2+1 = sum+2, then *2 = 2*(sum+2)
	sumRaw := result.Raw[0] + result.Raw[1]
	expected := 2 * (sumRaw + 2)
	if result.Final[0] != expected {
		t.Errorf("Final[0] = %d, expected %d (raw=%v)", result.Final[0], expected, result.Raw)
	}
}

// Integration test: parse + roll
func TestIntegrationParseAndRoll(t *testing.T) {
	roller := newTestRoller("integration")
	expr := "2d6:+2e:/2" // two d6, add 2 to each, then divide by 2 (integer division)
	dp, err := ParseExpression(expr)
	if err != nil {
		t.Fatalf("ParseExpression(%q) error: %v", expr, err)
	}
	// Expect modifiers: AddToEach{2}, Sum{}, Divide{2} after sorting? Priorities: AddToEach 30, Sum 100, Divide 120 -> order: AddToEach, Sum, Divide.
	// So: raw rolls -> add 2 to each -> sum -> divide by 2.
	result := roller.rollPool(dp)
	if len(result.Final) != 1 {
		t.Errorf("Final length = %d, want 1", len(result.Final))
	}
	// Compute expected: ((roll1+2)+(roll2+2))/2 = (roll1+roll2+4)/2
	sumRaw := result.Raw[0] + result.Raw[1]
	expected := (sumRaw + 4) / 2
	if result.Final[0] != expected {
		t.Errorf("Final[0] = %d, expected %d (raw=%v)", result.Final[0], expected, result.Raw)
	}
}

// Test for special dice types (D/DD) – they have Faces=0 and will panic if rolled.
// So we test parsing only, not rolling.
func TestSpecialDiceParsing(t *testing.T) {
	expr := "D20"
	dp, err := ParseExpression(expr)
	if err != nil {
		t.Fatalf("ParseExpression(D20) error: %v", err)
	}
	if len(dp.Dice) != 1 {
		t.Fatalf("Expected 1 die, got %d", len(dp.Dice))
	}
	die := dp.Dice[0]
	if die.Faces != 0 {
		t.Errorf("Special die Faces = %d, want 0", die.Faces)
	}
	if die.Metadata["special"] != "D" {
		t.Errorf("Special die metadata[special] = %q, want 'D'", die.Metadata["special"])
	}
	// Ensure that rolling would panic? We'll skip rolling test.
}
