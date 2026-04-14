package dice

import "fmt"

// Roll return sum of dices. It panic if any error occur.
func (m *Manager) Roll(expr string, dm ...int) int {
	if err := m.RollSafe(expr); err != nil {
		panic(err)
	}
	sum, err := m.rollState.Sum()
	if err != nil {
		panic(err)
	}
	for _, val := range dm {
		sum += val
	}
	return sum
}

// Roll uses default roll manager and returns sum of dices. It panic if any error occurs.
func Roll(expr string, dm ...int) int {
	return defaultManager.Roll(expr, dm...)
}

func (m *Manager) D66(mods ...int) string {
	r1 := m.Roll("1d6")
	r2 := m.Roll("1d6")
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

func D66(mods ...int) string {
	return defaultManager.D66(mods...)
}

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
