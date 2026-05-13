package rules

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
	RetrogradeMarker                     string
	Separator                            string
}

var StarSequance = []string{"Alpha", "Beta", "Gamma"}
var PlanetSequance = []string{"0", "1", "2"}
var SignificatnMoonSequance = []string{"a", "b", "c"}
var InsignivicantMoonSequance = insignificantMoons()

func insignificantMoons() []string {
	//fill from 001 to 999
	return []string{"001", "002", "003"}
}

// draft
func (odr *OrbitDesignationRules) DesignateOrbit(indexes ...int) string {
	code := ""
	for i, index := range indexes {
		switch i {
		case IndexStar:
			code = odr.StarDesignationSequance[index]
		case IndexPlanet:
			code += odr.Separator + odr.PlanetDesignationSequance[index]
		default:
			//other types
		}
	}
	return code
}
