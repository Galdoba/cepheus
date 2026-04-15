package ehex

import (
	"testing"
)

func TestValueToCodeMapping(t *testing.T) {
	t.Run("completeness", func(t *testing.T) {
		for v := 0; v <= 33; v++ {
			if _, ok := valueToCode[v]; !ok {
				t.Errorf("valueToCode missing entry for %d", v)
			}
		}
	})

	t.Run("uniqueness", func(t *testing.T) {
		seen := make(map[string]int)
		for v, code := range valueToCode {
			if prev, exists := seen[code]; exists {
				t.Errorf("duplicate code %q for values %d and %d", code, prev, v)
			}
			seen[code] = v
			if len(code) != 1 {
				t.Errorf("code %q for value %d is not a single character", code, v)
			}
		}
	})
}

func TestCodeToValueConsistency(t *testing.T) {
	for code, val := range codeToValue {
		expectedCode, ok := valueToCode[val]
		if !ok {
			t.Errorf("valueToCode missing entry for value %d (from code %q)", val, code)
		} else if expectedCode != code {
			t.Errorf("mismatch: codeToValue[%q]=%d but valueToCode[%d]=%q", code, val, val, expectedCode)
		}
	}
}

func TestFromValue(t *testing.T) {
	t.Run("valid range 0-33", func(t *testing.T) {
		for v := 0; v <= 33; v++ {
			e := FromValue(v)
			if e.Value() != v {
				t.Errorf("FromValue(%d).Value() = %d", v, e.Value())
			}
			expectedCode := valueToCode[v]
			if e.Code() != expectedCode {
				t.Errorf("FromValue(%d).Code() = %q, want %q", v, e.Code(), expectedCode)
			}
			if e.Description() != "" {
				t.Errorf("FromValue(%d).Description() should be empty, got %q", v, e.Description())
			}
		}
	})

	t.Run("out of range returns Unknown", func(t *testing.T) {
		invalidVals := []int{-999, -100, -1, 34, 100, 999}
		for _, v := range invalidVals {
			e := FromValue(v)
			if e != Unknown {
				t.Errorf("FromValue(%d) = %+v, want Unknown", v, e)
			}
		}
	})

	t.Run("special values are not returned", func(t *testing.T) {
		specialVals := []int{-101, -102, -103, -104, -105, -106, -107, -108, -109}
		for _, v := range specialVals {
			e := FromValue(v)
			if e != Unknown {
				t.Errorf("FromValue(%d) returned special value, want Unknown", v)
			}
		}
	})
}

func TestFromCode(t *testing.T) {
	t.Run("basic codes 0-9 A-Z (excl I,O)", func(t *testing.T) {
		for v := 0; v <= 33; v++ {
			code := valueToCode[v]
			e := FromCode(code)
			if e.Value() != v {
				t.Errorf("FromCode(%q).Value() = %d, want %d", code, e.Value(), v)
			}
			if e.Code() != code {
				t.Errorf("FromCode(%q).Code() = %q", code, e.Code())
			}
		}
	})

	t.Run("special symbols", func(t *testing.T) {
		tests := []struct {
			code  string
			value int
			ehex  Ehex
		}{
			{"?", -101, Unknown},
			{"*", -102, Any},
			{"!", -103, Invalid},
			{"#", -104, Default},
			{"-", -105, Ignore},
			{"&", -106, Reserved},
			{"~", -107, Masked},
			{">", -108, Extension},
			{".", -109, Placeholder},
		}
		for _, tt := range tests {
			e := FromCode(tt.code)
			if e != tt.ehex {
				t.Errorf("FromCode(%q) = %+v, want %+v", tt.code, e, tt.ehex)
			}
		}
	})

	t.Run("extended aliases (s,r -> 0)", func(t *testing.T) {
		aliases := []string{"s", "r"}
		for _, code := range aliases {
			e := FromCode(code)
			if e.Value() != 0 {
				t.Errorf("FromCode(%q).Value() = %d, want 0", code, e.Value())
			}
			if e.Code() != code {
				t.Errorf("FromCode(%q).Code() = %q, want %q", code, e.Code(), code)
			}
			// These are custom Ehex instances, not predefined constants.
		}
	})

	t.Run("invalid codes return Unknown", func(t *testing.T) {
		invalid := []string{"", " ", "/", "\\", "[", "]", "|", "a", "i", "o", "AA", "??", "I", "O"}
		for _, code := range invalid {
			e := FromCode(code)
			if e != Unknown {
				t.Errorf("FromCode(%q) = %+v, want Unknown", code, e)
			}
		}
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("value->code->value (0-33)", func(t *testing.T) {
		for v := 0; v <= 33; v++ {
			e1 := FromValue(v)
			code := e1.Code()
			e2 := FromCode(code)
			if e2.Value() != v {
				t.Errorf("round-trip failed: %d -> %q -> %d", v, code, e2.Value())
			}
			if e2.Code() != code {
				t.Errorf("round-trip code mismatch: %d -> %q -> %q", v, code, e2.Code())
			}
		}
	})

	t.Run("code->value->code (basic)", func(t *testing.T) {
		for _, code := range []string{"0", "5", "A", "Z", "H", "P"} {
			e1 := FromCode(code)
			val := e1.Value()
			e2 := FromValue(val)
			// For standard codes, FromValue should return the same code
			if e2.Code() != code {
				t.Errorf("round-trip failed: %q -> %d -> %q", code, val, e2.Code())
			}
		}
	})

	t.Run("special symbols round-trip via FromCode only", func(t *testing.T) {
		specials := []string{"?", "*", "!", "#", "-", "&", "~", ">", "."}
		for _, sym := range specials {
			e1 := FromCode(sym)
			code := e1.Code()
			e2 := FromCode(code)
			if e2 != e1 {
				t.Errorf("special symbol %q round-trip mismatch: %+v -> %+v", sym, e1, e2)
			}
		}
	})
}

func TestSpecialConstants(t *testing.T) {
	specials := []struct {
		name  string
		ehex  Ehex
		code  string
		value int
		desc  string
	}{
		{"Unknown", Unknown, "?", -101, "unknown"},
		{"Any", Any, "*", -102, "any"},
		{"Invalid", Invalid, "!", -103, "invalid"},
		{"Default", Default, "#", -104, "default"},
		{"Ignore", Ignore, "-", -105, "ignore"},
		{"Reserved", Reserved, "&", -106, "reserved"},
		{"Masked", Masked, "~", -107, "masked"},
		{"Extension", Extension, ">", -108, "extension"},
		{"Placeholder", Placeholder, ".", -109, "placeholder"},
	}

	for _, sp := range specials {
		t.Run(sp.name, func(t *testing.T) {
			if sp.ehex.Code() != sp.code {
				t.Errorf("Code() = %q, want %q", sp.ehex.Code(), sp.code)
			}
			if sp.ehex.Value() != sp.value {
				t.Errorf("Value() = %d, want %d", sp.ehex.Value(), sp.value)
			}
			if sp.ehex.Description() != sp.desc {
				t.Errorf("Description() = %q, want %q", sp.ehex.Description(), sp.desc)
			}
			if sp.ehex.String() != sp.code {
				t.Errorf("String() = %q, want %q", sp.ehex.String(), sp.code)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("basic custom creation", func(t *testing.T) {
		e := New(100, "X", "custom description")
		if e.Value() != 100 {
			t.Errorf("Value() = %d, want 100", e.Value())
		}
		if e.Code() != "X" {
			t.Errorf("Code() = %q, want X", e.Code())
		}
		if e.Description() != "custom description" {
			t.Errorf("Description() = %q, want 'custom description'", e.Description())
		}
	})

	t.Run("without description", func(t *testing.T) {
		e := New(42, "?")
		if e.Value() != 42 {
			t.Errorf("Value() = %d, want 42", e.Value())
		}
		if e.Description() != "" {
			t.Errorf("Description() = %q, want empty", e.Description())
		}
	})

	t.Run("multiple descriptions ignored except first", func(t *testing.T) {
		e := New(7, "Z", "first", "second", "third")
		if e.Description() != "first" {
			t.Errorf("Description() = %q, want 'first'", e.Description())
		}
	})
}

func TestWithDescription(t *testing.T) {
	original := FromValue(5)
	modified := original.WithDescription("five")

	// Original unchanged
	if original.Description() != "" {
		t.Error("WithDescription mutated original instance")
	}
	// Modified has new description
	if modified.Description() != "five" {
		t.Errorf("WithDescription() description = %q, want 'five'", modified.Description())
	}
	// Other fields preserved
	if modified.Value() != original.Value() || modified.Code() != original.Code() {
		t.Error("WithDescription altered value or code")
	}

	// Works on custom instances too
	custom := New(99, "@", "initial").WithDescription("updated")
	if custom.Description() != "updated" {
		t.Errorf("custom description = %q, want 'updated'", custom.Description())
	}
}

func TestZeroValueEhex(t *testing.T) {
	var e Ehex
	if e.Code() != "" {
		t.Errorf("zero value Code() = %q, want empty string", e.Code())
	}
	if e.Value() != 0 {
		t.Errorf("zero value Value() = %d, want 0", e.Value())
	}
	if e.Description() != "" {
		t.Errorf("zero value Description() = %q, want empty", e.Description())
	}
	if e.String() != "" {
		t.Errorf("zero value String() = %q, want empty", e.String())
	}
}

func TestEhexEquality(t *testing.T) {
	// Ehex is comparable (all fields are comparable types)
	t.Run("identical instances", func(t *testing.T) {
		a := FromValue(10)
		b := FromValue(10)
		if a != b {
			t.Error("identical FromValue instances should be equal")
		}
	})

	t.Run("different values", func(t *testing.T) {
		if Unknown == Any {
			t.Error("Unknown and Any should not be equal")
		}
	})

	t.Run("custom equal", func(t *testing.T) {
		a := New(100, "@", "desc")
		b := New(100, "@", "desc")
		if a != b {
			t.Error("custom instances with same fields should be equal")
		}
	})

	t.Run("different description", func(t *testing.T) {
		a := FromValue(1).WithDescription("one")
		b := FromValue(1).WithDescription("1")
		if a == b {
			t.Error("instances with different descriptions should not be equal")
		}
	})
}
