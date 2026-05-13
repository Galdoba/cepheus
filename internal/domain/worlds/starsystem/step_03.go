package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/tables"
)

const (
	DwarfType = "Dwarf Type"
)

func (b *Builder) Step_03(ss *StarSystem) error {
	tableCollection, err := dwarfTypeTable()
	if err != nil {
		return err
	}
	for i, s := range b.precursor.stars {
		if s.StellarClass == "BD" {
			dwarfType, err := tableCollection.Roll(b.dice, ObjectType)
			if err != nil {
				return fmt.Errorf("failed to roll table collection %q: %w", tableCollection.Name, err)
			}
			b.precursor.stars[i].StellarClass = dwarfType
			fmt.Println("set", b.precursor.stars[i].StellarClass)
		}
	}
	b.nextStep = step_05

	return nil
}

func dwarfTypeTable() (*tables.Collection, error) {
	tableCollection, err := tables.NewCollection(
		ObjectType, tables.New(
			ObjectType, "1d100",
			map[string]string{
				"01 - 50":  "L",
				"51 - 75":  "T",
				"76 - 100": "Y",
			}))
	if err != nil {
		return nil, fmt.Errorf("failed to create dwarf type tables: %w", err)
	}
	return tableCollection, nil
}
