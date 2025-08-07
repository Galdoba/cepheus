package cepheus

import (
	"github.com/Galdoba/cepheus/dice"
	"github.com/Galdoba/cepheus/ehex"
	"github.com/Galdoba/cepheus/uwp"
)

const (
	dices = "2d6"
)

func Generate(dp *dice.Dicepool) map[string]int {
	data := make(map[string]int)
	data[uwp.Size] = dp.Sum(dices) - 2

	data[uwp.Atmo] = dp.Sum(dices) - 7 + data[uwp.Size] + modIfLower(-99, data[uwp.Size], 1)
	data[uwp.Atmo] = bound(data[uwp.Atmo], 0, 15)

	data[uwp.Hydr] = dp.Sum(dices) - 7 + data[uwp.Size] +
		modIfEqual(-4, data[uwp.Atmo], 0, 1, 10, 11, 12) +
		modIfEqual(-2, data[uwp.Atmo], 14) +
		modIfLower(-99, data[uwp.Size], 2)
	data[uwp.Hydr] = bound(data[uwp.Hydr], 0, 10)

	data[uwp.Pops] = dp.Sum(dices) - 2 +
		modIfEqual(-2, data[uwp.Atmo], 0, 1, 2, 3, 4, 7, 9, 10, 11, 12, 13, 14, 15) +
		modIfEqual(3, data[uwp.Atmo], 6) +
		modIfEqual(1, data[uwp.Atmo], 7) +
		modIfEqual(1, data[uwp.Atmo], 5, 8)
	if data[uwp.Hydr] == 0 {
		data[uwp.Pops] += modIfEqual(-1, data[uwp.Atmo], 0, 1, 2, 3, 14)
	}
	data[uwp.Pops] = bound(data[uwp.Pops], 0, 10)

	data[uwp.Govr] = dp.Sum(dices) - 7 + data[uwp.Pops] + modIfLower(-99, data[uwp.Pops], 1)
	data[uwp.Govr] = bound(data[uwp.Govr], 0, 15)

	data[uwp.Laws] = dp.Sum(dices) - 7 + data[uwp.Govr] + modIfLower(-99, data[uwp.Govr], 1)
	data[uwp.Laws] = bound(data[uwp.Laws], 0, 10)

	portRoll := dp.Sum("2d6") + 7 - data[uwp.Pops]
	if portRoll <= 2 {
		data[uwp.Port] = ehex.FromString("X").Value()
	}
	return data
}

func bound(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func modIfLower(mod, value, target int) int {
	if value < target {
		return mod
	}
	return 0
}

func modIfEqual(mod, value int, target ...int) int {
	for _, val := range target {
		if val == value {
			return mod
		}
	}
	return 0
}
