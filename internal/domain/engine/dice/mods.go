package dice

import "slices"

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
	Apply([]int) []int
	Priority() int
}

type None struct{}

func (m None) Apply(raw []int) []int {
	return raw
}

func (m None) Priority() int {
	return 0
}

type Sum struct{}

func (m Sum) Priority() int { return PrioritySum }

func (m Sum) Apply(raw []int) []int {
	sum := 0
	for _, v := range raw {
		sum += v
	}
	return []int{sum}
}

type AddConst struct {
	value int
}

func (m AddConst) Priority() int { return PriorityAddToSum }

func (m AddConst) Apply(raw []int) []int {
	out := []int{}
	for _, v := range raw {
		out = append(out, v+m.value)
	}
	return out
}

type AddToEach struct {
	value int
}

func (m AddToEach) Priority() int { return PriorityAddToEach }

func (m AddToEach) Apply(raw []int) []int {
	out := []int{}
	for _, v := range raw {
		out = append(out, v+m.value)
	}
	return out
}

type AddIndividual struct {
	position int
	value    int
}

func (m AddIndividual) Priority() int { return PriorityAddIndividual }

func (m AddIndividual) Apply(raw []int) []int {
	out := []int{}
	for k, v := range raw {
		switch k == m.position {
		case false:
			out = append(out, v)
		case true:
			out = append(out, v+m.value)
		}
	}
	return out
}

type DropLowest struct {
	quantity int
}

func (m DropLowest) Priority() int { return PriorityDropLowest }

func (m DropLowest) Apply(raw []int) []int {
	drop := min(m.quantity, len(raw)-1)
	if drop <= 0 {
		return slices.Clone(raw)
	}
	out := slices.Clone(raw)
	slices.Sort(out)
	return out[drop:]
}

type DropHighest struct {
	quantity int
}

func (m DropHighest) Priority() int { return PriorityDropHighest }

func (m DropHighest) Apply(raw []int) []int {
	drop := min(m.quantity, len(raw)-1)
	if drop <= 0 {
		return slices.Clone(raw)
	}
	out := slices.Clone(raw)
	slices.Sort(out)
	return out[:len(out)-drop]
}

type Divide struct {
	value int
}

func (m Divide) Priority() int { return PriorityDivide }

func (m Divide) Apply(raw []int) []int {
	out := []int{}
	for _, v := range raw {
		out = append(out, v/m.value)
	}
	return out
}

type Multiply struct {
	value int
}

func (m Multiply) Priority() int { return PriorityMultiply }

func (m Multiply) Apply(raw []int) []int {
	out := []int{}
	for _, v := range raw {
		out = append(out, v*m.value)
	}
	return out
}
