package value

import (
	"github.com/Galdoba/cepheus/internal/domain/core/values/modifier"
)

const (
	ValueTypeCharacteristic = "characteristic"
	ValueTypeSkill          = "skill"
	minAdj                  = -6
	maxAdj                  = 6
)

// AdjustableValue represents a numeric value that can be modified and adjusted within defined limits
type AdjustableValue struct {
	exist        bool
	ebsenceValue int
	baseValue    int
	modifiers    map[string]modifier.Modifier
	adjustment   int
	highLimit    int
	valueType    string
	dmFunc       func(int) int
}

// minmax constrains a value between minimum and maximum bounds
func minmax(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (av *AdjustableValue) SetBase(v int) {
	av.baseValue = minmax(v, 0, av.highLimit)
}

func (av *AdjustableValue) SetModifier(mod modifier.Modifier) {
	av.modifiers[mod.Category()] = mod
}

func (av *AdjustableValue) GetModifier(category string) modifier.Modifier {
	if m, ok := av.modifiers[category]; ok {
		return m
	}
	return modifier.Modifier{}
}

func (av *AdjustableValue) SumModifiers(category ...string) int {
	sum := 0
	for _, mod := range av.modifiers {
		sum += mod.Value()
	}
	return sum
}

func (av *AdjustableValue) ListModifiers() []modifier.Modifier {
	list := []modifier.Modifier{}
	for _, m := range av.modifiers {
		list = append(list, m)
	}
	return list
}

func (av *AdjustableValue) RemoveModifier(category string) {
	delete(av.modifiers, category)
}
