package float

import "math"

//AI generated
func Round(fl float64) float64 {
	// Сначала округляем до миллиардных (9 знаков)
	intermediate := math.Round(fl*1e9) / 1e9

	n := 9
	if fl > 0.000001 {
		n = 6
	}
	if fl > 0.001 {
		n = 3
	}
	if fl > 10 {
		n = 1
	}

	// Затем округляем до требуемого количества знаков
	pow := math.Pow(10, float64(n))
	return math.Round(intermediate*pow) / pow
}

func RoundN(fl float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(fl*pow) / pow
}
