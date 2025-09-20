package starsystem

import (
	"github.com/Galdoba/cepheus/iiss/types/orbit"
	"github.com/Galdoba/cepheus/iiss/types/star"
)

type StarSystem struct {
	// Primary     *star.Star
	Stars       map[string]*star.Star
	Orbits      map[float64]*orbit.Orbit
	presenceGG  int
	presenceBT  int
	presenceTP  int
	totalWorlds int
	MAO         float64
}

func NewStarSystem() StarSystem {
	ss := StarSystem{
		// Primary:    &star.Star{},
		Stars:      make(map[string]*star.Star),
		Orbits:     make(map[float64]*orbit.Orbit),
		presenceGG: -1,
		presenceBT: -1,
		presenceTP: -1,
	}
	return ss
}

func (ss *StarSystem) primary() *star.Star {
	return ss.Stars["Aa"]
}
