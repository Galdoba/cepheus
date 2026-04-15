package dice

import "fmt"

// Roll evaluates the expression, applies the modifiers, and returns the sum.
// It is safe for concurrent use.
func (m *Manager) Roll(expr string, mods ...int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.roll(expr); err != nil {
		return 0, fmt.Errorf("roll %q %v failed: %w", expr, mods, err)
	}
	intr := &stdInterpreter{}
	intrpr, err := intr.interpret(m.rollState)
	if err != nil {
		return 0, fmt.Errorf("roll state recovery %q %v failed: %w", expr, mods, err)
	}
	sum := intrpr.sum
	for _, dm := range mods {
		sum += dm
	}
	return sum, nil
}

// Roll uses the default manager to evaluate an expression.
func Roll(expr string, mods ...int) (int, error) {
	return defaultManager.Roll(expr, mods...)
}

// MustRoll is like Roll but panics on error.
func (m *Manager) MustRoll(expr string, dm ...int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.roll(expr); err != nil {
		panic(err)
	}
	intr := &stdInterpreter{}
	intrpr, err := intr.interpret(m.rollState)
	if err != nil {
		panic(err)
	}
	sum := intrpr.sum
	for _, val := range dm {
		sum += val
	}
	return sum
}

// MustRoll uses the default manager and panics on error.
func MustRoll(expr string, dm ...int) int {
	return defaultManager.MustRoll(expr, dm...)
}

// D66 rolls two six-sided dice and returns a two-digit string (each digit 0‑9).
// Optional modifiers are applied to the first and second die respectively.
func (m *Manager) D66(mods ...int) string {
	r1 := m.MustRoll("1d6")
	r2 := m.MustRoll("1d6")
	for i, mod := range mods {
		switch i {
		case 0:
			r1 += mod
		case 1:
			r2 += mod
		}
	}
	r1 = setBounds(r1, 0, 9)
	r2 = setBounds(r2, 0, 9)
	return fmt.Sprintf("%d%d", r1, r2)
}

// D66 uses the default manager.
func D66(mods ...int) string {
	return defaultManager.D66(mods...)
}

// Flux returns (first die minus second die) plus any modifiers.
func (m *Manager) Flux(dm ...int) int {
	r1 := m.MustRoll("1d6")
	r2 := m.MustRoll("1d6")
	return r1 - r2 + sum(dm...)
}

// Flux uses the default manager.
func Flux(dm ...int) int {
	return defaultManager.Flux(dm...)
}

// FluxGood returns (higher die minus lower die) plus modifiers.
func (m *Manager) FluxGood(dm ...int) int {
	r1 := m.MustRoll("1d6")
	r2 := m.MustRoll("1d6")
	return max(r1, r2) - min(r1, r2) + sum(dm...)
}

// FluxGood uses the default manager.
func FluxGood(dm ...int) int {
	return defaultManager.FluxGood(dm...)
}

// FluxBad returns (lower die minus higher die) plus modifiers.
func (m *Manager) FluxBad(dm ...int) int {
	r1 := m.MustRoll("1d6")
	r2 := m.MustRoll("1d6")
	return min(r1, r2) - max(r1, r2) + sum(dm...)
}

// FluxBad uses the default manager.
func FluxBad(dm ...int) int {
	return defaultManager.FluxBad(dm...)
}

// Variance returns a random float64 in the range [0.0, 1.0] using the same random source.
func (m *Manager) Variance() float64 {
	r := m.MustRoll("1d1001-1")
	return float64(r) / 1000
}

// Variance uses the default manager.
func Variance() float64 {
	return defaultManager.Variance()
}

// sum is a helper that adds all integers in a slice.
func sum(values ...int) int {
	s := 0
	for _, v := range values {
		s += v
	}
	return s
}

// setBounds clamps val between min and max (order‑insensitive).
func setBounds(val, min, max int) int {
	if min > max {
		min, max = max, min
	}
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}