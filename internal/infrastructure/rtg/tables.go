package rtg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Galdoba/cepheus/pkg/tttable"
)

type RandomTableGenerator interface {
	Roll(string, ...string) (int, string, error)
	FindByIndex(string, int) (string, error)
}

type randomTableGenerator struct {
	tableCollection *tttable.TableCollection
}

func (rg *randomTableGenerator) Roll(name string, mods ...string) (int, string, error) {
	index, outcome, err := rg.tableCollection.Roll(name, mods...)
	if err != nil {
		return tttable.AndLess, "", fmt.Errorf("failed to roll RTG: %v", err)
	}
	return index, outcome, nil
}

func (rg *randomTableGenerator) FindByIndex(tableName string, index int) (string, error) {
	return rg.tableCollection.FindByIndex(tableName, index)
}

func tableFileName(tableName string) string {
	name := strings.ToLower(tableName)
	name = strings.ReplaceAll(name, " ", "_")
	return name + ".json"
}

func CreateRandomTable(dir, name, diceExpr string, rows map[string]string, mods map[string]int) error {
	rowList := []tttable.Row{}
	for key, val := range rows {
		rowList = append(rowList, tttable.NewRow(key, val))
	}
	tab, err := tttable.NewTable(name,
		tttable.WithDiceExpression(diceExpr),
		tttable.WithRows(rowList...),
		tttable.WithMods(mods),
	)
	if err != nil {
		return fmt.Errorf("failed to create table %v: %v", name, err)
	}
	if err := tab.SaveAs(filepath.Join(dir, tableFileName(name))); err != nil {
		return fmt.Errorf("failed to save table %v: %v", name, err)
	}
	return nil
}

func AssertTable(path string) error {
	tab, err := tttable.Load(path)
	if err != nil {
		return err
	}
	return tab.Validate()
}
