package tables

import (
	"errors"
	"slices"
	"testing"
)

// ---------------------------------------------------------------------
// Test helpers and mock roller
// ---------------------------------------------------------------------

// mockRoller implements TableRoller for testing.
type mockRoller struct {
	d66Result   string
	rollResults map[string]int // expression -> result
	rollErr     error
}

func (m *mockRoller) D66(...int) string {
	return m.d66Result
}

func (m *mockRoller) Roll(expr string, mods ...int) (int, error) {
	if m.rollErr != nil {
		return 0, m.rollErr
	}
	if m.rollResults == nil {
		return 0, nil
	}
	if val, ok := m.rollResults[expr]; ok {
		return val, nil
	}
	// Default: return 0
	return 0, nil
}

// Helper to generate []int from n down to DefaultLowerBound.
func evenOrLessThan(n int) []int {
	nums := []int{}
	for i := n; i >= DefaultLowerBound; i-- {
		nums = append(nums, i)
	}
	slices.Sort(nums)
	return nums
}

// Helper to generate []int from n up to DefaultUpperBound.
func evenOrMoreThan(n int) []int {
	nums := []int{}
	for i := n; i <= DefaultUpperBound; i++ {
		nums = append(nums, i)
	}
	slices.Sort(nums)
	return nums
}

// ---------------------------------------------------------------------
// Tests for GameTable.Validate
// ---------------------------------------------------------------------

func TestGameTableValidate(t *testing.T) {
	tests := []struct {
		name    string
		table   GameTable
		wantErr bool
	}{
		{
			name:    "valid table",
			table:   New("test", "d6", map[string]string{"1": "a", "2": "b", "3": "c", "4": "d", "5": "e", "6": "f"}),
			wantErr: false,
		},
		{
			name:    "empty name",
			table:   New("", "d6", map[string]string{"1": "a", "2": "b"}),
			wantErr: true,
		},
		{
			name:    "only one entry",
			table:   New("test", "d6", map[string]string{"1": "a"}),
			wantErr: true,
		},
		{
			name:    "empty value",
			table:   New("test", "d6", map[string]string{"1": "a", "2": "", "3": "c"}),
			wantErr: true,
		},
		{
			name:    "hole in range for non-D66",
			table:   New("test", "d6", map[string]string{"1": "a", "2": "b", "4": "d", "5": "e", "6": "f"}),
			wantErr: true,
		},
		{
			name:    "index out of bounds",
			table:   New("test", "d6", map[string]string{"1": "a", "2": "b", "1001": "f"}),
			wantErr: true,
		},
		{
			name:    "invalid dice expression",
			table:   New("test", "xyz", map[string]string{"1": "a", "2": "b"}),
			wantErr: true,
		},
		{
			name:    "valid D66 table with all 36 entries",
			table:   New("test", "d66", map[string]string{"11": "a", "12": "b", "13": "c", "14": "d", "15": "e", "16": "f", "21": "g", "22": "h", "23": "i", "24": "j", "25": "k", "26": "l", "31": "m", "32": "n", "33": "o", "34": "p", "35": "q", "36": "r", "41": "s", "42": "t", "43": "u", "44": "v", "45": "w", "46": "x", "51": "y", "52": "z", "53": "aa", "54": "ab", "55": "ac", "56": "ad", "61": "ae", "62": "af", "63": "ag", "64": "ah", "65": "ai", "66": "aj"}),
			wantErr: false,
		},
		{
			// D66 tables are allowed to have holes – this should NOT error.
			name:    "D66 table with holes is valid",
			table:   New("test", "d66", map[string]string{"11": "a", "12": "b", "66": "c"}),
			wantErr: false,
		},
		{
			name:    "duplicate index",
			table:   New("test", "d6", map[string]string{"1": "a", "2": "b", "3": "c", "3 - 4": "d"}),
			wantErr: true,
		},
		{
			name:    "valid table with modifier",
			table:   New("test", "2d6+1", map[string]string{"3": "a", "4": "b", "5": "c", "6": "d", "7": "e", "8": "f", "9": "g", "10": "h", "11": "i", "12": "j", "13": "k"}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.table.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------
// Tests for validateExpression (internal function)
// ---------------------------------------------------------------------

func TestValidateExpression(t *testing.T) {
	tests := []struct {
		expr    string
		wantErr bool
	}{
		{"d6", false},
		{"2d10", false},
		{"3d6+2", false},
		{"d20-1", false},
		{"D66", false},
		{"d66", false},
		{"", true},
		{"d", true},
		{"d0", true},
		{"0d6", true},
		{"d6+", true},
		{"d6++", true},
		{"d6-", true},
		{"xd6", true},
		{"d6.5", true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			err := validateExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateExpression(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------
// Tests for indexesToString and stringToIndexes
// ---------------------------------------------------------------------

func TestIndexesToString(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		want    string
		wantErr bool
	}{
		{"empty", []int{}, "", false},
		{"single", []int{5}, "5", false},
		{"multiple singles", []int{1, 3, 5}, "1, 3, 5", false},
		{"consecutive range", []int{2, 3, 4, 5}, "2 - 5", false},
		{"mixed ranges", []int{2, 3, 4, 5, 8, 9, 15, 21}, "2 - 5, 8 - 9, 15, 21", false},
		{"andAbove valid", []int{2, andAbove}, "2+", false},
		{"andBelow valid", []int{-50, andBelow}, "-50-", false},
		{"andAbove wrong args", []int{2, andAbove, 5}, "", true},
		{"andBelow wrong args", []int{andBelow}, "", true},
		{"unsorted input", []int{10, 5, 7, 6}, "5 - 7, 10", false},
		{"negative range", []int{-5, -4, -3, -2}, "-5 - -2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := indexesToString(tt.input...)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexesToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("indexesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToIndexes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{"empty", "", nil, true},
		{"whitespace only", "   ", nil, true},
		{"single number", "5", []int{5}, false},
		{"negative single", "-3", []int{-3}, false},
		{"range", "2 - 4", []int{2, 3, 4}, false},
		{"negative range", "-3 - -1", []int{-3, -2, -1}, false},
		{"mixed", "2 - 4, 6, 8+", append([]int{2, 3, 4, 6}, evenOrMoreThan(8)...), false},
		{"plus range expands to upper bound", "98+", evenOrMoreThan(98), false},
		{"minus range expands to lower bound", "0-", evenOrLessThan(0), false},
		{"invalid range min gt max", "5 - 3", nil, true},
		{"invalid token", "2 - 4, abc", nil, true},
		{"duplicates should be deduplicated", "1, 2, 1", []int{1, 2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToIndexes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("stringToIndexes() len = %d, want %d", len(got), len(tt.want))
					return
				}
				for i, v := range got {
					if v != tt.want[i] {
						t.Errorf("stringToIndexes()[%d] = %v, want %v", i, v, tt.want[i])
						return
					}
				}
			}
		})
	}
}

// ---------------------------------------------------------------------
// Tests for Collection
// ---------------------------------------------------------------------

func TestNewCollection(t *testing.T) {
	t.Run("valid collection", func(t *testing.T) {
		table1 := New("table1", "d6", map[string]string{"1": "a", "2": "b", "3": "c", "4": "d", "5": "e", "6": "f"})
		table2 := New("table2", "2d6", map[string]string{"2": "a", "3": "b", "4": "c", "5": "d", "6": "e", "7": "f", "8": "g", "9": "h", "10": "i", "11": "j", "12": "k"})
		coll, err := NewCollection("testcoll", table1, table2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if coll.Name != "testcoll" {
			t.Errorf("name = %q, want %q", coll.Name, "testcoll")
		}
		if len(coll.Tables) != 2 {
			t.Errorf("got %d tables, want 2", len(coll.Tables))
		}
	})

	t.Run("duplicate table names", func(t *testing.T) {
		table1 := New("same", "d6", map[string]string{"1": "a", "2": "b"})
		table2 := New("same", "2d6", map[string]string{"2": "a", "3": "b"})
		_, err := NewCollection("dupe", table1, table2)
		if err == nil {
			t.Fatal("expected error for duplicate table names, got nil")
		}
	})

	t.Run("empty collection name", func(t *testing.T) {
		table := New("t", "d6", map[string]string{"1": "a", "2": "b"})
		_, err := NewCollection("", table)
		if err == nil {
			t.Fatal("expected error for empty collection name, got nil")
		}
	})

	t.Run("no tables provided", func(t *testing.T) {
		_, err := NewCollection("empty")
		if err == nil {
			t.Fatal("expected error for no tables, got nil")
		}
	})

	t.Run("invalid table inside collection", func(t *testing.T) {
		badTable := New("bad", "d6", map[string]string{"1": "a"}) // only one entry
		_, err := NewCollection("withbad", badTable)
		if err == nil {
			t.Fatal("expected error due to invalid table, got nil")
		}
	})
}

func TestCollectionRoll(t *testing.T) {
	// Prepare a valid collection with normal and D66 tables.
	normalTable := New("normal", "2d6", map[string]string{
		"2":  "snake eyes",
		"3":  "three",
		"4":  "four",
		"5":  "five",
		"6":  "six",
		"7":  "seven",
		"8":  "eight",
		"9":  "nine",
		"10": "ten",
		"11": "eleven",
		"12": "boxcars",
	})
	d66Table := New("d66", "d66", map[string]string{
		"11": "one one",
		"12": "one two",
		"66": "six six",
	})
	coll, err := NewCollection("test", normalTable, d66Table)
	if err != nil {
		t.Fatalf("failed to create test collection: %v", err)
	}

	t.Run("roll normal table", func(t *testing.T) {
		roller := &mockRoller{
			rollResults: map[string]int{"2d6": 7},
		}
		result, err := coll.Roll(roller, "normal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "seven" {
			t.Errorf("got %q, want %q", result, "seven")
		}
		// Check that rollSequence and results were recorded.
		if len(coll.rollSequence) != 1 || coll.rollSequence[0] != "normal" {
			t.Errorf("rollSequence = %v, want [normal]", coll.rollSequence)
		}
		if len(coll.results) != 1 || coll.results[0] != "seven" {
			t.Errorf("results = %v, want [seven]", coll.results)
		}
	})

	t.Run("roll D66 table", func(t *testing.T) {
		coll.Reset() // clear previous history
		roller := &mockRoller{
			d66Result: "12",
		}
		result, err := coll.Roll(roller, "d66")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "one two" {
			t.Errorf("got %q, want %q", result, "one two")
		}
	})

	t.Run("missing table", func(t *testing.T) {
		roller := &mockRoller{}
		_, err := coll.Roll(roller, "nonexistent")
		if err == nil {
			t.Fatal("expected error for missing table, got nil")
		}
	})

	t.Run("nil roller", func(t *testing.T) {
		_, err := coll.Roll(nil, "normal")
		if err == nil {
			t.Fatal("expected error for nil roller, got nil")
		}
	})

	t.Run("roller error", func(t *testing.T) {
		roller := &mockRoller{
			rollErr: errors.New("dice exploded"),
		}
		_, err := coll.Roll(roller, "normal")
		if err == nil {
			t.Fatal("expected error from roller, got nil")
		}
	})

	t.Run("result not found in table (out-of-range roll)", func(t *testing.T) {
		// For normal table, roll a value not covered by keys.
		roller := &mockRoller{
			rollResults: map[string]int{"2d6": 13},
		}
		_, err := coll.Roll(roller, "normal")
		if err == nil {
			t.Fatal("expected error for missing result, got nil")
		}
	})
}

func TestCollectionRollCascade(t *testing.T) {
	// Create tables where results point to other tables.
	tableA := New("A", "d6", map[string]string{
		"1": "B",
		"2": "C",
		"3": "end",
		"4": "end",
		"5": "end",
		"6": "end",
	})
	tableB := New("B", "d6", map[string]string{
		"1": "C",
		"2": "end",
		"3": "end",
		"4": "end",
		"5": "end",
		"6": "end",
	})
	tableC := New("C", "d6", map[string]string{
		"1": "final",
		"2": "final",
		"3": "final",
		"4": "final",
		"5": "final",
		"6": "final",
	})
	coll, err := NewCollection("cascade", tableA, tableB, tableC)
	if err != nil {
		t.Fatalf("failed to create collection: %v", err)
	}

	t.Run("cascade stops when result is not a table name", func(t *testing.T) {
		roller := &mockRoller{
			rollResults: map[string]int{
				"d6": 3, // from A -> "end"
			},
		}
		result, err := coll.RollCascade(roller, "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "end" {
			t.Errorf("got %q, want %q", result, "end")
		}
		// Should have rolled only on A.
		if len(coll.rollSequence) != 1 || coll.rollSequence[0] != "A" {
			t.Errorf("rollSequence = %v, want [A]", coll.rollSequence)
		}
	})

	t.Run("cascade follows chain", func(t *testing.T) {
		coll.Reset()
		// Simulate: A rolls 1 -> B, B rolls 1 -> C, C rolls 1 -> "final"
		roller := &mockRoller{
			rollResults: map[string]int{
				"d6": 1,
			},
		}
		result, err := coll.RollCascade(roller, "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "final" {
			t.Errorf("got %q, want %q", result, "final")
		}
		expectedSeq := []string{"A", "B", "C"}
		for i, name := range expectedSeq {
			if i >= len(coll.rollSequence) || coll.rollSequence[i] != name {
				t.Errorf("rollSequence[%d] = %q, want %q", i, coll.rollSequence[i], name)
			}
		}
	})

	t.Run("cascade error from Roll propagates", func(t *testing.T) {
		coll.Reset()
		roller := &mockRoller{
			rollErr: errors.New("boom"),
		}
		_, err := coll.RollCascade(roller, "A")
		if err == nil {
			t.Fatal("expected error from Roll, got nil")
		}
	})

	t.Run("cascade max depth exceeded", func(t *testing.T) {
		// Create a loop: table Loop points to itself.
		loopTable := New("Loop", "d6", map[string]string{"1": "Loop", "2": "Loop"})
		coll2, _ := NewCollection("loopcoll", loopTable)
		roller := &mockRoller{
			rollResults: map[string]int{"d6": 1},
		}
		_, err := coll2.RollCascade(roller, "Loop")
		if err == nil {
			t.Fatal("expected max depth error, got nil")
		}
	})
}

func TestCollectionReset(t *testing.T) {
	table := New("t", "d6", map[string]string{"1": "a", "2": "b"})
	coll, _ := NewCollection("test", table)
	roller := &mockRoller{
		rollResults: map[string]int{"d6": 1},
	}
	coll.Roll(roller, "t")
	coll.Roll(roller, "t")

	if len(coll.rollSequence) != 2 {
		t.Fatalf("expected 2 rolls, got %d", len(coll.rollSequence))
	}
	coll.Reset()
	if len(coll.rollSequence) != 0 || len(coll.results) != 0 {
		t.Errorf("Reset did not clear sequences: seq=%v, results=%v", coll.rollSequence, coll.results)
	}
}

func TestCollectionValidate(t *testing.T) {
	validTable := New("t", "d6", map[string]string{"1": "a", "2": "b"})

	t.Run("valid collection", func(t *testing.T) {
		coll, _ := NewCollection("good", validTable)
		if err := coll.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		coll := &Collection{Name: "", Tables: map[string]GameTable{"t": validTable}}
		if err := coll.Validate(); err == nil {
			t.Error("expected error for empty name, got nil")
		}
	})

	t.Run("no tables", func(t *testing.T) {
		coll := &Collection{Name: "notables", Tables: map[string]GameTable{}}
		if err := coll.Validate(); err == nil {
			t.Error("expected error for missing tables, got nil")
		}
	})

	t.Run("invalid table", func(t *testing.T) {
		badTable := New("bad", "d6", map[string]string{"1": "a"}) // only one entry
		coll := &Collection{Name: "badcoll", Tables: map[string]GameTable{"bad": badTable}}
		if err := coll.Validate(); err == nil {
			t.Error("expected error for invalid table, got nil")
		}
	})
}
