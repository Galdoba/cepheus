# Tables Package

Generic RPG table lookup system with index parsing, validation, and cascade support.

## Overview

The `tables` package provides a system for defining and rolling on RPG lookup tables.
Tables map dice roll results (or D66 indexes) to string values. A `Collection` groups
multiple tables and supports cascading rolls — where a result on one table triggers
a re-roll on another.

The package is decoupled from the dice engine via the `TableRoller` interface,
allowing any dice implementation to be used.

## Import

```go
import "github.com/Galdoba/cepheus/internal/domain/engine/tables"
```

## Quick Start

```go
// Create a table: name, dice expression, index→value map
encounterTable := tables.New("encounters", "2d6", map[string]string{
    "2":  "Nothing",
    "3":  "Drift debris",
    "4":  "Merchant ship",
    "5":  "Patrol craft",
    "6":  "Pirates",
    "7":  "Nomads convoy",
    "8":  "Scout vessel",
    "9":  "Corsair raiders",
    "10": "Imperial fleet",
    "11": "Alien artifact",
    "12": "Space hulk",
})

// Validate before use
if err := encounterTable.Validate(); err != nil {
    log.Fatal(err)
}

// Create a collection for multi-table lookups
collection, err := tables.NewCollection("subsector", encounterTable)
if err != nil {
    log.Fatal(err)
}

// Roll on a table — requires a TableRoller (e.g., dice.Manager)
mgr, _ := dice.New("")
result, err := collection.Roll(mgr, "encounters")
// result → "Pirates" (if 2d6 rolled 6)
```

---

## GameTable

A `GameTable` maps index strings to result values:

```go
type GameTable struct {
    Name       string            // Table identifier
    Expression string            // Dice expression (e.g., "2d6", "d66")
    Data       map[string]string // Index → result mapping
    D66        bool              // Whether this uses D66 indexing
}
```

### Creating a Table

```go
func New(name, expression string, data map[string]string) GameTable
```

The `D66` flag is set automatically if `expression` is `"d66"` or `"D66"`.

### Index String Format

Keys in the `Data` map support flexible notation:

| Format | Example | Meaning |
|--------|---------|---------|
| Single value | `"7"` | Matches roll result 7 |
| Range | `"3-5"` | Matches 3, 4, or 5 |
| Open-ended high | `"11+"` | Matches 11 through +1000 |
| Open-ended low | `"2-"` | Matches -1000 through 2 |
| Comma-separated | `"2, 4, 6"` | Matches 2, 4, or 6 |
| Mixed | `"2-4, 6, 9+"` | Matches 2–4, 6, or 9+ |

Ranges are **inclusive** and automatically deduplicated.

### Validation

```go
func (t GameTable) Validate() error
```

Checks performed:

| Check | Description |
|-------|-------------|
| Name not empty | Table must have a name |
| Minimum 2 entries | At least 2 data entries required |
| Expression validity | Expression matches `XdY[+Z/-Z]` or `d66` |
| Index parseability | Every key parses via `stringToIndexes` |
| No duplicate indexes | Set-based detection across all keys (catches overlapping ranges) |
| No range holes (non-D66) | Expanded indexes must be contiguous from min to max |
| Bounds check [-1000, 1000] | All indexes must be within range |
| No sentinel leakage | `andAbove` (1001) / `andBelow` (-1001) must not appear as keys |
| No empty values | Every result string must be non-empty |

D66 tables are exempt from the range-holes check (sparse D66 tables are valid).

---

## Collection

A `Collection` groups multiple tables and tracks roll history:

```go
type Collection struct {
    Name         string
    Tables       map[string]GameTable
    // unexported: rollSequence, results
}
```

### Creating a Collection

```go
func NewCollection(name string, tables ...GameTable) (*Collection, error)
```

- Rejects duplicate table names at construction time
- Validates every table during construction
- Requires at least one table
- Requires a non-empty collection name

### Rolling

```go
func (tc *Collection) Roll(roller TableRoller, name string, mods ...int) (string, error)
```

- Looks up the table by `name`
- For D66 tables: calls `roller.D66(mods...)` to get a string index
- For standard tables: calls `roller.Roll(expression, mods...)` to get an int index, then matches against keys
- Records the table name and result in `rollSequence` and `results`
- Returns an error if the roll result doesn't match any key

Error messages include the table name, numeric index, and D66 index string for debugging.

### Cascading Roll

```go
func (tc *Collection) RollCascade(roller TableRoller, name string) (string, error)
```

Rolls on the starting table. If the result matches another table's name in the
collection, re-rolls on that table. Continues until a result is reached that is
**not** a table name.

Maximum depth: 1000 rolls. Exceeding this returns an error (detects infinite loops).

Example:
```
Table A: "1" → "B", "2" → "end"
Table B: "1" → "final", "2" → "end"

RollCascade(roller, "A") with roll=1 on A:
  A rolls 1 → "B" (is a table name)
  B rolls 1 → "final" (not a table name)
  Returns "final"
```

### Reset

```go
func (tc *Collection) Reset()
```

Clears `rollSequence` and `results` history. Useful for reusing a collection
across multiple generation runs.

### Validation

```go
func (tc *Collection) Validate() error
```

Validates every table in the collection. Called automatically during `NewCollection`.

---

## TableRoller Interface

The dice-rolling abstraction that `Collection` depends on:

```go
type TableRoller interface {
    D66(mods ...int) string
    Roll(expression string, mods ...int) (int, error)
}
```

The `dice.Manager` type from the `dice` package implements this interface.
Custom implementations can be provided for testing or alternative dice engines.

### Mock Roller (for testing)

```go
type mockRoller struct {
    d66Result   string
    rollResults map[string]int
    rollErr     error
}
```

The test suite includes a mock implementation — copy the pattern for your own tests.

---

## JSON Persistence

Tables can be saved to and loaded from JSON files:

```go
// Save a table to a JSON file
err := tables.Save(myTable, "/path/to/my_table.json")

// Load a table from a JSON file
loadedTable, err := tables.Load("/path/to/my_table.json")
```

The JSON format matches the `GameTable` struct tags:

```json
{
  "name": "encounters",
  "expression": "2d6",
  "data": {
    "2": "Nothing",
    "3": "Drift debris"
  },
  "d_66": false
}
```

---

## Index Parsing

### `stringToIndexes`

Converts an index string to a sorted, deduplicated slice of integers:

```go
indexes, err := stringToIndexes("2-4, 6, 9+")
// → []int{2, 3, 4, 6, 9, 10, 11, ..., 1000}
```

Supported patterns:

| Pattern | Regex | Behavior |
|---------|-------|----------|
| Range | `^(-?\d+)\s*-\s*(-?\d+)$` | Expands min–max inclusive |
| Open high | `^(-?\d+)\+$` | Expands to `DefaultUpperBound` (+1000) |
| Open low | `^(-?\d+)\-$` | Expands from `DefaultLowerBound` (-1000) |
| Single | `^(-?\d+)$` | Returns single value |

### `indexesToString`

Converts a slice of integers back to a compact string representation:

```go
s, err := indexesToString(2, 3, 4, 5, 8, 9, 15, 21)
// → "2 - 5, 8 - 9, 15, 21"

s, err := indexesToString(2, andAbove)
// → "2+"

s, err := indexesToString(-50, andBelow)
// → "-50-"
```

Sentinel values `andAbove` (1001) and `andBelow` (-1001) produce `+`/`-` suffix notation.
Exactly 2 values are required when using sentinels.

---

## Architecture

### Files

| File | Purpose |
|------|---------|
| `table.go` | `GameTable`, `Validate()`, index parsing (`stringToIndexes`, `indexesToString`), expression validation, `Save`/`Load` |
| `collection.go` | `Collection`, `NewCollection`, `Roll`, `RollCascade`, `Reset`, `Validate` |
| `roller.go` | `TableRoller` interface definition |
| `table_test.go` | 29 tests covering validation, parsing, collection ops, cascade |

### Design

```
┌─────────────────────────────────────┐
│           Collection                 │
│  ┌──────────┐  ┌──────────┐         │
│  │ Table A  │  │ Table B  │  ...    │
│  │ (2d6)    │  │ (d66)    │         │
│  └────┬─────┘  └────┬─────┘         │
│       │             │                │
│       ▼             ▼                │
│  ┌─────────────────────────┐         │
│  │     TableRoller         │         │
│  │  Roll() / D66()         │         │
│  └─────────────────────────┘         │
│         ▲                            │
│         │ (implemented by)           │
│  ┌──────┴──────┐                     │
│  │ dice.Manager │                     │
│  └─────────────┘                     │
└─────────────────────────────────────┘
```

---

## Bounds and Defaults

| Constant | Value | Purpose |
|----------|-------|---------|
| `DefaultUpperBound` | 1000 | Maximum index for `N+` expansion |
| `DefaultLowerBound` | -1000 | Minimum index for `N-` expansion |
| `andAbove` | 1001 | Sentinel for `+` notation in `indexesToString` |
| `andBelow` | -1001 | Sentinel for `-` notation in `indexesToString` |

Index bounds for validation: **[-1000, 1000]**. Values outside this range
are rejected by `Validate()`.

---

## Error Messages

All errors are contextual:

```
// Collection construction
"failed to create collection: duplicating table names provided: \"encounters\""

// Missing table in collection
"table \"nonexistent\" not found in collection \"subsector\""

// Nil roller
"nil roller provided"

// Roll result not found
"result is empty in table \"encounters\" (index=13 (or <not set>))"

// Cascade depth exceeded
"cascade exceeded max depth 1000"

// Validation errors
"table name cannot be empty"
"table \"test\" must have at least 2 entries"
"table \"test\" has holes in index range [1, 6]"
"table \"test\": index duplication: 3"
"table \"test\" has index 1001 out of bounds [-1000, 1000]"
"table \"test\" contains marker index 1001"
"table \"test\" has empty value"
```

---

## Contributing

When adding new validation checks to `GameTable.Validate()`:

1. Place the check after existing ones in logical order
2. Use `fmt.Errorf` with the table name for context
3. Add a test case in `table_test.go`
4. Document the check in the validation table above

When modifying `stringToIndexes` or `indexesToString`:

1. Update the supported patterns table
2. Add test cases for new patterns
3. Verify round-trip consistency: `indexesToString(stringToIndexes(s)) == s`
