package tttable

import (
	"fmt"
	"io"
)

// TableCollection manages multiple tables
type TableCollection struct {
	Tables map[string]*Table
	roller Roller
	writer io.Writer
}

// CollectionOption is a functional option for configuring a TableCollection
type CollectionOption func(*TableCollection) error

// WithTables adds multiple tables to the collection
func WithTables(tables ...*Table) CollectionOption {
	return func(tc *TableCollection) error {
		for _, table := range tables {
			if err := tc.AddTable(table); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithRoller adds a roller to the collection
func WithRoller(roller Roller) CollectionOption {
	return func(tc *TableCollection) error {
		if tc.roller != nil {
			return fmt.Errorf("option duplicated: roller")
		}
		tc.roller = roller
		return nil
	}
}

// WithWriter adds a writer to the collection
func WithWriter(w io.Writer) CollectionOption {
	return func(tc *TableCollection) error {
		if tc.writer != nil {
			return fmt.Errorf("option duplicated: writer")
		}
		tc.writer = w
		return nil
	}
}

// NewTableCollection creates a new table collection with options
func NewTableCollection(opts ...CollectionOption) (*TableCollection, error) {
	tc := &TableCollection{
		Tables: make(map[string]*Table),
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(tc); err != nil {
			return nil, err
		}
	}

	return tc, nil
}

// SetRoller sets the random number generator for the entire collection
func (tc *TableCollection) SetRoller(roller Roller) {
	tc.roller = roller
}

// SetRoller sets the random number generator for the entire collection
func (tc *TableCollection) SetWriter(w io.Writer) {
	tc.writer = w
}

// AddTable adds a table to the collection
func (tc *TableCollection) AddTable(table *Table) error {
	if _, exists := tc.Tables[table.Name]; exists {
		return fmt.Errorf("table with name %s already exists", table.Name)
	}
	tc.Tables[table.Name] = table
	return nil
}

// GetTable retrieves a table by name
func (tc *TableCollection) GetTable(name string) (*Table, bool) {
	table, exists := tc.Tables[name]
	return table, exists
}

// RemoveTable removes a table from the collection
func (tc *TableCollection) RemoveTable(name string) {
	delete(tc.Tables, name)
}

// Roll performs a cascade roll starting from the specified table
// If the result matches another table name, it continues rolling
// Any mods provided will substitute own table mods globally
// Returns the final result and any error encountered
func (tc *TableCollection) Roll(tableName string, mods ...string) (string, error) {
	_, results, err := tc.rollCascadeInternal(tableName, nil, 0, mods...)
	if err != nil {
		return "", err
	}

	// Return the last result
	if len(results) == 0 {
		return "", fmt.Errorf("no results generated")
	}
	return results[len(results)-1], nil
}

// RollCascade performs a cascade roll and returns all intermediate results
func (tc *TableCollection) RollCascade(tableName string, mods ...string) ([]int, []string, error) {
	return tc.rollCascadeInternal(tableName, nil, 0, mods...)
}

// FindByIndex return result from specified table by index
func (tc *TableCollection) FindByIndex(tableName string, index int) (string, error) {
	if tab, ok := tc.Tables[tableName]; ok {
		return tab.FindByRoll(index)
	}
	return "", fmt.Errorf("no table %v in collection", tableName)
}

// rollCascadeInternal is the internal implementation for cascade rolls
// visited is a map to detect cycles, depth is current recursion depth
// Any mods provided will substitute own table mods globally
func (tc *TableCollection) rollCascadeInternal(tableName string, visited map[string]bool, depth int, mods ...string) ([]int, []string, error) {
	// Check recursion depth to prevent infinite loops
	if depth > 100 {
		return nil, nil, ErrCascadeTooDeep
	}

	// Initialize visited map on first call
	if visited == nil {
		visited = make(map[string]bool)
	}

	// Check for cycles
	if visited[tableName] {
		return nil, nil, fmt.Errorf("cycle detected: table %s already visited", tableName)
	}
	visited[tableName] = true

	// Get the table
	table, ok := tc.GetTable(tableName)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrTableNotFound, tableName)
	}

	// Check if roller is set
	if tc.roller == nil {
		return nil, nil, ErrRollerNotSet
	}

	if tc.writer != nil {
		table.SetWriter(tc.writer)
	}

	// Roll on the table
	index, result, err := table.roll(tc.roller, mods...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to roll on table %s: %w", tableName, err)
	}

	// If the result is another table in the collection, continue cascading
	if _, ok := tc.GetTable(result); ok {
		// Recursively roll on the next table
		nextIndexes, nextResults, err := tc.rollCascadeInternal(result, visited, depth+1, mods...)
		if err != nil {
			return nil, nil, err
		}

		// Combine results: current result (which is a table name) + all next results
		allResults := make([]string, 0, len(nextResults)+1)
		allResults = append(allResults, result)
		allResults = append(allResults, nextResults...)
		allIndexes := make([]int, 0, len(nextIndexes)+1)
		allIndexes = append(allIndexes, index)
		allIndexes = append(allIndexes, nextIndexes...)

		return allIndexes, allResults, nil
	}

	// Result is not a table name, return it as the final result
	return []int{index}, []string{result}, nil
}

// GetTableNames returns all table names in the collection
func (tc *TableCollection) GetTableNames() []string {
	names := make([]string, 0, len(tc.Tables))
	for name := range tc.Tables {
		names = append(names, name)
	}
	return names
}

// Clear removes all tables from the collection
func (tc *TableCollection) Clear() {
	tc.Tables = make(map[string]*Table)
}
