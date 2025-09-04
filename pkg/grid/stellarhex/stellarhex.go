package stellarhex

import (
	"github.com/Galdoba/cepheus/pkg/grid/coordinates"
	"github.com/Galdoba/cepheus/pkg/dice"
)

const (
	Density_Undefined = 0
	Density_Void      = 1
	Density_Rift      = 2
	Density_Sparse    = 3
	Density_Scattered = 4
	Density_Standard  = 5
	Density_Dense     = 6
	Density_Cluster   = 7
	Density_Core      = 8
)

type StellarHex struct {
	Coordiantes     coordinates.SpaceCoordinates
	Density         int
	StarSystem      string
	D               string
	BD              string
	LGG             string
	MGG             []string //2
	SGG             []string //3
	RogueWorldBig   []string //4
	RogueWorldSmall []string //10
	Planetoids      []string //20
	PSR             string
	NS              string
	BH              string
	NB              string
}

func New(crd coordinates.SpaceCoordinates, options ...StellarHexOption) *StellarHex {
	sh := &StellarHex{
		Coordiantes: crd,
		Density:     Density_Undefined,
		//stars
		StarSystem: "?",
		NS:         "?",
		BH:         "?",
		D:          "?",
		BD:         "?",
		//planets
		LGG:             "?",
		MGG:             []string{"?", "?"},
		SGG:             []string{"?", "?", "?"},
		RogueWorldBig:   []string{"?", "?", "?", "?"},
		RogueWorldSmall: []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
		PSR:             "?",
		NB:              "?",
	}

	for _, enrich := range options {
		enrich(sh)
	}
	return sh
}

func (sh *StellarHex) GenerateMissingDetails(dp *dice.Dicepool) error {
	if sh.Density == Density_Undefined {
		sh.Density = generateDensity(dp)
	}
	if err := sh.defineStarLikeObjectsPresence(dp); err != nil {
		return err
	}
	if err := sh.defineRoguePlanetsPresence(dp); err != nil {
		return err
	}
	sh.confirmAbsenseOfTheRest()
	return nil
}

type StellarHexOption func(*StellarHex)

func Density(density int) StellarHexOption {
	return func(sh *StellarHex) {
		sh.Density = density
	}
}

////////////////
