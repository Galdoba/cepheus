package starsystem

type StarSystem struct {
	SectorType string
}

type StarSystebObject interface {
	Type() string
	Orbit() float64
}
