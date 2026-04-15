package dice

import (
	"fmt"
	"slices"
)

const (
	modNone          = "none"
	modAddToEach     = "add to each"
	modAddIndividual = "add individual"
	modAddToSum      = "add to sum"
	modDropLowest    = "drop lowest"
	modDropHighest   = "drop highest"
	modDivide        = "divide"
	modMultiply      = "multiply"
	modSum           = "sum"

	priorityNone          = 0
	priorityAddIndividual = 20
	priorityAddToEach     = 30
	priorityDropLowest    = 70
	priorityDropHighest   = 71
	prioritySum           = 100
	priorityAddToSum      = 110
	priorityDivide        = 120
	priorityMultiply      = 130
)

type mod interface {
	apply([]int) ([]int, error)
	priority() int
}

type none struct{}

func (m none) apply(raw []int) ([]int, error) { return raw, nil }
func (m none) priority() int                  { return priorityNone }

type summ struct{}

func (m summ) priority() int { return prioritySum }
func (m summ) apply(raw []int) ([]int, error) {
	s := 0
	for _, v := range raw {
		s += v
	}
	return []int{s}, nil
}

type addConst struct{ value int }

func (m addConst) priority() int { return priorityAddToSum }
func (m addConst) apply(raw []int) ([]int, error) {
	out := make([]int, len(raw))
	for i, v := range raw {
		out[i] = v + m.value
	}
	return out, nil
}

type addToEach struct{ value int }

func (m addToEach) priority() int { return priorityAddToEach }
func (m addToEach) apply(raw []int) ([]int, error) {
	out := make([]int, len(raw))
	for i, v := range raw {
		out[i] = v + m.value
	}
	return out, nil
}

type addIndividual struct {
	position int
	value    int
}

func (m addIndividual) priority() int { return priorityAddIndividual }
func (m addIndividual) apply(raw []int) ([]int, error) {
	if m.position < 1 || m.position > len(raw) {
		return nil, fmt.Errorf("position %d out of range (1..%d)", m.position, len(raw))
	}
	out := make([]int, len(raw))
	for i, v := range raw {
		if i == m.position-1 {
			out[i] = v + m.value
		} else {
			out[i] = v
		}
	}
	return out, nil
}

type dropLowest struct{ quantity int }

func (m dropLowest) priority() int { return priorityDropLowest }
func (m dropLowest) apply(raw []int) ([]int, error) {
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

type dropHighest struct{ quantity int }

func (m dropHighest) priority() int { return priorityDropHighest }
func (m dropHighest) apply(raw []int) ([]int, error) {
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

type divide struct{ value int }

func (m divide) priority() int { return priorityDivide }
func (m divide) apply(raw []int) ([]int, error) {
	if m.value == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	out := make([]int, len(raw))
	for i, v := range raw {
		out[i] = v / m.value
	}
	return out, nil
}

type multiply struct{ value int }

func (m multiply) priority() int { return priorityMultiply }
func (m multiply) apply(raw []int) ([]int, error) {
	out := make([]int, len(raw))
	for i, v := range raw {
		out[i] = v * m.value
	}
	return out, nil
}
