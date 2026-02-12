package rtg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/cepheus/internal/infrastructure/filepaths"
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/tttable"
)

const (
	systemPlanetsSubfolder = "system_planets"

	MOD_SingleClassV         = "System consists of a single Class V star"
	MOD_PrimaryIsBD          = "Primary star is Brown Dwarf"
	MOD_PrimaryIsPostStellar = "Primary star is post stellar object"
	MOD_PerEveryPostStellar  = "Per every post stellar object"
	MOD_FourOrMoreStars      = "System consists of four or more stars"

	TableGGQuantity = "Gas Giant Quantity"
)

func NewSystemPlanetsDeterminationGenerator(roller *dice.Roller, opts ...string) (RandomTableGenerator, error) {
	options := make(map[string]bool)
	for _, o := range opts {
		options[o] = true
	}
	tables := []tttable.RollableTable{}
	for name, path := range StarSystemPlanetsTableMap() {
		table, err := tttable.Load(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load table %v: %v", name, err)
		}
		tables = append(tables, table)
	}
	tc, err := tttable.NewTableCollection(
		tttable.WithTables(tables...),
	)
	tc.SetRoller(roller)
	if err != nil {
		return nil, fmt.Errorf("failed to create random tables colection: %v", err)
	}

	generator := randomTableGenerator{}
	generator.tableCollection = tc
	return &generator, nil
}

func StarSystemPlanetsTableMap() map[string]string {
	tableMap := make(map[string]string)
	for _, name := range starSystemPlanetsTableNames() {
		tableMap[name] = filepath.Join(filepaths.RandomTablesDirectory(), systemPlanetsSubfolder, tableFileName(name))
	}
	return tableMap
}

func starSystemPlanetsTableNames() []string {
	return []string{
		TableGGQuantity,
	}
}

func InitStarSystemPlanetsDeterminationTables() error {
	for name, path := range StarSystemPlanetsTableMap() {
		dir := filepath.Dir(path)
		fmt.Fprintf(os.Stderr, "init table: %v                \r", name)
		err := AssertTable(path)
		if err == nil {
			continue
		}
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0777); err != nil {
				return fmt.Errorf("failed to create random tables directory: %v", err)
			}
			path = dir
			switch name {
			case TableGGQuantity:
				err = CreateRandomIndexTable(path, name, "2d6",
					map[string]string{
						"4-":   "1",
						"5-6":  "2",
						"7-8":  "3",
						"9-11": "4",
						"12":   "5",
						"13+":  "6",
					},
					NewMods().
						AddMod(tttable.Flat, MOD_SingleClassV, 1).
						AddMod(tttable.Flat, MOD_PrimaryIsBD, -2).
						AddMod(tttable.Flat, MOD_PrimaryIsPostStellar, -2).
						AddMod(tttable.Cumulative, MOD_PerEveryPostStellar, -1).
						AddMod(tttable.Flat, MOD_FourOrMoreStars, -1),
				)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "\n")
				return fmt.Errorf("failed to initiate table: %v", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "\n")
			return fmt.Errorf("failed to initiate table %v: %v", name, err)
		}
	}
	return nil
}
