package port

import (
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/uwp"
)

type Port struct {
	Code         string
	Fuel         string
	Repair       string
	Shipyard     string
	NavalBase    bool
	ScoutBase    bool
	MilitaryBase bool
	CorsairBase  bool
	Sensors      int
	DetailedData map[string]int
	// Highport     bool TODO: поместить в виде значения Триария в DetailedData
}

type portDetailGenerator struct {
	rng *dice.Roller
	uwp uwp.UWP
}
