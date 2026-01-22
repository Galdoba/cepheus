package dice

import (
	"fmt"
	"testing"
)

// TestDiceExpression_Comprehensive runs a comprehensive test suite for dice expression parsing.
// It tests both valid expressions and error cases for sum-based rolls.
func TestDiceExpression_Comprehensive(t *testing.T) {
	tests := []struct {
		name        string // Descriptive test name
		expr        string // Dice expression to test
		expectError bool   // Whether the expression should produce an error
	}{
		// Sum tests (1-60)
		{"sum basic", "2d6", false},
		{"sum big numbers", "10d10", false},
		{"sum minimal", "1d1", false},
		{"sum additive positive", "3d6+2", false},
		{"sum additive negative", "3d6-2", false},
		{"sum multiple additive", "3d6+2-1+3", false},
		{"sum additive zero positive", "3d6+0", false},
		{"sum additive zero negative", "3d6-0", false},
		{"sum large additive", "3d6+100", false},
		{"sum large negative additive", "3d6-100", false},
		{"sum multiplicative", "2d6x2", false},
		{"sum multiple multiplicative", "2d6x2x3", false},
		{"sum negative multiplicative", "2d6x-2", false},
		{"sum zero multiplicative", "2d6x0", false},
		{"sum multiplicative with plus", "2d6x+2", false},
		{"sum deletive", "2d6/3", false},
		{"sum multiple deletive", "2d6/2/3", false},
		{"sum negative deletive", "2d6/-2", false},
		{"sum replace simple", "2d6r1:2", false},
		{"sum replace multiple sources", "2d6r1;2:3", false},
		{"sum multiple replaces", "2d6r1:2r3:4", false},
		{"sum replace three sources", "2d6r1;2;3:6", false},
		{"sum drop low", "4d6dl1", false},
		{"sum drop high", "4d6dh1", false},
		{"sum drop low and high", "4d6dl1dh1", false},
		{"sum drop 2 low", "4d6dl2", false},
		{"sum drop 2 high", "4d6dh2", false},
		{"sum individual positive", "3d6i+1", false},
		{"sum individual negative", "3d6i-1", false},
		{"sum multiple individual", "3d6i+1i-2", false},
		{"sum individual no sign", "3d6i1", false},
		{"sum three individuals", "3d6i+1i+2i-3", false},
		{"sum minimum", "2d6min5", false},
		{"sum maximum", "2d6max10", false},
		{"sum min and max", "2d6min5max10", false},
		{"sum negative minimum", "2d6min-5", false},
		{"sum negative maximum", "2d6max-1", false},
		{"sum minimum with plus", "2d6min+5", false},
		{"sum reroll simple", "2d6rr1", false},
		{"sum reroll multiple values", "2d6rr1;2", false},
		{"sum multiple reroll tokens", "2d6rr1rr2", false},
		{"sum reroll three values", "2d6rr1;2;3", false},
		{"sum reroll many values", "2d6rr1;2;3;4;5;6", false},
		{"sum complex", "4d6+2x2/2r1:6dl1dh1i+1min10max20rr1;2", false},
		{"sum another complex", "10d10+5x2/2r1;2:10dl2dh1i-2min0max100rr1;2;3", false},
		{"sum zero dice error", "0d6", true},
		{"sum zero faces error", "2d0", true},
		{"sum negative dice error", "-2d6", true},
		{"sum negative faces error", "2d-6", true},
		{"sum drop low too many error", "2d6dl5", true},
		{"sum drop high too many error", "2d6dh5", true},
		{"sum duplicate dl error", "2d6dl1dl1", true},
		{"sum duplicate dh error", "2d6dh1dh1", true},
		{"sum duplicate replace error", "2d6r1:2r1:3", true},
		{"sum duplicate rr tokens error", "2d6rr1rr1", true},
		{"sum duplicate rr values error", "2d6rr1;1", true},
		{"sum incomplete multiplicative error", "2d6x", true},
		{"sum incomplete additive error", "2d6+", true},
		{"sum unknown token error", "2d6unknown", true},

		// Note: Concat tests are commented out in this test suite
		// but would follow the same pattern as sum tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := DiceExpression(tt.expr).ParseRoll()

			if tt.expectError {
				if err == nil {
					t.Errorf("Parse() expected error for %q but got none", tt.expr)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error for %q: %v", tt.expr, err)
				return
			}

			// Print parsed directives for debugging/verification
			fmt.Println(tt.expr, rd)
		})
	}

	// Demonstrate that different seeds produce different results
	fmt.Println(Roll("2d6"))
	fmt.Println(Roll("2d6"))
	fmt.Println(Roll("2d6"))
}
