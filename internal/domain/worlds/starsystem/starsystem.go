package starsystem

import (
	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
	"github.com/Galdoba/cepheus/internal/domain/engine/tables"
)

const (
	Star          SystemObject = "Star"
	BrownDwarf    SystemObject = "Brown Dwarf"
	RoguePlanet   SystemObject = "Rogue Planet"
	RogueGasGiant SystemObject = "Rogue Gas Giant"
	NeutronStar   SystemObject = "Neutron Star"
	Nebula        SystemObject = "Nebula"
	BlackHole     SystemObject = "Black Hole"
)

type SystemObject string

type StarSystem struct {
	CentralSystemObject SystemObject
}

func DetermineSystemObject(dp *dice.Manager) SystemObject {
	col, err := getCollection()
	if err != nil {
		panic(err)
	}
	ot, err := col.Roll(dp, tableObjectType)
	if err != nil {
		panic(0)
	}
	return SystemObject(ot)
}

const (
	tableObjectType = "Star System Object Type"
)

func getCollection() (*tables.Collection, error) {
	return tables.NewCollection("Star System Tables",
		tables.New(tableObjectType, "1d100", map[string]string{
			"01 - 80": string(Star),
			"81 - 88": string(BrownDwarf),
			"89 - 94": string(RoguePlanet),
			"95 - 97": string(RogueGasGiant),
			"98":      string(NeutronStar),
			"99":      string(Nebula),
			"00":      string(BlackHole),
		}),
	)
}
