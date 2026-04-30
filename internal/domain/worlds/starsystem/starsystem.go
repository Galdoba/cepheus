package starsystem

const (
// Star          SystemObject = "Star"
// BrownDwarf    SystemObject = "Brown Dwarf"
// RoguePlanet   SystemObject = "Rogue Planet"
// RogueGasGiant SystemObject = "Rogue Gas Giant"
// NeutronStar   SystemObject = "Neutron Star"
// Nebula        SystemObject = "Nebula"
// BlackHole     SystemObject = "Black Hole"
)

type SystemObject string
type Orbit float64

type StarSystem struct {
	CentralSystemObject SystemObject
	Objects             []Star
	Objects2            map[Orbit]Star
}

type Star struct {
	StarFields string
	Planets    []Planet
	Planets2   map[Orbit]Planet
}

type Planet struct {
	PlanetFields string
	Moons        []Moon
	Moons2       map[Orbit]Moon
}

type Moon struct {
	MoonFields string
}
