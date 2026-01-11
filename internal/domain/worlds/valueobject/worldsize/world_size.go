package worldsize

import (
	"fmt"
	"math"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/ehex"
	"github.com/Galdoba/cepheus/pkg/float"
)

const (
	ExoticIce       = "Exotic Ice"
	MostlyIce       = "Mostly Ice"
	MostlyRock      = "Mostly Rock"
	RockAndMetal    = "Rock and Metal"
	MostlyMetal     = "Mostly Metal"
	CompressedMetal = "Compressed Metal"
	diameterOfTerra = 12742
	massOfTerra     = 1.0
)

type WorldSize struct {
	Code           ehex.Ehex
	Diameter       int
	Composition    string
	Dencity        float64
	Gravity        float64
	Mass           float64
	EscapeVelocity float64
	ResourceRating float64
}

type DetailGenerator struct {
	rng               *dice.Dicepool
	hzco              float64
	systemAge         float64
	forcedSeed        string
	forcedDiameter    int
	forcedComposition string
	forcedDencity     float64
	forcedGravity     float64
	forcedMass        float64
}

func newDetailGenerator() *DetailGenerator {
	dg := DetailGenerator{}
	dg.rng = dice.NewDicepool()
	dg.forcedGravity = -1
	return &dg
}

func New(code ehex.Ehex, details ...SizeDetails) WorldSize {
	dg := newDetailGenerator()
	for _, apply := range details {
		apply(dg)
	}
	worldSize := WorldSize{}
	worldSize.Code = code
	if dg.forcedSeed != "" {
		dg.rng = dice.NewDicepool(dice.WithSeedString(dg.forcedSeed))
	}
	switch worldSize.Code.Code() {
	case "0":
		panic("belt is not implemented")
	default:
		worldSize.Diameter = dg.diameter(worldSize.Code)
		worldSize.Composition = dg.terrestrialComposition(worldSize.Code)
		worldSize.Dencity = float.RoundN(dg.terrestrialDencity(worldSize.Composition), 3)
		worldSize.Gravity = dg.gravity(worldSize.Dencity, worldSize.Diameter)
		worldSize.Mass = dg.mass(worldSize.Dencity, worldSize.Diameter)
		worldSize.EscapeVelocity = escapeVelocity(worldSize.Mass, worldSize.Diameter)
	}

	return worldSize
}

func (ws WorldSize) Size() int {
	switch ws.Code.Code() {
	case "R", "S":
		return 0
	default:
		return ws.Code.Value()
	}
}

func (ws WorldSize) Profile() string {
	return fmt.Sprintf("%v-%vkm-%v-%v-%v", ws.Code.Code(), ws.Diameter, ws.Dencity, ws.Gravity, ws.Mass)
}

func (dg *DetailGenerator) diameter(size ehex.Ehex) int {
	switch dg.forcedDiameter {
	default:
		return dg.forcedDiameter
	case 0:
		return baseDiameter(size) + diameterIncrease(dg.rng)
	}
}

func baseDiameter(code ehex.Ehex) int {
	switch code.Code() {
	case "0", "R":
		return 0
	case "S":
		return 600
	case "1":
		return 1600
	case "2":
		return 3200
	case "3":
		return 4800
	case "4":
		return 6400
	case "5":
		return 8000
	case "6":
		return 9600
	case "7":
		return 11200
	case "8":
		return 12800
	case "9":
		return 14400
	case "A":
		return 16000
	case "B":
		return 17600
	case "C":
		return 19200
	case "D":
		return 20800
	case "E":
		return 22400
	case "F":
		return 24000
	}
	return -1
}

func diameterIncrease(rng *dice.Dicepool) int {
	done := false
	increment := 0
	for !done {
		increment = 0
		r1 := rng.Sum("1d3")
		r2 := rng.Sum("1d6")
		switch r1 {
		case 1:
			increment = 0
		case 2:
			increment = 600
		case 3:
			increment = 1200
		}
		switch r2 {
		case 1:
			increment += 0
		case 2:
			increment += 100
		case 3:
			increment += 200
		case 4:
			increment += 300
		case 5:
			if r1 == 3 {
				continue
			}
			increment += 400
		case 6:
			if r1 == 3 {
				continue
			}
			increment += 500
		}
		increment += rng.Sum("1d100") - 1
		done = true
	}
	return increment
}

func (dg *DetailGenerator) terrestrialComposition(sizeCode ehex.Ehex) string {
	if dg.forcedComposition != "" {
		return dg.forcedComposition
	}
	dm := 0
	switch sizeCode.Code() {
	case "0", "1", "2", "3", "4", "R", "S":
		dm += -1
	case "6", "7", "8", "9":
		dm += 1
	case "A", "B", "C", "D", "E", "F":
		dm += 3
	}
	switch dg.hzco < 1 {
	case true:
		dm += 1
	case false:
		dm += -1 - int(dg.hzco)
	}
	if dg.systemAge > 10 {
		dm += 1
	}
	r := bound(dg.rng.Sum("2d6")+dm, -4, 15)
	switch r {
	case -4:
		return ExoticIce
	case -3, -2, -1, 0, 1, 2:
		return MostlyIce
	case 3, 4, 5, 6:
		return MostlyRock
	case 7, 8, 9, 10, 11:
		return RockAndMetal
	case 12, 13, 14:
		return MostlyMetal
	case 15:
		return CompressedMetal
	}
	return ""
}

func bound(i, min, max int) int {
	if min > max {
		min, max = max, min
	}
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}

func (dg *DetailGenerator) terrestrialDencity(composition string) float64 {
	if dg.forcedDencity != 0 {
		return dg.forcedDencity
	}
	array := []float64{}
	switch composition {
	case ExoticIce:
		array = dencityArray(0.03, 0.03)
	case MostlyIce:
		array = dencityArray(0.18, 0.03)
	case MostlyRock:
		array = dencityArray(0.5, 0.03)
	case RockAndMetal:
		array = dencityArray(0.82, 0.03)
	case MostlyMetal:
		array = dencityArray(1.15, 0.03)
	case CompressedMetal:
		array = dencityArray(1.5, 0.05)
	}
	return array[dg.rng.Sum("2d6")-2]
}

func dencityArray(start, step float64) []float64 {
	array := []float64{start}
	for len(array) < 11 {
		last := array[len(array)-1]
		array = append(array, last+step)
	}
	return array
}

func (dg *DetailGenerator) gravity(dencity float64, diameter int) float64 {
	if dg.forcedGravity > 0 {
		return dg.forcedGravity
	}
	return gravity(dencity, diameter)
}

func gravity(dencity float64, diameter int) float64 {
	gr := (dencity * float64(diameter)) / float64(diameterOfTerra)
	return float.RoundN(gr, 2)
}

func (dg *DetailGenerator) mass(dencity float64, diameter int) float64 {
	if dg.forcedMass > 0 {
		return dg.forcedMass
	}
	return mass(dencity, diameter)
}

func mass(dencity float64, diameter int) float64 {
	base := float64(diameter) / float64(diameterOfTerra)
	return float.RoundN(dencity*(base*base*base), 3)
}

func escapeVelocity(mass float64, diameter int) float64 {
	m := mass / massOfTerra
	d := float64(diameter) / float64(diameterOfTerra)
	velocity := math.Sqrt(m / d)
	return float.RoundN(velocity*11.186, 1)
}

func orbitalVelocity(escapeVelocity float64) float64 {
	return float.RoundN(escapeVelocity+math.Sqrt2, 3)
}
