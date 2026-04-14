package tables

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

type Collection struct {
	Name         string
	Tables       map[string]GameTable
	rollSequence []string
	results      []string
}

func NewCollection(name string, tables ...GameTable) (*Collection, error) {
	tc := Collection{
		Name:         name,
		Tables:       make(map[string]GameTable, len(tables)),
		rollSequence: []string{},
		results:      []string{},
	}
	nameDetected := make(map[string]int)
	for _, t := range tables {
		tc.Tables[t.Name] = t
		nameDetected[t.Name]++
		if nameDetected[t.Name] > 1 {
			return nil, fmt.Errorf("failed to create collection: duplicating table names provided: %q", t.Name)
		}
	}
	if err := tc.Validate(); err != nil {
		return nil, err
	}
	return &tc, nil
}

func (tc *Collection) Reset() {
	tc.results = []string{}
}

func (tc *Collection) Roll(dm *dice.Manager, name string, mods ...int) (string, error) {
	table := GameTable{}
	if found, ok := tc.Tables[name]; !ok {
		return "", fmt.Errorf("table %q not found in collection %q", name, tc.Name)
	} else {
		table = found
	}
	index := -1002 //imposible index
	indexStr := "<not set>"
	result := ""
	switch table.D66 {
	case true:
		indexStr = dm.D66(mods...)
		result = table.Data[indexStr]
	case false:
		index = dm.Roll(table.Expression, mods...)
	search_loop:
		for indexStr, value := range table.Data {
			candidates, _ := stringToIndexes(indexStr) //skip error because it supposed to be validated by now
			for _, match := range candidates {
				if index != match {
					continue
				}
				result = value
				break search_loop
			}
		}
	}
	if result == "" {
		return result, fmt.Errorf("result is empty in table %q (index=%d (or %q))", table.Name, index, indexStr)
	}
	tc.results = append(tc.results, result)
	tc.rollSequence = append(tc.rollSequence, table.Name)
	return result, nil
}

func (tc *Collection) RollCascade(dm *dice.Manager, name string) (string, error) {
	maxDepth := 1000
	currentTable := name

	for depth := 0; depth < maxDepth; depth++ {
		result, err := tc.Roll(dm, currentTable)
		if err != nil {
			return "", fmt.Errorf("cascade failed at depth %d: %w", depth, err)
		}

		nextTable, ok := tc.Tables[result]
		if !ok {
			return result, nil
		}

		currentTable = nextTable.Name
	}

	return "", fmt.Errorf("cascade exceeded max depth %d", maxDepth)
}

func (tc *Collection) Validate() error {
	if len(tc.Name) == 0 {
		return fmt.Errorf("collection name cannot be empty")
	}
	if len(tc.Tables) < 1 {
		return fmt.Errorf("collection %q must have at least one table", tc.Name)
	}
	for _, t := range tc.Tables {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("failed to validate collection %q: %w", tc.Name, err)
		}
	}
	return nil
}
