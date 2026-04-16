# Ehex Package

Extended Hexadecimal encoding for Traveller UWP (Universal World Profile) codes.

## Overview

The `ehex` package implements the Traveller extended hexadecimal encoding scheme,
which maps integers 0–33 to characters `0-9A-HJ-NP-Z` (excluding `I` and `O` to
avoid confusion with digits `1` and `0`). It also provides a set of special
sentinel codes for unknown, masked, reserved, and extension values.

This encoding is used throughout the Cepheus Engine for UWP fields such as
Starport, Size, Atmosphere, Hydrographics, Population, Government, and Law Level.

## Import

```go
import "github.com/Galdoba/cepheus/internal/domain/engine/ehex"
```

## Quick Start

```go
// Encode an integer to Ehex
e := ehex.FromValue(15)
fmt.Println(e.Code()) // "F"

// Decode a character to Ehex
e = ehex.FromCode("Z")
fmt.Println(e.Value()) // 33

// Special values
fmt.Println(ehex.Unknown.Code()) // "?"
fmt.Println(ehex.Any.Value())    // -102

// Custom codes for application-specific extensions
custom := ehex.New(100, "@", "custom meaning")
fmt.Println(custom.Code())       // "@"
fmt.Println(custom.Value())      // 100
```

---

## Ehex Encoding

### Standard Mapping (0–33)

The standard encoding uses 34 characters, excluding `I` and `O`:

| Value | Code | | Value | Code | | Value | Code |
|-------|------|-|-------|------|-|-------|------|
| 0 | `0` | | 11 | `B` | | 22 | `N` |
| 1 | `1` | | 12 | `C` | | 23 | `P` |
| 2 | `2` | | 13 | `D` | | 24 | `Q` |
| 3 | `3` | | 14 | `E` | | 25 | `R` |
| 4 | `4` | | 15 | `F` | | 26 | `S` |
| 5 | `5` | | 16 | `G` | | 27 | `T` |
| 6 | `6` | | 17 | `H` | | 28 | `U` |
| 7 | `7` | | 18 | `J` | | 29 | `V` |
| 8 | `8` | | 19 | `K` | | 30 | `W` |
| 9 | `9` | | 20 | `L` | | 31 | `X` |
| 10 | `A` | | 21 | `M` | | 32 | `Y` |
| | | | | | | 33 | `Z` |

Note: `I` (value 18→`J`) and `O` (value 24→`P`) are skipped.

### Special Codes

| Constant | Code | Value | Description |
|----------|------|-------|-------------|
| `Unknown` | `?` | -101 | Value is unknown or not determined |
| `Any` | `*` | -102 | Any value is acceptable |
| `Invalid` | `!` | -103 | No valid value exists |
| `Default` | `#` | -104 | Use the default value |
| `Ignore` | `-` | -105 | This field should be ignored |
| `Reserved` | `&` | -106 | Reserved for future use |
| `Masked` | `~` | -107 | Value is masked or hidden |
| `Extension` | `>` | -108 | Extended value (requires further parsing) |
| `Placeholder` | `.` | -109 | Placeholder to be filled in later |

### Extended Codes

Additional codes decoded by `FromCode`:

| Code | Value | Meaning |
|------|-------|---------|
| `s` | 0 | Alias for zero (lowercase) |
| `r` | 0 | Alias for zero (lowercase) |

All unrecognized codes return `Unknown`.

---

## API

### Ehex Type

```go
type Ehex struct {
    // unexported fields
}
```

An `Ehex` is a value type (struct, not pointer). It is comparable — two
`Ehex` instances are equal if all three fields (`code`, `value`, `description`)
match.

#### Methods

```go
// Code returns the character string (e.g., "A", "Z", "?").
func (e Ehex) Code() string

// Value returns the integer value (0-33, or negative for special codes).
func (e Ehex) Value() int

// Description returns the human-readable description (may be empty).
func (e Ehex) Description() string

// String implements fmt.Stringer — returns the same as Code().
func (e Ehex) String() string

// WithDescription returns a copy with the description set.
// The original Ehex is unchanged (immutable).
func (e Ehex) WithDescription(desc string) Ehex
```

### Construction

```go
// FromValue encodes an integer. Returns Unknown for values outside 0-33.
// Special negative values (-101 to -109) are NOT mapped — use the
// predefined constants directly.
func FromValue(value int) Ehex

// FromCode decodes a character. Checks basic map first, then extended map.
// Returns Unknown for unrecognized codes.
func FromCode(code string) Ehex

// New creates a custom Ehex with arbitrary value, code, and optional description.
// Use this for application-specific codes outside the standard range.
func New(value int, code string, description ...string) Ehex
```

### Predefined Constants

```go
var (
    Unknown     Ehex  // "?"  value=-101
    Any         Ehex  // "*"  value=-102
    Invalid     Ehex  // "!"  value=-103
    Default     Ehex  // "#"  value=-104
    Ignore      Ehex  // "-"  value=-105
    Reserved    Ehex  // "&"  value=-106
    Masked      Ehex  // "~"  value=-107
    Extension   Ehex  // ">"  value=-108
    Placeholder Ehex  // "."  value=-109
)
```

---

## Usage Patterns

### UWP Field Encoding

```go
// Encode a starport value
starport := ehex.FromValue(5) // "5" → actually "5" for value 5
starport = starport.WithDescription("Good starport")

// Encode a special "unknown" starport
unknown := ehex.Unknown
unknown = unknown.WithDescription("Starport not surveyed")
```

### UWP Field Decoding

```go
// Parse a UWP character
ch := "A"
eh := ehex.FromCode(ch)
if eh == ehex.Unknown {
    fmt.Println("Unknown value")
} else {
    fmt.Printf("Value: %d\n", eh.Value())
}
```

### Custom Extended Codes

```go
// Define an application-specific code
custom := ehex.New(42, "@", "Custom world rule")
fmt.Println(custom) // "@"
fmt.Println(custom.Value()) // 42
```

### Equality Checking

```go
a := ehex.FromValue(10)
b := ehex.FromCode("A")
fmt.Println(a == b) // true

// WithDescription creates a distinct instance
c := a.WithDescription("letter A")
fmt.Println(a == c) // false (descriptions differ)
```

---

## Behavior Details

### FromValue

| Input Range | Output |
|-------------|--------|
| 0–33 | Standard Ehex with matching code |
| Outside 0–33 | `Unknown` |
| Special values (-101 to -109) | `Unknown` (not mapped) |

`FromValue` never returns the predefined special constants (`Any`, `Invalid`,
etc.). These must be accessed directly.

### FromCode

Lookup order:
1. Check `codeToValue` (basic 0–9, A–Z excl I/O)
2. Check `codeToValueExtended` (special symbols, lowercase aliases)
3. Return `Unknown`

Special codes (`?`, `*`, `!`, etc.) return the corresponding predefined
constant by exact value match.

### Zero Value

```go
var e ehex.Ehex
e.Code()        // ""
e.Value()       // 0
e.Description() // ""
e.String()      // ""
```

The zero value is a valid `Ehex` with empty code and zero value. It is not
equal to any predefined constant.

---

## Architecture

### Files

| File | Purpose |
|------|---------|
| `ehex.go` | Core types, encoding maps (`valueToCode`, `codeToValue`, `codeToValueExtended`), `FromValue`, `FromCode`, `New`, methods, predefined constants |
| `ehex_test.go` | Comprehensive tests: map completeness, uniqueness, round-trip, special codes, custom creation, equality, zero value |

### Internal Maps

| Map | Direction | Coverage |
|-----|-----------|----------|
| `valueToCode` | int → string | 0–33 |
| `codeToValue` | string → int | `0`–`9`, `A`–`H`, `J`–`N`, `P`–`Z` |
| `codeToValueExtended` | string → int | Special symbols + lowercase aliases |

The two maps are tested for bidirectional consistency:
`FromCode(code) → Value → FromValue → Code` must round-trip for all standard codes.

---

## Testing

The test suite covers:

| Test | What it verifies |
|------|-----------------|
| `TestValueToCodeMapping/completeness` | All values 0–33 present |
| `TestValueToCodeMapping/uniqueness` | No duplicate codes, all single-char |
| `TestCodeToValueConsistency` | Bidirectional map consistency |
| `TestFromValue/valid range 0-33` | Correct encoding for all values |
| `TestFromValue/out of range` | Returns `Unknown` for invalid values |
| `TestFromValue/special values` | Special values not returned by `FromValue` |
| `TestFromCode/basic codes` | All standard codes decode correctly |
| `TestFromCode/special symbols` | All 9 special codes return correct constants |
| `TestFromCode/extended aliases` | `s` and `r` map to value 0 |
| `TestFromCode/invalid codes` | Unrecognized codes return `Unknown` |
| `TestRoundTrip` | `value → code → value` and `code → value → code` |
| `TestSpecialConstants` | All 9 predefined constants have correct fields |
| `TestNew` | Custom code creation with/without description |
| `TestWithDescription` | Immutability — original unchanged |
| `TestZeroValueEhex` | Zero value has empty code, value 0 |
| `TestEhexEquality` | Equality/inequality for identical and different instances |

---

## Contributing

When adding new standard codes:

1. Add entries to `valueToCode` and `codeToValue` (keep them in sync)
2. Add a completeness test if the range expands beyond 0–33
3. Run `go test` to verify round-trip consistency

When adding special codes:

1. Define a package-level `Ehex` variable with the sentinel value
2. Add the code → value mapping in `codeToValueExtended`
3. Add a case in `FromCode`'s switch to return the predefined constant
4. Add a test in `TestSpecialConstants` and `TestFromCode/special symbols`
