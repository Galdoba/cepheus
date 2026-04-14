package dice

import (
	"fmt"
	"slices"
)

const (
	ModNone          = "none"
	ModAddToEach     = "add to each"
	ModAddIndividual = "add individual"
	ModAddToSum      = "add to sum"
	ModDropLowest    = "drop lowest"
	ModDropHighest   = "drop highest"
	ModDivide        = "divide"
	ModMultiply      = "multiply"
	ModSum           = "sum"

	PriorityNone          = 0
	PriorityAddIndividual = 20
	PriorityAddToEach     = 30
	PriorityDropLowest    = 70
	PriorityDropHighest   = 71
	PrioritySum           = 100
	PriorityAddToSum      = 110
	PriorityDivide        = 120
	PriorityMultiply      = 130
)

type Mod interface {
	Apply([]int) ([]int, error)
	Priority() int
}

type None struct{}

func (m None) Apply(raw []int) ([]int, error) {
	return raw, nil
}

func (m None) Priority() int {
	return 0
}

// Sum is aggegation mod used by normal dice
type Sum struct{}

func (m Sum) Priority() int { return PrioritySum }

func (m Sum) Apply(raw []int) ([]int, error) {
	sum := 0
	for _, v := range raw {
		sum += v
	}
	return []int{sum}, nil
}

// AddConst add value after Sum took place
type AddConst struct {
	value int
}

func (m AddConst) Priority() int { return PriorityAddToSum }

func (m AddConst) Apply(raw []int) ([]int, error) {
	out := []int{}
	for _, v := range raw {
		out = append(out, v+m.value)
	}
	return out, nil
}

// AddToEach add value before Sum took place
type AddToEach struct {
	value int
}

func (m AddToEach) Priority() int { return PriorityAddToEach }

func (m AddToEach) Apply(raw []int) ([]int, error) {
	out := []int{}
	for _, v := range raw {
		out = append(out, v+m.value)
	}
	return out, nil
}

// AddIndividual add value to individual dice
type AddIndividual struct {
	position int
	value    int
}

func (m AddIndividual) Priority() int { return PriorityAddIndividual }

func (m AddIndividual) Apply(raw []int) ([]int, error) {
	if m.position < 0 {
		return nil, fmt.Errorf("position index cannot be negative: %d", m.position)
	}
	if m.position >= len(raw) {
		return nil, fmt.Errorf("position %d out of range, have %d dice", m.position, len(raw))
	}
	out := []int{}
	for k, v := range raw {
		switch k == m.position {
		case false:
			out = append(out, v)
		case true:
			out = append(out, v+m.value)
		}
	}
	return out, nil
}

// DropLowest discards dice with lowest value (used for Boon rolls)
type DropLowest struct {
	quantity int
}

func (m DropLowest) Priority() int { return PriorityDropLowest }

func (m DropLowest) Apply(raw []int) ([]int, error) {
	if m.quantity < 0 {
		return nil, fmt.Errorf("drop quantity cannot be negative: %d", m.quantity)
	}
	if m.quantity > len(raw)-1 {
		return nil, fmt.Errorf("cannot drop %d dice from pool of %d", m.quantity, len(raw))
	}
	if m.quantity == 0 {
		return slices.Clone(raw), nil
	}
	out := slices.Clone(raw)
	slices.Sort(out)
	return out[m.quantity:], nil
}

// DropHighest discards dice with highest values (used for Bane rolls)
type DropHighest struct {
	quantity int
}

func (m DropHighest) Priority() int { return PriorityDropHighest }

func (m DropHighest) Apply(raw []int) ([]int, error) {
	if m.quantity < 0 {
		return nil, fmt.Errorf("drop quantity cannot be negative: %d", m.quantity)
	}
	if m.quantity > len(raw)-1 {
		return nil, fmt.Errorf("cannot drop %d dice from pool of %d", m.quantity, len(raw))
	}
	if m.quantity == 0 {
		return slices.Clone(raw), nil
	}
	out := slices.Clone(raw)
	slices.Sort(out)
	return out[:len(out)-m.quantity], nil
}

// Divide divides values after Sum
type Divide struct {
	value int
}

func (m Divide) Priority() int { return PriorityDivide }

func (m Divide) Apply(raw []int) ([]int, error) {
	if m.value == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	out := []int{}
	for _, v := range raw {
		out = append(out, v/m.value)
	}
	return out, nil
}

// Multiply multiplies after Sum
type Multiply struct {
	value int
}

func (m Multiply) Priority() int { return PriorityMultiply }

func (m Multiply) Apply(raw []int) ([]int, error) {
	out := []int{}
	for _, v := range raw {
		out = append(out, v*m.value)
	}
	return out, nil
}
