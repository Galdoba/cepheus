package population

import (
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/uwp"
)

type PopulationData struct {
	Code                          string
	ValueDetailed                 float64
	PopulationConcentrationRating int
	MajorCities                   int
	UrbanPopulation               int
	MajorCityPopulation           int
}

func GenerateDetails(dp *dice.Dicepool, profile *uwp.UWP) *PopulationData {
	pd := PopulationData{}
	pd.Code = profile.CodeOf(uwp.Pops)
	return &pd
}
