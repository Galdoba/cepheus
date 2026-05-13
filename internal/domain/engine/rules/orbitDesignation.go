package rules

import "fmt"

const (
	IndexStar = iota
	IndexPlanet
	IndexSMoon
	IndexIMoon
)

type OrbitDesignationRules struct {
	StarDesignationSequance              []string
	PlanetDesignationSequance            []string
	SignificantMoonDesignationSequance   []string
	InsignificantMoonDesignationSequance []string
}

var odr = &OrbitDesignationRules{}

func init() {
	odr.StarDesignationSequance = StarSequance
	odr.PlanetDesignationSequance = PlanetSequance
	odr.StarDesignationSequance = SignificantMoonSequance
	odr.StarDesignationSequance = InsignivicantMoonSequance
}

var StarSequance = []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"}
var PlanetSequance = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
var SignificantMoonSequance = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var InsignivicantMoonSequance = insignificantMoons()

func insignificantMoons() []string {
	s := []string{}
	for i := range 300 {
		s = append(s, fmt.Sprintf("%00d", i+1))
	}
	return s
}

func DesignateStar(i int) string {
	return odr.StarDesignationSequance[i]
}

func DesignatePlanet(s, p int) string {
	return odr.StarDesignationSequance[s] + " " + odr.PlanetDesignationSequance[p]
}

func DesignateSattelite(s, p, m int) string {
	return odr.StarDesignationSequance[s] + " " + odr.PlanetDesignationSequance[p] + " " + odr.SignificantMoonDesignationSequance[m]
}

func DesignateInsignificant(s, p, i int) string {
	return odr.StarDesignationSequance[s] + " " + odr.PlanetDesignationSequance[p] + " " + odr.InsignificantMoonDesignationSequance[i]
}
