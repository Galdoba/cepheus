package interpolate

import "math"

// ai generated
func interpolate(massMap map[int]float64, index int) float64 {
	// Если значение уже есть в карте, возвращаем его
	if mass, exists := massMap[index]; exists {
		return mass
	}
	// Находим ближайшие известные индексы для интерполяции
	lowerIndex, upperIndex := findClosestIndices(massMap, index)
	// Если не найдены подходящие индексы, возвращаем 0
	if lowerIndex == -1 || upperIndex == -1 {
		return 0
	}
	// Линейная интерполяция
	lowerMass := massMap[lowerIndex]
	upperMass := massMap[upperIndex]
	return lowerMass + (float64(index-lowerIndex))*(upperMass-lowerMass)/float64(upperIndex-lowerIndex)
}

// ai generated
func findClosestIndices(massMap map[int]float64, target int) (int, int) {
	var lower, upper int = -1, -1

	// Ищем ближайший меньший индекс
	for i := target - 1; i >= 10; i-- {
		if _, exists := massMap[i]; exists {
			lower = i
			break
		}
	}

	// Ищем ближайший больший индекс
	for i := target + 1; i <= 939; i++ {
		if _, exists := massMap[i]; exists {
			upper = i
			break
		}
	}

	return lower, upper
}

//ai generated
var WhiteDwarfCooling = map[float64]float64{
	0.0:  100000,
	0.1:  25000,
	0.5:  10000,
	1.0:  8000,
	1.5:  7000,
	2.5:  5500,
	5.0:  5000,
	10.0: 4000,
	13.0: 3800,
}

//ai generated
// WhiteDwarfTemp интерполирует температуру белого карлика для заданного возраста
func WhiteDwarfTemp(age float64) float64 {
	// Если возраст выходит за границы, возвращаем крайние значения
	if age <= 0.0 {
		return WhiteDwarfCooling[0.0]
	}
	if age >= 13.0 {
		return WhiteDwarfCooling[13.0]
	}

	// Если значение уже известно, возвращаем его
	if temp, exists := WhiteDwarfCooling[age]; exists {
		return temp
	}

	// Находим ближайшие известные возрасты для интерполяции
	var lowerAge, upperAge float64
	minDiffLower := math.MaxFloat64
	minDiffUpper := math.MaxFloat64

	for knownAge := range WhiteDwarfCooling {
		if knownAge < age && (age-knownAge) < minDiffLower {
			lowerAge = knownAge
			minDiffLower = age - knownAge
		}
		if knownAge > age && (knownAge-age) < minDiffUpper {
			upperAge = knownAge
			minDiffUpper = knownAge - age
		}
	}

	// Линейная интерполяция в логарифмическом масштабе (так как температура меняется нелинейно)
	// Используем логарифмы для более точной интерполяции
	logLowerTemp := math.Log(WhiteDwarfCooling[lowerAge])
	logUpperTemp := math.Log(WhiteDwarfCooling[upperAge])

	// Интерполируем в логарифмическом пространстве
	ratio := (age - lowerAge) / (upperAge - lowerAge)
	logInterpolated := logLowerTemp + ratio*(logUpperTemp-logLowerTemp)

	// Возвращаемся из логарифмического пространства
	return math.Exp(logInterpolated)
}
