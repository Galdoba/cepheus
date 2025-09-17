package star

import (
	"math"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func originalMass(dp *dice.Dicepool, currentMass float64) float64 {
	return float64(dp.Sum("1d3")) / 2 * currentMass
}

func MainSequanceLifespan(mass float64) float64 {
	denominator := math.Pow(mass, 2.5)
	msa := 10.0 / denominator
	return msa
}

func smallStarAge(dp *dice.Dicepool) float64 {
	age := float64(dp.Sum("1d6"))*2 + float64(dp.Sum("1d3")-2) + variance(dp)
	age = adjust(dp, age)
	return roundFloat(age)
}

func largeStarAge(dp *dice.Dicepool, msa float64) float64 {
	age := msa * variance(dp)
	age = adjust(dp, age)
	return roundFloat(age)
}

func variance(dp *dice.Dicepool) float64 {
	return float64(dp.Sum("1d100")) / 100.0
}

func roundFloat(x float64) float64 {
	// Сначала округляем до миллионных (6 знаков)
	intermediate := math.Round(x*1e6) / 1e6

	n := 6
	if x > 0.001 {
		n = 3
	}
	if x > 10 {
		n = 1
	}

	// Затем округляем до требуемого количества знаков
	pow := math.Pow(10, float64(n))
	return math.Round(intermediate*pow) / pow
}
