package rtg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Galdoba/cepheus/pkg/tttable"
)

type RandomTableGenerator interface {
	Roll(string, ...string) (string, error)
	Find(string, string) (string, error)
}

type randomTableGenerator struct {
	tableCollection *tttable.TableCollection
}

func (rg *randomTableGenerator) Roll(name string, mods ...string) (string, error) {
	outcome, err := rg.tableCollection.Roll(name, mods...)
	if err != nil {
		return "", fmt.Errorf("failed to roll RTG: %v", err)
	}
	return outcome, nil
}

func (rg *randomTableGenerator) Find(tableName string, index string) (string, error) {
	return rg.tableCollection.Find(tableName, index)
}

func tableFileName(tableName string) string {
	name := strings.ToLower(tableName)
	name = strings.ReplaceAll(name, " ", "_")
	return name + ".json"
}

func CreateRandomIndexTable(dir, name, diceExpr string, rows map[string]string, mods map[tttable.ModType]map[string]int) error {
	rowList := []tttable.TableEntry{}
	for key, val := range rows {
		rowList = append(rowList, tttable.NewTableEntry(key, val))
	}
	tab, err := tttable.NewTable(name,
		tttable.WithDiceExpression(diceExpr),
		tttable.WithIndexEntries(rowList...),
		tttable.WithIndexMods(tttable.Flat, mods[tttable.Flat]),
		tttable.WithIndexMods(tttable.Cumulative, mods[tttable.Cumulative]),
		tttable.WithIndexMods(tttable.Max, mods[tttable.Max]),
		tttable.WithIndexMods(tttable.Min, mods[tttable.Min]),
	)
	if err != nil {
		return fmt.Errorf("failed to create table %v: %v", name, err)
	}
	if err := tttable.SaveAs(tab, filepath.Join(dir, tableFileName(name))); err != nil {
		return fmt.Errorf("failed to save table %v: %v", name, err)
	}
	return nil
}

type mods map[tttable.ModType]map[string]int

func NewMods() mods {
	modMap := make(map[tttable.ModType]map[string]int)
	return modMap
}

func (mm mods) AddMod(mType tttable.ModType, descr string, value int) mods {
	if mm[mType] == nil {
		mm[mType] = make(map[string]int)
	}
	mm[mType][descr] = value
	return mm
}

func AssertTable(path string) error {
	tab, err := tttable.Load(path)
	if err != nil {
		return err
	}
	return tab.Validate()
}
