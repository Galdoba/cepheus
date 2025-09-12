package star

import "math"

func (st Star) index() int {
	index := 0
	switch st.Class {
	case "Ia":
		index += 100
	case "Ib":
		index += 200
	case "II":
		index += 300
	case "III":
		index += 400
	case "IV":
		index += 500
	case "V":
		index += 600
	case "VI":
		index += 700
	case "VII", "D":
		index += 800
	case "BD":
		index += 900
	}

	switch st.Type {
	case "O", "L":
		index += 10
	case "B", "T":
		index += 20
	case "A", "Y":
		index += 30
	case "F":
		index += 40
	case "G":
		index += 50
	case "K":
		index += 60
	case "M":
		index += 70
	}
	if st.SubType != nil {
		index += *st.SubType
	}
	return index
}

func fromIndex(i int) (string, string, int) {
	cl := ""
	switch i / 100 {
	case 1:
		cl = "Ia"
	case 2:
		cl = "Ib"
	case 3:
		cl = "II"
	case 4:
		cl = "III"
	case 5:
		cl = "IV"
	case 6:
		cl = "V"
	case 7:
		cl = "VI"
	case 8:
		cl = "D"
	case 9:
		cl = "BD"
	}
	noclass := i % 100
	scl := ""
	switch noclass / 10 {
	case 1:
		scl = "O"
		if cl == "BD" {
			scl = "L"
		}
	case 2:
		scl = "B"
		if cl == "BD" {
			scl = "T"
		}
	case 3:
		scl = "A"
		if cl == "BD" {
			scl = "Y"
		}
	case 4:
		scl = "F"
	case 5:
		scl = "G"
	case 6:
		scl = "K"
	case 7:
		scl = "M"
	}
	sub := noclass % 10
	return cl, scl, sub
}

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
// interpolateWhiteDwarfTemp интерполирует температуру белого карлика для заданного возраста
func interpolateWhiteDwarfTemp(age float64) float64 {
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
