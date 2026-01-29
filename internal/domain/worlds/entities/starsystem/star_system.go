package starsystem

import (
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
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
}

type Star struct {
	Type        string
	SubType     string
	Class       string
	Dead        bool
	Protostar   bool
	Mass        float64
	Temperature float64
	Diameter    float64
	Luminocity  float64
	Age         float64
	Anomaly     string
}
