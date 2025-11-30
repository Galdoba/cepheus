package value

import (
	"github.com/Galdoba/cepheus/internal/domain/core/values/modifier"
)

const (
	defaultCharactristicLimit = 15
)

type CharacteristicValue struct {
	AdjustableValue
}

func NewCharacteristicValue() *CharacteristicValue {
	c := CharacteristicValue{}
	c.exist = true
	c.valueType = ValueTypeCharacteristic
	c.highLimit = defaultCharactristicLimit
	c.dmFunc = charactristicDM
	c.modifiers = make(map[string]modifier.Modifier)
	return &c
}

func charactristicDM(val int) int {
	if val < 0 {
		val = 0
	}
	switch val {
	case 0:
		return -3
	case 1, 2:
		return -2
	default:
		return (val / 3) - 2
	}
}

func (c *CharacteristicValue) Value() int {
	v := minmax(c.baseValue, 0, c.highLimit) + c.adjustment
	return max(v, 0)
}

func (c *CharacteristicValue) ValueModded(mods ...string) int {
	return max(c.Value()+c.SumModifiers(mods...), 0)
}

func (c *CharacteristicValue) Add(add int) {
	c.baseValue = minmax(c.baseValue+add, 0, c.highLimit)
}

func (c *CharacteristicValue) DM() int {
	return c.dmFunc(c.Value())
}

func (c *CharacteristicValue) DM_With(mods ...string) int {
	return c.dmFunc(c.ValueModded(mods...))
}
