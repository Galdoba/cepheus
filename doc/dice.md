# Dice Package

Cepheus Engine dice rolling implementation for RPG game mechanics.

## Overview

The `dice` package provides a flexible dice rolling system with expression parsing, modifier chaining, expression caching, and special RPG mechanics (D66, Flux, Variance). It is the random-roll foundation for all other packages in the project.

## Import

```go
import "github.com/Galdoba/cepheus/internal/domain/engine/dice"
```

## Quick Start

```go
// Simple roll — returns (int, error)
result, err := dice.Roll("2d6")

// With a constant modifier added after the sum
result, err := dice.Roll("2d6", +3)

// MustRoll — panics on error (for hardcoded expressions)
result := dice.MustRoll("3d6:dl1")

// D66 — returns a two-digit string ("00"–"99")
idx := dice.D66()

// Flux — first d6 minus second d6, plus optional mods
f := dice.Flux()

// Variance — random float64 in [0.0, 1.0]
v := dice.Variance()
```

## Seeded Manager

The package-level functions use a default manager seeded from the current timestamp.
For deterministic results (tests, reproducible generation), create your own:

```go
mgr, err := dice.New("my-seed-string")
if err != nil {
    // handle error
}
result, err := mgr.Roll("2d6")
idx := mgr.D66()
```

The seed string is converted to an `int64` via a custom hash function
(`stringToInt64`), so any string produces a deterministic RNG state.

---

## Expression Syntax

### Grammar

```
expression   → dicePart [ ":" complexMods ] [ simpleAdditive ]
dicePart     → [count] ("d" | "D" | "DD") sides
simpleAdditive → ("+" | "-") number
complexMods  → complexMod [ ":" complexMod ]*
```

### Dice Part

| Expression | Meaning | Status |
|------------|---------|--------|
| `2d6` | 2 six-sided dice | ✅ |
| `d20` | 1 twenty-sided die (count defaults to 1) | ✅ |
| `3d10+5` | 3d10, add 5 to the sum (simple additive) | ✅ |
| `D6` | Concat dice — not implemented | ❌ returns error |
| `2DD10` | Destructive dice — not implemented | ❌ returns error |

### Complex Modifiers (after `:`)

Multiple modifiers can be chained, separated by `:`. They are applied in
**priority order**, not left-to-right.

| Modifier | Syntax | Description | Example |
|----------|--------|-------------|---------|
| Add to each | `+Ne` | Add N to every individual die before sum | `2d6:+1e` |
| Add individual | `+N>>M` | Add N to die at position M (1-based) | `2d6:+1>>1` |
| Add individual (neg) | `-N>>M` | Subtract N from die at position M | `2d6:-1>>2` |
| Drop lowest | `dlN` | Remove N lowest dice from the pool | `4d6:dl1` |
| Drop highest | `dhN` | Remove N highest dice from the pool | `3d6:dh1` |
| Divide | `/N` | Integer divide each value by N | `2d6:/2` |
| Multiply | `xN` or `*N` | Multiply each value by N | `2d6:x100` |

### Simple Additive (after expression, no colon)

A trailing `+N` or `-N` (without a colon prefix) adds a constant to the
**summed** result. This has priority 110, applied after summing but before
divide/multiply.

```
2d6+5   → roll 2d6, sum, add 5
3d10-2  → roll 3d10, sum, subtract 2
```

### Full Examples

```
2d6                 → 2d6, sum
2d6+5               → 2d6, sum, +5
d20-3               → 1d20, sum, -3
3d6:dl1             → 3d6, drop lowest 1, sum
3d6:dh1             → 3d6, drop highest 1, sum
2d6:+1e             → 2d6, +1 to each die, sum
2d6:+1e:/2          → 2d6, +1 each, sum, /2
4d6:dl1:x100        → 4d6, drop lowest 1, sum, ×100
2d6:+1>>1           → 2d6, +1 to first die only, sum
2d6:-2>>2           → 2d6, -2 from second die only, sum
3d6:dl1:+1e:*2      → 3d6, drop 1 lowest, +1 each, sum, ×2
```

---

## Modifier Priority Order

Modifiers are sorted by priority and applied lowest-first:

| Priority | Modifier | When applied |
|----------|----------|--------------|
| 20 | AddIndividual | To each die at a specific position |
| 30 | AddToEach | To every die in the pool |
| 70 | DropLowest | Remove lowest N dice |
| 71 | DropHighest | Remove highest N dice |
| 100 | Sum | Collapse pool to single summed value |
| 110 | AddConst | Add/subtract a constant to the sum |
| 120 | Divide | Integer divide by N |
| 130 | Multiply | Multiply by N |

Example: `3d6:dl1:+1e` → `addToEach` (30) then `dropLowest` (70) then `sum` (100).
The parser reorders modifiers regardless of input order.

---

## Public API

### Package-Level Functions

These use the default manager (random seed from current timestamp).

```go
// Roll parses, rolls, applies modifiers, and returns the sum.
func Roll(expr string, mods ...int) (int, error)

// MustRoll panics on parse or roll error. Use for hardcoded expressions.
func MustRoll(expr string, dm ...int) int

// D66 rolls two d6 and returns concatenated digits as a string.
// Each die is clamped to 0–9 after applying optional mods.
func D66(mods ...int) string

// Flux returns (first d6 - second d6) + mods.
func Flux(dm ...int) int

// FluxGood returns (max(d1,d2) - min(d1,d2)) + mods. Always non-negative.
func FluxGood(dm ...int) int

// FluxBad returns (min(d1,d2) - max(d1,d2)) + mods. Always non-positive.
func FluxBad(dm ...int) int

// Variance returns a random float64 in [0.0, 1.0].
func Variance() float64

// ValidateExpression checks whether an expression is syntactically valid.
func ValidateExpression(expr string) error
```

### Manager Methods

For seeded/deterministic use:

```go
// New creates a Manager seeded from the given string.
// An empty string produces a random seed.
func New(seed string) (*Manager, error)

// Roll evaluates the expression and returns the sum.
func (m *Manager) Roll(expr string, mods ...int) (int, error)

// MustRoll panics on error.
func (m *Manager) MustRoll(expr string, dm ...int) int

// D66 rolls two d6 and returns concatenated digits.
func (m *Manager) D66(mods ...int) string

// Flux, FluxGood, FluxBad — same as package-level versions.
func (m *Manager) Flux(dm ...int) int
func (m *Manager) FluxGood(dm ...int) int
func (m *Manager) FluxBad(dm ...int) int

// Variance returns a random float64 in [0.0, 1.0].
func (m *Manager) Variance() float64

// Result returns the last roll's raw dice. Thread-safe.
func (m *Manager) Result() Result
```

### Result Type

```go
type Result struct {
    // unexported fields — accessed via methods
}

func (r Result) Dice() []die   // copy of the dice objects
func (r Result) Raw() []int    // copy of the raw roll values
```

The `die` type is unexported; consumers interact with values through `Raw()`.

---

## Architecture

### Files

| File | Purpose |
|------|---------|
| `spec.go` | Core types: `Manager`, `Result`, `Roller` interface, default init |
| `api.go` | Public API: `Roll`, `MustRoll`, `D66`, `Flux`, `FluxGood`, `FluxBad`, `Variance` |
| `parse.go` | Expression parser: `parseExpression`, modifier parsing, `ValidateExpression` |
| `mods.go` | Modifier types: `addToEach`, `addIndividual`, `dropLowest`, `dropHighest`, `divide`, `multiply`, `addConst`, `summ` |
| `roll.go` | `basicRoll` — rolls every die in a dicepool |
| `manager.go` | `Manager` struct, `roll` method, `stdInterpreter` (applies mods to raw roll) |
| `roller.go` | `randRoller` — `math/rand`-based roller, `stringToInt64` seed hashing |
| `cache.go` | `expressionCache` — thread-safe `sync.RWMutex` cache of parsed expressions |
| `die.go` | `die` and `dicepool` types with builder methods |

### Roll Flow

```
Roll("2d6+5")
  │
  ├─► exprCache.get() ──hit──► use cached expression
  │                │
  │              miss
  │                ▼
  │         parseExpression()
  │           ├─► parseDicePart()  → []die
  │           ├─► parseComplexModifiers() → []mod
  │           └─► parseSimpleAdditive() → addConst
  │                │
  │                ▼
  │         exprCache.set()  (for next time)
  │                │
  ├────────────────┘
  ▼
basicRoll(roller, dicepool)
  │  └─► roller.roll(die) for each die → raw []int
  ▼
stdInterpreter.interpret()
  │  └─► apply each mod in priority order → final []int
  │  └─► sum final values → int
  ▼
return sum (+ optional mods)
```

### Thread Safety

`Manager` uses a `sync.Mutex` around the entire roll pipeline (parse → roll →
interpret), ensuring concurrent calls are serialized. The expression cache uses
its own `sync.RWMutex` for concurrent reads with exclusive writes.

---

## Error Handling

All errors are descriptive and include context:

```go
// Empty expression
dice.Roll("")
// → error: "empty expression"

// Invalid dice type
dice.Roll("3x6")
// → error: "invalid dice type"

// Missing sides
dice.Roll("3d")
// → error: "missing sides after d"

// Drop more dice than available
dice.Roll("2d6:dl2")
// → error: "cannot drop 2 dice from pool of 2"

// Division by zero
dice.Roll("2d6:/0")
// → error: "invalid divisor"

// Position out of range (AddIndividual)
dice.Roll("2d6:+1>>5")
// → error: "position 5 out of range (1..2)"

// Unknown modifier
dice.Roll("2d6:xyz")
// → error: "unknown complex modifier: xyz"
```

---

## D66 Implementation

`D66` rolls two d6 independently, applies optional per-die modifiers, clamps
each to 0–9, and concatenates:

```go
dice.D66()        // e.g., "37"
dice.D66(1, -1)   // +1 to first die, -1 to second
```

This differs from a standard `2d6` sum — the result is a two-digit string
suitable for table index lookup.

---

## Future Plans

### 1. Special Dice Types

**Concat Dice (`D`)**: Treat dice positions as digits to concatenate rather than sum.
Parser recognizes `D` but rolling returns an error.

**Destructive Dice (`DD`)**: Roll N dice, keep best M (Mongoose Traveller 2E advantage).
Parser recognizes `DD` but rolling returns an error.

### 2. Logger Integration

Add a `Logger` interface for roll auditing:

```go
type Logger interface {
    LogRoll(expr string, result Result)
    LogModifier(mod mod, input, output []int)
    LogError(expr string, err error)
}
```

### 3. Exploding Dice

`d6!` — reroll and add when the die shows its maximum value.

### 4. Reroll Mechanics

Reroll dice that show specific values (e.g., reroll 1s).

### 5. Probability Calculator

Compute probability distributions for expressions:

```go
type ProbabilityDistribution struct {
    Min    int
    Max    int
    Mean   float64
    StdDev float64
    Table  map[int]float64
}
```

---

## Contributing

When adding new modifiers:

1. Implement the `mod` interface (`apply([]int) ([]int, error)`, `priority() int`)
2. Add name and priority constants in `mods.go`
3. Add parser support in `parse.go` (`parseOneComplexModifier`)
4. Add tests
5. Ensure error messages include context (values, bounds, pool size)
