package star

import "github.com/Galdoba/cepheus/internal/domain/worlds/entities/abstractions"

type Star struct {
	Designation   string
	SpectralClass string
	SubClass      int
	Size          string
	Diameter      float64
	Luminocity    float64
	OrbitN        float64
	Planets       abstractions.WorldObject
}
