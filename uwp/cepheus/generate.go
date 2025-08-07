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

	data[uwp.Atmosphere] = dp.Sum(dices) - 7 + data[uwp.Size] + modIfLower(-99, data[uwp.Size], 1)
	data[uwp.Atmosphere] = bound(data[uwp.Atmosphere], 0, 15)

	data[uwp.Hydrosphere] = dp.Sum(dices) - 7 + data[uwp.Size] +
		modIfEqual(-4, data[uwp.Atmosphere], 0, 1, 10, 11, 12) +
		modIfEqual(-2, data[uwp.Atmosphere], 14) +
		modIfLower(-99, data[uwp.Size], 2)
	data[uwp.Hydrosphere] = bound(data[uwp.Hydrosphere], 0, 10)

	data[uwp.Population] = dp.Sum(dices) - 2 +
		modIfEqual(-2, data[uwp.Atmosphere], 0, 1, 2, 3, 4, 7, 9, 10, 11, 12, 13, 14, 15) +
		modIfEqual(3, data[uwp.Atmosphere], 6) +
		modIfEqual(1, data[uwp.Atmosphere], 7) +
		modIfEqual(1, data[uwp.Atmosphere], 5, 8)
	if data[uwp.Hydrosphere] == 0 {
		data[uwp.Population] += modIfEqual(-1, data[uwp.Atmosphere], 0, 1, 2, 3, 14)
	}
	data[uwp.Population] = bound(data[uwp.Population], 0, 10)

	data[uwp.Government] = dp.Sum(dices) - 7 + data[uwp.Population] + modIfLower(-99, data[uwp.Population], 1)
	data[uwp.Government] = bound(data[uwp.Government], 0, 15)

	data[uwp.Laws] = dp.Sum(dices) - 7 + data[uwp.Government] + modIfLower(-99, data[uwp.Government], 1)
	data[uwp.Laws] = bound(data[uwp.Laws], 0, 10)

	portRoll := dp.Sum("2d6") + 7 - data[uwp.Population]
	if portRoll <= 2 {
		data[uwp.Port] = ehex.FromString("X").Value()
	}:q
	:

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
