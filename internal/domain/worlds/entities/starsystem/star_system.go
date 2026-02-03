package starsystem

import (
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/pkg/dice"
)

const (
	Primary     = "Aa"
	PrimaryComp = "Ab"
	Close       = "Ba"
	CloseComp   = "Bb"
	Near        = "Ca"
	NearComp    = "Cb"
	Far         = "Da"
	FarComp     = "Db"
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
	Designation string
	Dead        bool
	Protostar   bool
	Mass        float64
	Temperature float64
	Diameter    float64
	Luminocity  float64
	Age         float64
	Anomaly     string
}

func rollNebula(r *dice.Roller) int {
	switch r.Roll("1d6") {
	case 1:
		return 1
	default:
		return 2
	}
}
