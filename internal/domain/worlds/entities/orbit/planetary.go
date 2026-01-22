package orbit

import (
	"math"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/au"
	"github.com/Galdoba/cepheus/pkg/dice"
)

type Planetary struct {
	Position float64 //0 - 20
	Period   float64 //days
}

type planetaryOrbitGenerator struct {
	rng              *dice.Roller
	forcedPosition   float64
	forcedPeriod     float64
	forcedCenterMass float64
}

func newPlanetaryOrbitGenerator() *planetaryOrbitGenerator {
	pog := planetaryOrbitGenerator{}
	pog.rng = dice.New("")
	pog.forcedPeriod = -1
	pog.forcedCenterMass = 1
	return &pog
}

type planetaryOrbitDetails func(*planetaryOrbitGenerator)

func New(details ...planetaryOrbitDetails) *Planetary {
	pog := newPlanetaryOrbitGenerator()
	for _, apply := range details {
		apply(pog)
	}
	o := Planetary{}
	o.Position = pog.orbitN()
	o.Period = pog.period(o.Position)

	return &o
}

func (pog *planetaryOrbitGenerator) orbitN() float64 {
	if pog.forcedPosition > 0 {
		return pog.forcedPosition
	}
	return float64(pog.rng.Roll("1d6")) + (float64(pog.rng.Flux()) * 0.1)
}

func (pog *planetaryOrbitGenerator) period(orbitN float64) float64 {
	if pog.forcedPeriod > 0 {
		return pog.forcedPeriod
	}
	dist := math.Pow(float64(au.FromOrbitNumber(orbitN)), 3.0)
	return float.RoundN(dist/pog.forcedCenterMass, 2)
}
