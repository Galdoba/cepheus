package dice_test

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

func TestNew(t *testing.T) {
	// m, err := dice.New("seed 42")
	// if err != nil {
	// 	fmt.Println("err New:", err)
	// 	return
	// }
	// for i := range 5 {
	// 	fmt.Println("loop", i+1)
	// 	fmt.Println(dice.MustRoll("3d6"), dice.MustRoll("3d6"), dice.MustRoll("3d6"), dice.MustRoll("3d6"))
	// 	fmt.Println(m.MustRoll("3d6"), m.MustRoll("3d6"), m.MustRoll("3d6"), m.MustRoll("3d6"))
	// }
}

// ----------------------------------------------------------------------
// Helper: newSeededManager creates a manager with a fixed seed for determinism.
func newSeededManager(t *testing.T, seed string) *dice.Manager {
	t.Helper()
	m, err := dice.New(seed)
	if err != nil {
		t.Fatalf("dice.New(%q) failed: %v", seed, err)
	}
	return m
}

// ----------------------------------------------------------------------
// Expression Parsing Tests

func TestParseValidExpressions(t *testing.T) {
	tests := []struct {
		expr string
		desc string
	}{
		{"1d6", "single die"},
		{"2d20", "multiple dice"},
		{"10d100", "large dice count and sides"},
		{"3d6+5", "positive simple additive"},
		{"3d6-2", "negative simple additive"},
		{"4d8:dl1", "drop lowest one"},
		{"4d8:dh1", "drop highest one"},
		{"3d6:2e", "add 2 to each die"},
		{"3d6:2>>2", "add 2 to second die"},
		{"3d6:-1>>3", "subtract 1 from third die"},
		{"3d6:/2", "divide each die by 2"},
		{"3d6:x2", "multiply each die by 2"},
		{"3d6:*3", "multiply each die by 3 (alternate syntax)"},
		// {"4d6:dl1+2", "complex: drop lowest then add 2 to sum"},
		{"4d6:dl1:2e", "multiple complex modifiers"},
		{"1d6", "bare die with implicit count 1"},
		{"d6", "implicit count 1 with d prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			err := dice.ValidateExpression(tt.expr)
			if err != nil {
				t.Errorf("ValidateExpression(%q) = %v, want nil", tt.expr, err)
			}
		})
	}
}

func TestParseInvalidExpressions(t *testing.T) {
	tests := []struct {
		expr      string
		errSubstr string
	}{
		{"", "empty expression"},
		{"abc", "invalid dice type"},
		{"d", "missing sides"},
		{"1d", "missing sides"},
		{"1d0", "invalid sides"},
		{"3d6:invalid", "unknown complex modifier"},
		{"3d6:dl0", "invalid dl count"},
		{"3d6:dh-1", "invalid dh count"},
		{"3d6:/0", "invalid divisor: 0"},
		// {"3d6:dl5", "cannot drop 5 dice from pool of 3"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			err := dice.ValidateExpression(tt.expr)
			if err == nil {
				t.Errorf("ValidateExpression(%q) succeeded, want error containing %q", tt.expr, tt.errSubstr)
				return
			}
			if !strings.Contains(err.Error(), tt.errSubstr) {
				t.Errorf("ValidateExpression(%q) error = %q, want substring %q", tt.expr, err.Error(), tt.errSubstr)
			}
		})
	}
}

// ----------------------------------------------------------------------
// Rolling Tests with Deterministic Seeds

func TestRollDeterministic(t *testing.T) {
	m := newSeededManager(t, "testseed123")

	// With a fixed seed, rolls should be repeatable.
	roll1, err := m.Roll("3d6")
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}
	roll2, err := m.Roll("3d6")
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}
	if roll1 == roll2 {
		// Not an error, but unlikely; we just note the seed is fixed.
		t.Logf("Rolls with same seed: %d, %d", roll1, roll2)
	}

	// Create a new manager with the same seed; first roll should match.
	m2 := newSeededManager(t, "testseed123")
	roll3, err := m2.Roll("3d6")
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}
	if roll1 != roll3 {
		t.Errorf("Deterministic rolls mismatch: first manager gave %d, second gave %d", roll1, roll3)
	}
}

func TestRollWithModifiers(t *testing.T) {
	m := newSeededManager(t, "modtest")

	tests := []struct {
		expr string
		min  int
		max  int
	}{
		{"1d6+10", 11, 16},
		{"1d6-1", 0, 5},
		{"2d6:2e", 6, 16},  // 2 dice, each +2 => min 2*(1+2)=6, max 2*(6+2)=16
		{"4d6:dl1", 3, 18}, // drop lowest of 4d6 => sum of top 3
		{"3d6:/2", 0, 9},   // each die halved (integer division)
		{"3d6:x2", 6, 36},  // each die doubled
		{"1d2:5>>1", 6, 7},
		{"3d6:1>>1", 4, 19},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			// Roll many times to check bounds.
			for i := 0; i < 100; i++ {
				val, err := m.Roll(tt.expr)
				if err != nil {
					t.Fatalf("Roll(%q) error: %v", tt.expr, err)
				}
				if val < tt.min || val > tt.max {
					t.Errorf("Roll(%q) = %d, want in [%d,%d]", tt.expr, val, tt.min, tt.max)
				}
			}
		})
	}
}

func TestMustRollPanicsOnInvalid(t *testing.T) {
	m := newSeededManager(t, "panicseed")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustRoll did not panic on invalid expression")
		}
	}()
	_ = m.MustRoll("invalid")
}

func TestMustRollReturnsSum(t *testing.T) {
	m := newSeededManager(t, "mustroll")
	// With seed fixed, we know the result.
	val := m.MustRoll("2d6")
	if val < 2 || val > 12 {
		t.Errorf("MustRoll returned %d, expected between 2 and 12", val)
	}
}

func TestRollWithInlineMods(t *testing.T) {
	m := newSeededManager(t, "inlinemods")

	base, err := m.Roll("1d6")
	if err != nil {
		t.Fatal(err)
	}
	m = newSeededManager(t, "inlinemods")
	withMod, err := m.Roll("1d6", 5)
	if err != nil {
		t.Fatal(err)
	}
	if withMod != base+5 {
		t.Errorf("Roll with inline mod: base=%d, base+5=%d, got %d", base, base+5, withMod)
	}
}

// ----------------------------------------------------------------------
// D66, Flux, FluxGood, FluxBad Tests

func TestD66(t *testing.T) {
	m := newSeededManager(t, "d66seed")

	// Without modifiers, result should be two digits 0-9.
	result := m.D66()
	if len(result) != 2 {
		t.Errorf("D66() = %q, want length 2", result)
	}
	for _, ch := range result {
		if ch < '0' || ch > '9' {
			t.Errorf("D66() contains non-digit character: %q", result)
		}
	}

	// With modifiers, bounds should still hold.
	resultMod := m.D66(3, -2) // first die +3, second die -2
	if len(resultMod) != 2 {
		t.Errorf("D66() with mods = %q, want length 2", resultMod)
	}
	// Digits still 0-9 due to clamping.
}

func TestFluxFunctions(t *testing.T) {
	m := newSeededManager(t, "fluxseed")

	flux := m.Flux()
	if flux < -5 || flux > 5 {
		t.Errorf("Flux() = %d, want in [-5,5]", flux)
	}

	fluxGood := m.FluxGood()
	if fluxGood < 0 || fluxGood > 5 {
		t.Errorf("FluxGood() = %d, want in [0,5]", fluxGood)
	}

	fluxBad := m.FluxBad()
	if fluxBad < -5 || fluxBad > 0 {
		t.Errorf("FluxBad() = %d, want in [-5,0]", fluxBad)
	}

	// Test with modifiers.
	fluxMod := m.Flux(10)
	if fluxMod < 5 || fluxMod > 15 {
		t.Errorf("Flux(10) = %d, want in [5,15]", fluxMod)
	}
}

func TestVariance(t *testing.T) {
	m := newSeededManager(t, "varseed")

	for i := 0; i < 100; i++ {
		v := m.Variance()
		if v < 0.0 || v > 1.0 {
			t.Errorf("Variance() = %f, want in [0.0, 1.0]", v)
		}
		// Check it's a multiple of 0.001 (since 1d1001-1 gives 0-1000)
		scaled := int(math.Round(v * 1000))
		if float64(scaled)/1000 != v {
			t.Errorf("Variance() = %f, not a multiple of 0.001", v)
		}
	}
}

// ----------------------------------------------------------------------
// Concurrency Tests

func TestConcurrentRolls(t *testing.T) {
	m := newSeededManager(t, "concurrent")

	var wg sync.WaitGroup
	workers := 50
	rollsPerWorker := 20
	errCh := make(chan error, workers*rollsPerWorker)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < rollsPerWorker; j++ {
				_, err := m.Roll("3d6+2:dl1")
				if err != nil {
					errCh <- fmt.Errorf("worker %d roll %d: %v", id, j, err)
				}
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}
}

func TestConcurrentDefaultManager(t *testing.T) {
	// Default manager is safe for concurrent use.
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = dice.MustRoll("2d6")
			_ = dice.D66()
		}()
	}
	wg.Wait()
}

// ----------------------------------------------------------------------
// Caching Tests

func TestExpressionCache(t *testing.T) {
	m := newSeededManager(t, "cacheseed")

	expr := "3d6+5"
	// First roll parses and caches.
	_, err := m.Roll(expr)
	if err != nil {
		t.Fatal(err)
	}

	// We can't inspect cache directly, but we can verify that invalidating
	// something doesn't happen. A second roll should use cache.
	_, err = m.Roll(expr)
	if err != nil {
		t.Fatal(err)
	}
	// No direct assertion, but if cache were broken we'd see parsing errors.
}

// ----------------------------------------------------------------------
// Result Method Tests

func TestResult(t *testing.T) {
	m := newSeededManager(t, "resultseed")
	_, err := m.Roll("2d6")
	if err != nil {
		t.Fatal(err)
	}

	res := m.Result()
	dice := res.Dice()
	raw := res.Raw()

	if len(dice) != 2 {
		t.Errorf("Result.Dice() length = %d, want 2", len(dice))
	}
	if len(raw) != 2 {
		t.Errorf("Result.Raw() length = %d, want 2", len(raw))
	}
	for _, v := range raw {
		if v < 1 || v > 6 {
			t.Errorf("raw roll value %d out of range [1,6]", v)
		}
	}

	// Modify the returned slice; should not affect internal state.
	if len(raw) > 0 {
		raw[0] = 999
		newRaw := res.Raw()
		if newRaw[0] == 999 {
			t.Errorf("Result.Raw() returned slice that shares underlying array")
		}
	}
}

// ----------------------------------------------------------------------
// Edge Cases and Boundaries

func TestEdgeCases(t *testing.T) {
	m := newSeededManager(t, "edgecases")

	// Expression with 1 die and drop lowest (should work, drop 0)
	val, err := m.Roll("1d6:dl1")
	if err == nil {
		// Should error because cannot drop 1 from 1 die (must leave at least one)
		t.Errorf("Roll 1d6:dl1 should fail, got %d", val)
	}

	// Drop all but one is okay.
	val, err = m.Roll("4d6:dl3")
	if err != nil {
		t.Errorf("Roll 4d6:dl3 failed: %v", err)
	}
	if val < 1 || val > 6 {
		t.Errorf("dl3 from 4d6 should leave one die, got %d", val)
	}

	// Very large number of dice.
	_, err = m.Roll("100d100")
	if err != nil {
		t.Errorf("Roll 100d100 failed: %v", err)
	}

	// Negative modifiers in expression.
	val, err = m.Roll("1d6-10")
	if err != nil {
		t.Fatal(err)
	}
	// Result may be negative.
	if val > -4 || val < -9 {
		t.Errorf("1d6-10 = %d, expected between -9 and -4", val)
	}
}

// ----------------------------------------------------------------------
// Benchmark for performance measurement (optional)

func BenchmarkRoll(b *testing.B) {
	m, _ := dice.New("benchmark")
	expr := "4d6:dl1+2"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Roll(expr)
	}
}

func BenchmarkMustRoll(b *testing.B) {
	m, _ := dice.New("benchmark")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.MustRoll("3d6")
	}
}
