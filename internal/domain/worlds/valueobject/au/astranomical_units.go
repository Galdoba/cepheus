package au

import "github.com/Galdoba/cepheus/pkg/float"

const (
	AU_to_Km = 149597870.9
)

type AU float64

type distanceValue struct {
	distance   float64
	difference float64
}

var distanceTable = map[float64]distanceValue{
	0:  distanceValue{0, 0.4},
	1:  distanceValue{0.4, 0.3},
	2:  distanceValue{0.7, 0.4},
	3:  distanceValue{1.0, 0.6},
	4:  distanceValue{1.6, 1.2},
	5:  distanceValue{2.8, 2.4},
	6:  distanceValue{5.2, 4.8},
	7:  distanceValue{10, 10},
	8:  distanceValue{20, 20},
	9:  distanceValue{40, 37},
	10: distanceValue{77, 77},
	11: distanceValue{154, 154},
	12: distanceValue{308, 307},
	13: distanceValue{615, 615},
	14: distanceValue{1230, 1270},
	15: distanceValue{2500, 2400},
	16: distanceValue{4900, 4900},
	17: distanceValue{9800, 9700},
	18: distanceValue{19500, 20000},
	19: distanceValue{39500, 39200},
	20: distanceValue{78700, 0},
}

func FromOrbitNumber(orbitN float64) AU {
	integer, fraction := float.SplitFloatAbs(orbitN)
	units := distanceTable[integer].distance + (distanceTable[integer].difference * fraction)
	units = float.RoundN(units, 6)
	if units > 0.001 {
		units = float.RoundN(units, 3)
	}
	if units > 0.01 {
		units = float.RoundN(units, 2)
	}
	if units > 0.1 {
		units = float.RoundN(units, 1)
	}
	return AU(units)
}

func (au AU) OrbitNumber() float64 {
	auValue := float64(au)
	fullOrbit := 0.0
	for orbitNum := 20.0; orbitNum >= 0; orbitNum-- {
		distance := distanceTable[orbitNum].distance
		if auValue >= distance {
			fullOrbit = orbitNum
			break
		}
	}

	dist := distanceTable[fullOrbit].distance
	diff := distanceTable[fullOrbit].difference

	// Если разница равна 0 (последняя орбита), возвращаем целое число
	if diff == 0 {
		return fullOrbit
	}

	remainder := auValue - dist
	fraction := remainder / diff

	orbitN := float.RoundN(fullOrbit+fraction, 3)
	if orbitN > 0.01 {
		orbitN = float.RoundN(orbitN, 2)
	}
	if orbitN > 10 {
		orbitN = float.RoundN(orbitN, 1)
	}
	return orbitN
}

func (au AU) DistanceKM() int {
	dist := float64(au) * AU_to_Km
	return int(dist)
}

func (au AU) Distance() float64 {
	return float.RoundN(float64(au)*AU_to_Km/1000.0, 3)
}
