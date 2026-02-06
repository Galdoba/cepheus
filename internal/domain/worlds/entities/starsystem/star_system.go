package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/pkg/dice"
)

type StarSystem struct {
	ID           string
	coordinates  coordinates.Cube
	importedData t5ss.WorldData
	Empty        bool
	Dead         bool
	Primordial   bool
	Clustered    bool
	NebulaType   int
	Anomaly      bool
	PrimaryStar  *Star
	Age          float64
	Stars        map[orbit.Orbit]*Star
}

type Star struct {
	Type        string
	SubType     string
	Class       string
	Designation stellar.StarDesignation
	Dead        bool
	Protostar   bool
	Mass        float64
	Temperature float64
	Diameter    float64
	Luminocity  float64
	Age         float64
	Anomaly     string
	Period      float64
}

func rollNebula(r *dice.Roller) int {
	switch r.Roll("1d6") {
	case 1:
		return 1
	default:
		return 2
	}
}

func (ss *StarSystem) Profile() string {
	p := ""
	si := newStarIterator(ss.Stars)
	starPos := 0
	for si.next() {
		starPos++
		o, s, err := si.getValues()
		if err != nil {
			panic(err)
		}
		switch starPos {
		case 1:
			p += fmt.Sprintf("%v-", len(ss.Stars))
			p += fmt.Sprintf("%v%v", s.Type, s.SubType)
			if s.Class != "" {
				p += fmt.Sprintf(" %v", s.Class)
			}
			p += fmt.Sprintf("-%v", float.RoundN(s.Mass, 3))
			p += fmt.Sprintf("-%v", float.RoundN(s.Diameter, 3))
			p += fmt.Sprintf("-%v", float.RoundN(s.Luminocity, 3))
			p += fmt.Sprintf("-%v", float.RoundN(ss.Age, 3))
		default:
			p += fmt.Sprintf(":%v", s.Designation)
			p += fmt.Sprintf("-%v", o.Distance)
			p += fmt.Sprintf("-%v", o.Eccentricity)
			p += fmt.Sprintf("-%v%v", s.Type, s.SubType)
			if s.Class != "" {
				p += fmt.Sprintf(" %v", s.Class)
			}
			p += fmt.Sprintf("-%v", float.RoundN(s.Mass, 3))
			p += fmt.Sprintf("-%v", float.RoundN(s.Diameter, 3))
			p += fmt.Sprintf("-%v", s.Luminocity)
		}

	}

	return p
}
