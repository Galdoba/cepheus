package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/tables"
)

const (
	ObjectType    = "Object Type"
	star          = "Star"
	brownDwarf    = "Brown Dwarf"
	roguePlanet   = "Rogue Planet"
	rogueGasGiant = "Rogue Gas Giant"
	neutronStar   = "Neutron Star"
	nebula        = "Nebula"
	blachHole     = "Black Hole"
)

func (b *Builder) Step_02(ss *StarSystem) error {
	tableCollection, err := tables.NewCollection(
		ObjectType, tables.New(
			ObjectType, "1d100",
			map[string]string{
				"01 - 80": star,
				"81 - 88": brownDwarf,
				"89 - 94": roguePlanet,
				"95 - 97": rogueGasGiant,
				"98":      neutronStar,
				"99":      nebula,
				"100":     blachHole,
			}))
	if err != nil {
		return fmt.Errorf("failed to create step 02 tables: %w", err)
	}
	objType, err := tableCollection.Roll(b.dice, ObjectType)
	if err != nil {
		return fmt.Errorf("failed to roll table collection %q: %w", tableCollection.Name, err)
	}
	b.precursor.celestialObject = objType
	switch b.precursor.celestialObject {
	case star:
		b.nextStep = step_04
		b.precursor.AddStar(star)
	case brownDwarf:
		b.nextStep = step_03
		b.precursor.AddStar(ClassBD)
	case roguePlanet:
		b.nextStep = "step_15"
	case rogueGasGiant:
		b.nextStep = "step_13"
	case neutronStar, nebula:
		b.nextStep = "step_18"
	case blachHole:
		b.nextStep = done
	}

	return nil
}
