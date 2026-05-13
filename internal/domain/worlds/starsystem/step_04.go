package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/tables"
)

const (
	StellarType = "Stellar Type"
	ClassO      = "O"
	ClassB      = "B"
	ClassA      = "A"
	ClassF      = "F"
	ClassG      = "G"
	ClassK      = "K"
	ClassM      = "M"
	ClassBD     = "BD"
	ClassL      = "L"
	ClassT      = "T"
	ClassY      = "Y"
)

type stellarClass string

func (b *Builder) Step_04(ss *StarSystem) error {
	tableCollection, err := stellarTypeTable()
	if err != nil {
		return err
	}
	for i, s := range b.precursor.stars {
		if s.StellarClass == star {
			starType, err := tableCollection.Roll(b.dice, string(b.cfg.SystemType))
			if err != nil {
				return fmt.Errorf("failed to roll table collection %q: %w", tableCollection.Name, err)
			}
			r := b.dice.MustRoll("1d6")
			switch r {
			case 1:
				starType = "BD"
				b.nextStep = step_03
			default:
				b.nextStep = "step_05"
			}
			b.precursor.stars[i].StellarClass = starType
			fmt.Println("set", starType)
		}
	}
	for _, s := range b.precursor.stars {
		fmt.Println(s)
	}
	return nil
}

func stellarTypeTable() (*tables.Collection, error) {
	tableCollection, err := tables.NewCollection(
		StellarType,
		tables.New(
			string(Realistic), "1d100",
			map[string]string{
				"01 - 80": ClassM,
				"81 - 88": ClassK,
				"89 - 94": ClassG,
				"95 - 97": ClassF,
				"98":      ClassA,
				"99":      ClassB,
				"100":     ClassO,
			},
		),
		tables.New(
			string(SemiRelistic), "1d100",
			map[string]string{
				"01 - 50": ClassM,
				"51 - 77": ClassK,
				"78 - 90": ClassG,
				"91 - 97": ClassF,
				"98":      ClassA,
				"99":      ClassB,
				"100":     ClassO,
			},
		),
		tables.New(
			string(Fantastic), "1d100",
			map[string]string{
				"01 - 25": ClassM,
				"26 - 50": ClassK,
				"51 - 75": ClassG,
				"76 - 97": ClassF,
				"98":      ClassA,
				"99":      ClassB,
				"100":     ClassO,
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dwarf type tables: %w", err)
	}
	return tableCollection, nil
}
