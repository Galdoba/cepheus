package orbit

import (
	"math"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/au"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
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
	case stellar.Near:
		o.Designation = string(stellar.Primary)
	case stellar.Far:
		o.Designation = string(stellar.Primary)
	case stellar.PrimaryComp:
		o.Distance = (float64(r.Roll("1d6")) / 10.0) + (float64(r.Roll("2d6-7") / 100.0))
	case stellar.CloseComp:
		o.Distance = (float64(r.Roll("1d6")) / 10.0) + (float64(r.Roll("2d6-7") / 100.0))
	case stellar.NearComp:
		o.Distance = (float64(r.Roll("1d6")) / 10.0) + (float64(r.Roll("2d6-7") / 100.0))
	case stellar.FarComp:
		o.Distance = (float64(r.Roll("1d6")) / 10.0) + (float64(r.Roll("2d6-7") / 100.0))
	}

}
