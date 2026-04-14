# Dice Package

Cepheus Engine dice rolling implementation for RPG game mechanics.

## Overview

The `dice` package provides dice expression parsing and rolling functionality for the Cepheus Engine rules system. It supports standard dice notation, modifiers, caching, and proper error handling.

## Installation

```go
import "github.com/Galdoba/cepheus/internal/domain/engine/dice"
```

## Quick Start

```go
// Simple roll
result := dice.Roll("2d6")

// With modifiers
result := dice.Roll("3d6:dl1") // roll 3d6, drop lowest 1

// Safe version with error handling
result, err := dice.RollSafe("2d6+5")
if err != nil {
    // handle error
}
```

## Expression Syntax

### Basic Dice

| Expression | Description |
|------------|-------------|
| `d6` | 1 six-sided die |
| `2d20` | 2 twenty-sided dice |
| `d10+5` | 1d10 plus 5 |
| `2d6-1` | 2d6 minus 1 |

### Modifiers (after colon)

| Modifier | Description | Example |
|----------|------------|---------|
| `dlN` | Drop lowest N dice | `3d6:dl1` |
| `dhN` | Drop highest N dice | `3d6:dh1` |
| `+Ne` | Add to each die | `2d6:+1e` |
| `/N` | Divide result by N | `2d6:/2` |
| `xN` or `*N` | Multiply result by N | `2d6:*2` |
| `+N>>M` | Add N to die at position M | `2d6:+1>>1` |

### Special Dice Types

| Type | Description | Status |
|------|------------|--------|
| `d` | Normal dice | ✅ Implemented |
| `D` | Concat dice (combine all rolls) | TODO |
| `DD` | Destructive dice (best N) | TODO |

### Full Expression Examples

```
2d6           -> 2 six-sided dice, sum
d20+5         -> 1d20 + 5
3d6:dl1       -> roll 3d6, drop lowest 1 (boon)
3d6:dh1       -> roll 3d6, drop highest 1 (bane)
2d6:+1e:/2    -> add 1 to each, sum, divide by 2
4d6:dl1:x100  -> roll 4d6, drop lowest, multiply by 100 (characteristic)
```

## Public API

### Roll

```go
func Roll(expression string) Result
```

Panics on invalid expression. Use for **hardcoded** expressions.

```go
result := dice.Roll("2d6")
```

### RollSafe

```go
func RollSafe(expression string) (Result, error)
```

Returns error for invalid expressions. Use for **user-provided** input.

```go
result, err := dice.RollSafe(userInput)
if err != nil {
    log.Printf("Invalid roll: %v", err)
}
```

## Result Structure

```go
type Result struct {
    Rolled         Dicepool // Original dicepool
    Raw            []int    // Individual die rolls
    Mods           []Mod    // Applied modifiers
    Final          []int    // Final values after modifiers
    FinalAsStrings []string // String representations (for UI)
}
```

## Architecture

### Files

| File | Purpose |
|------|---------|
| `spec.go` | Type definitions |
| `builders.go` | Constructors, Roller init |
| `mods.go` | Modifier implementations |
| `parse.go` | Expression parser |
| `roll.go` | Rolling logic |
| `api.go` | Public API |
| `cache.go` | Expression cache |

### Flow

```
Roll("2d6+5")
    │
    ├─► cache.Get() ──yes──► use cached Dicepool
    │                      │
    │          no ◄────────┘
    ▼
parseExpression("2d6+5")
    │
    ├─► parseDicePart() ──► []Die
    ├─► parseComplexMods() ──► []Mod
    └─► parseSimpleAdditive() ──► AddConst
    │
    ▼
new Dicepool(dice, modifiers)
    │
    └─► cache.Set() (for next time)
    ▼
Roller.rollDicepool()
    │
    ├─► roll each Die
    ├─► apply Modifiers in priority order
    └─► return Result
```

### Modifier Priority Order

Modifiers are applied in this order (lower priority = applied first):

| Priority | Modifier |
|----------|----------|
| 0 | None |
| 20 | AddIndividual |
| 30 | AddToEach |
| 70 | DropLowest |
| 71 | DropHighest |
| 100 | Sum |
| 110 | AddConst |
| 120 | Divide |
| 130 | Multiply |

## Error Handling

All modifiers return meaningful errors:

```go
// AddIndividual - position out of range
AddIndividual{position: 10}.Apply([]int{1,2,3})
// -> error: "position 10 out of range, have 3 dice"

// DropLowest - trying to drop more than available
DropLowest{quantity: 5}.Apply([]int{1,2})
// -> error: "cannot drop 5 dice from pool of 2"

// Divide - division by zero
Divide{value: 0}.Apply([]int{6})
// -> error: "division by zero"
```

## Internal Types

### Mod Interface

```go
type Mod interface {
    Apply([]int) ([]int, error)
    Priority() int
}
```

### Die

```go
type Die struct {
    Faces    int
    Codes    map[int]string    // e.g., 20 → "crit"
    Metadata map[string]string // color, name, tags
}
```

### Dicepool

```go
type Dicepool struct {
    Type      string
    Dice      []Die
    Modifiers []Mod
    Metadata  map[string]string
}
```

---

## Future Plans

### 1. Special Dice Types (TODO)

**Concat Dice (`D`)**: Combine all dice into single roll (not sum)
- `D6` → rolls 1d6, result is the single value (not summed)
- Used for characteristic generation

**Destructive Dice (`DD`)**: Roll N dice, keep best M
- Example: `2DD6` from Mongoose Traveller 2E
- Roll 2d6, keep highest (Equivalent to Advantage)

Implementation status: Parser recognizes `D` and `DD` (parse.go:58-63), but rolling is not implemented. Would need new modifier types to handle keeping best/highest.

### 2. Logger Integration

Add logging for dice rolls:

```go
type Logger interface {
    LogRoll(expr string, result Result)
    LogModifier(mod Mod, input, output []int)
    LogError(expr string, err error)
}
```

Benefits:
- Debug game sessions
- Audit player rolls
- Replay functionality

Implementation: Add logger interface to `spec.go`, inject via options pattern or context.

### 3. Probability Calculator

Calculate probabilities for dice expressions:

```go
type Calculator func(expression string) (ProbabilityDistribution, error)

type ProbabilityDistribution struct {
    Min     int
    Max     int
    Mean    float64
    Median  int
    Mode    int
    StdDev  float64
    Table   map[int]float64 // value → probability
}
```

Use cases:
- Character optimization
- Combat calculations
- Rule validation

Implementation approaches:
1. Brute force - roll all combinations (feasible for small dice pools)
2. Monte Carlo - random sampling for large pools
3. Analytical - mathematical formulas per modifier

### 4. Additional Features (Backlog)

- **Exploding dice** (`d6!`): Roll again on max value
- **Keeping dice**: Support for `D` / `DD` syntax
- **Rerolling**: Reroll on specific values
- **Advantage/Disadvantage**: Built-in boons/bane modifiers
- **Dice notation export**: Convert to standard notation

---

## Contributing

When adding new modifiers:
1. Implement `Mod` interface
2. Add constants for name and priority in `mods.go`
3. Add parser support in `parse.go` if expression-based
4. Add tests in `mods_test.go`
5. Ensure error messages are descriptive