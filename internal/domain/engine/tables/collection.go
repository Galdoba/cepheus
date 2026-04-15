package tables

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	tc.rollSequence = []string{}
	tc.results = []string{}
}

func (tc *Collection) Roll(roller TableRoller, name string, mods ...int) (string, error) {
	table := GameTable{}
	var err error
	if roller == nil {
		return "", fmt.Errorf("nil roller provided")
	}
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
		indexStr = roller.D66(mods...)
		result = table.Data[indexStr]
	case false:
		index, err = roller.Roll(table.Expression, mods...)
		if err != nil {
			return "", fmt.Errorf("roll on table %q (expression %q %v) failed: %w", table.Name, table.Expression, mods, err)
		}
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

func (tc *Collection) RollCascade(roller TableRoller, name string) (string, error) {
	maxDepth := 1000
	currentTable := name
	if name == "" {
		return "", fmt.Errorf("no name for starting table")
	}
	for depth := range maxDepth {
		result, err := tc.Roll(roller, currentTable)
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

func Save(t GameTable, path string) error {
	data, err := json.MarshalIndent(&t, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal table: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open/create file: %w", err)
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("failed to clean table file: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write table to file: %w", err)
	}
	return nil
}

func Load(path string) (GameTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GameTable{}, fmt.Errorf("failed to read table file: %w", err)
	}
	tab := GameTable{}
	if err := json.Unmarshal(data, &tab); err != nil {
		return GameTable{}, fmt.Errorf("failed to unmarshal table data: %w", err)
	}
	return tab, nil
}
