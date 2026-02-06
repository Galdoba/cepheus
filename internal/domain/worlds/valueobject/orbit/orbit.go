package orbit

import (
	"math"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/au"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
)

const (
	DM_ObjectIsStar                 = 2
	DM_PerParentObject              = 1
	DM_PerOrbitsBelow1ForYangSystem = -1
	DM_ObjectIsBelt                 = 1
)

type Orbit struct {
	Distance          float64 //Diameter of parent body or OrbitN
	Eccentricity      float64
	Rotation          string //retrograde?
	ParentDesignation string //позиционный код родителя
	Designation       string
}

func StarOrbitPeriod(m1, m2 float64, distance au.AU) float64 {
	return math.Sqrt((math.Pow(float64(distance), 3)) / (m1 + m2))
}

func RollStellarOrbit(r stellar.Roller, star stellar.StarDesignation) Orbit {
	o := Orbit{}
	switch star {
	case stellar.Primary:
		return o
	case stellar.Close:
		o.Designation = string(stellar.Primary)
		//1D-1 (A result of 0 = Orbit# 0.5 or 0.2AU)
		o.Distance = float64(r.Roll("1d6-1")) + fluxDecimals(r)
		if o.Distance < 1 {
			o.Distance = 0.5
		}
	case stellar.Near:
		o.Designation = string(stellar.Primary)
		o.Distance = float64(r.Roll("1d6+5")) + fluxDecimals(r)
	case stellar.Far:
		o.Designation = string(stellar.Primary)
		o.Distance = float64(r.Roll("1d6+11")) + fluxDecimals(r)
	case stellar.PrimaryComp:
		o.Distance = rollCompanionOrbit(r)
		o.Designation = string(stellar.Primary)
	case stellar.CloseComp:
		o.Distance = rollCompanionOrbit(r)
		o.Designation = string(stellar.Close)
	case stellar.NearComp:
		o.Distance = rollCompanionOrbit(r)
		o.Designation = string(stellar.Near)
	case stellar.FarComp:
		o.Distance = rollCompanionOrbit(r)
		o.Designation = string(stellar.Far)
	}
	o.Eccentricity = -1
	return o
}

func rollCompanionOrbit(r stellar.Roller) float64 {
	return (float64(r.Roll("1d6")) / 10.0) + (float64(r.Roll("2d6-7") / 100.0))
}

func fluxDecimals(r stellar.Roller) float64 {
	r1 := r.Roll("2d6")
	r2 := r.Roll("2d6")
	return float64(r1-r2) * 0.1
}

func RollStarEccentricity(r stellar.Roller, dm int) float64 {
	r1 := r.Roll("2d6") + dm
	r2 := 0.0
	base := 0.0
	switch r1 {
	case 6, 7:
		r2 = float64(r.Roll("1d6")) * 0.005
		base = 0.0
	case 8, 9:
		r2 = float64(r.Roll("1d6")) * 0.01
		base = 0.03
	case 10:
		r2 = float64(r.Roll("1d6")) * 0.05
		base = 0.05
	case 11:
		r2 = float64(r.Roll("2d6")) * 0.05
		base = 0.05
	default:
		if r1 <= 5 {
			r2 = float64(r.Roll("1d6")) * 0.001
			base = -0.001
		}
		if r1 >= 12 {
			r2 = float64(r.Roll("1d6")) * 0.05
			base = 0.3
		}
	}

	return base + r2

}
