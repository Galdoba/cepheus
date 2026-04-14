# System Generation Implementation Plan

## Overview

Implement the System Generation algorithm from `System_Generation_Extended.md` as a Go package under `internal/domain/engine/systemgen/`. The generator produces procedurally generated star systems following the 18-step RPG ruleset, using the existing dice package for all random rolls.

---

## Architecture

### Package Structure

```
internal/domain/engine/systemgen/
├── types.go          # Core domain types (StarSystem, Star, Planet, etc.)
├── tables.go         # Static data tables (stellar classes, zones, etc.)
├── generator.go      # Main Generator struct and public API
├── steps.go          # Individual generation step implementations (1-18)
├── roller.go         # Dice rolling adapter (wraps dice.RollSafe)
└── errors.go         # Custom error types
```

### Design Principles

1. **Dependency on dice package**: All random rolls go through `dice.Roll()` or `dice.RollSafe()` — no direct `math/rand` usage
2. **Immutable data tables**: Lookup tables are read-only, defined as package-level variables or constants
3. **Step-by-step generation**: Each generation step is a separate function for testability
4. **Seedable RNG**: The Generator accepts a string seed, converted to int64 via the dice package's seeding mechanism
5. **No side effects**: Generation functions take input, return output — no global state mutation

---

## Implementation Steps

### Phase 1: Domain Types (`types.go`)

**Core Types:**

```go
type StarSystem struct {
    HexCoords   HexCoords
    ObjectType  ObjectType      // Star, BrownDwarf, RoguePlanet, etc.
    PrimaryBody CelestialBody
    Companions  []CompanionStar
    GasGiants   []GasGiant
    AsteroidBelts []AsteroidBelt
    RockyPlanets  []RockyPlanet
    ZimmPoints    []ZimmPoint
    Quirks        []SystemQuirk
}

type ObjectType int
const (
    ObjectStar ObjectType = iota
    ObjectBrownDwarf
    ObjectRoguePlanet
    ObjectRogueGasGiant
    ObjectNeutronStar
    ObjectNebula
    ObjectBlackHole
)

type CelestialBody struct {
    Star          *StarData
    BrownDwarf    *BrownDwarfData
    RoguePlanet   *RoguePlanetData
}

type StarData struct {
    StellarClass    StellarClass   // O, B, A, F, G, K, M
    NumericClass    int            // 0-9
    LuminosityClass LuminosityClass // V, IV, III, II, Ia, Ib, VI, D
    Temperature     int            // Kelvin
    Mass            float64        // Solar masses
    Luminosity      float64        // Solar units
    InnerLimit      float64        // AU
    HabitableZone   ZoneRange      // Inner, Outer AU
    SnowLine        float64        // AU (0 = N/A)
    OuterLimit      float64        // AU
}

type BrownDwarfData struct {
    Class      BrownDwarfClass // L, T, Y
    NumericClass int
    // ... same zone data as StarData
}

type StellarClass int
const (
    ClassO StellarClass = iota
    ClassB
    ClassA
    ClassF
    ClassG
    ClassK
    ClassM
)

type LuminosityClass int
const (
    LumV LuminosityClass = iota  // Main Sequence
    LumIV                         // Subgiant
    LumIII                        // Giant
    LumII                         // Bright Giant
    LumIa                         // Luminous Supergiant
    LumIb                         // Supergiant
    LumVI                         // Subdwarf
    LumD                          // Dwarf (white dwarf)
)

type ZoneRange struct {
    Inner float64
    Outer float64
}

type CompanionStar struct {
    StarData StarData
    Distance CompanionDistance // Contact, Close, Near, Far, Distant
    AUs      float64
}

type CompanionDistance int
const (
    DistContact CompanionDistance = iota
    DistClose
    DistNear
    DistFar
    DistDistant
)

type GasGiant struct {
    SizeType GasGiantSize // Neptunian or Jovian
    SizeCode int
    OrbitAU  float64
    Migrated bool
    OriginalOrbitAU float64
}

type GasGiantSize int
const (
    SizeNeptunian GasGiantSize = iota
    SizeJovian
)

type AsteroidBelt struct {
    OrbitAU float64
    Width   float64
}

type RockyPlanet struct {
    PlanetType      PlanetType    // Dwarf, Mercurian, Subterran, Terran, Superterran
    SizeCode        int
    AtmosphereCode  int
    Hydrographics   int
    Population      int
    Government      int
    LawLevel        int
    TechLevel       int
    OrbitAU         float64
    Eccentricity    float64
    IsMainWorld     bool
    Moons           []Moon
}

type PlanetType int
const (
    TypeDwarf PlanetType = iota
    TypeMercurian
    TypeSubterran
    TypeTerran
    TypeSuperterran
)

type Moon struct {
    SizeCode       int
    AtmosphereCode int
    // ... similar to RockyPlanet
}

type ZimmPoint struct {
    Type       ZimmPointType // Incoming or Outgoing
    OrbitAU    float64
    TargetName string
}

type ZimmPointType int
const (
    ZPIncoming ZimmPointType = iota
    ZPOutgoing
)

type SystemQuirk struct {
    RollResult int
    Description string
}

type HexCoords struct {
    Row int
    Col int
}
```

### Phase 2: Data Tables (`tables.go`)

**Tables to implement:**

1. **SubsectorType → ObjectPresence threshold** (Empty=5, Scattered=20, etc.)
2. **ObjectType lookup** (d100 → Star/BrownDwarf/RoguePlanet/etc.)
3. **BrownDwarfClass lookup** (d100 → L/T/Y)
4. **StellarClass tables** (Realistic, Semi-realistic, Fantastic)
5. **LuminosityClass lookup** (d100 → V/IV/III/etc.)
6. **MultipleStarSystem tables** (by stellar class → Single/Binary/Trinary/etc.)
7. **CompanionDistance tables** (d100 → Contact/Close/Near/Far/Distant + AU sub-tables)
8. **StarZoneData table** — the massive table mapping every star type to Temperature/Mass/Luminosity/InnerLimit/HabitableZone/SnowLine/OuterLimit
   - Main Sequence (V): O0-M9
   - Luminous Supergiant (Ia): O0-M9
   - Less Luminous Supergiant (Ib): O0-M8
   - Bright Giant (II): O0-M9
   - Giant (III): O0-M9
   - Subgiant (IV): O0-M9
   - Supergiant (Ib duplicate?): O0-M9
   - Dwarf (D): D0-D9
   - Brown Dwarfs: L0-Y9

**Table structure:**

```go
type StarZoneEntry struct {
    Name         string  // e.g., "O0 V", "G2 V", "L3"
    Temperature  int
    Mass         float64
    Luminosity   float64
    InnerLimit   float64
    HabitableMin float64
    HabitableMax float64
    SnowLine     float64 // 0 = N/A
    OuterLimit   float64
}

var MainSequenceZones = map[string]StarZoneEntry{...}
var SupergiantIaZones = map[string]StarZoneEntry{...}
// ... etc for each luminosity class
var BrownDwarfZones = map[string]StarZoneEntry{...}
```

### Phase 3: Generator API (`generator.go`)

**Public API:**

```go
type Generator struct {
    subsectorType SubsectorType
    seed          string
    // internal state
}

type SubsectorType int
const (
    SubEmpty SubsectorType = iota
    SubScattered
    SubDispersed
    SubAverage
    SubCrowded
    SubDense
)

type GenerationResult struct {
    System  StarSystem
    Log     []string  // step-by-step generation log
    Warnings []string // non-fatal issues (e.g., planet outside inner limit)
}

// New creates a Generator with the given subsector density and RNG seed
func New(subsectorType SubsectorType, seed string) *Generator

// GenerateHex generates a complete star system for the given hex coordinates
func (g *Generator) GenerateHex(row, col int) (*GenerationResult, error)

// GenerateSubsector generates all hexes in a subsector (typically 8x5 = 40 hexes)
func (g *Generator) GenerateSubsector() ([]*GenerationResult, error)
```

### Phase 4: Generation Steps (`steps.go`)

Each step corresponds to a numbered step in the document:

```go
// step1: Determine object presence in hex
func step1(g *Generator) (bool, error)

// step2: Determine object type (Star/BrownDwarf/etc.)
func step2(g *Generator) (ObjectType, error)

// step3: Brown dwarf subtype (L/T/Y)
func step3(g *Generator) (BrownDwarfClass, error)

// step4: Star type (Realistic/Semi-realistic/Fantastic)
func step4(g *Generator, realism RealismMode) (StellarClass, error)

// step5: Numeric classification (0-9)
func step5(g *Generator) (int, error)

// step6: Luminosity class
func step6(g *Generator) (LuminosityClass, error)

// step7: Multiple star system determination
func step7(g *Generator, primary StellarClass) (int, error) // returns number of stars

// step8: Companion star distances
func step8(g *Generator, numCompanions int) ([]CompanionStar, error)

// step9: Gas giant count
func step9(g *Generator, primaryType CelestialBody) (int, error)

// step10: Asteroid belt count
func step10(g *Generator) (int, error)

// step11: Habitable zone lookup (not a roll, just table lookup)
func step11(body CelestialBody) (StarZoneEntry, error)

// step12: Main world placement (user-provided or skip)
func step12(g *Generator, mainWorld *RockyPlanet) error

// step13: Gas giant placement + migration
func step13(g *Generator, system *StarSystem, zones StarZoneEntry) error

// step14: Orbital placement (doubling or varied distance)
func step14(g *Generator, system *StarSystem, zones StarZoneEntry, placementMethod PlacementMethod) error

// step15: Rocky planet details
func step15(g *Generator, system *StarSystem, zones StarZoneEntry) error

// step15A-C: Planet detail by zone (inner/habitable/outer)
func step15A(g *Generator, planet *RockyPlanet) error
func step15B(g *Generator, planet *RockyPlanet, zones StarZoneEntry) error
func step15C(g *Generator, planet *RockyPlanet) error

// step16: Moon generation
func step16(g *Generator, planet *RockyPlanet) error

// step17: Zimm Point placement
func step17(g *Generator, system *StarSystem, adjacentSystems []string) error

// step18: System quirks
func step18(g *Generator) ([]SystemQuirk, error)
```

### Phase 5: Dice Roller Adapter (`roller.go`)

```go
// Roll is a convenience wrapper around dice.Roll
func Roll(expression string) (int, error)

// MustRoll panics on error (for internal use where expression is hardcoded)
func MustRoll(expression string) int
```

### Phase 6: Error Types (`errors.go`)

```go
var (
    ErrInvalidExpression  = errors.New("invalid dice expression")
    ErrNoObjectInHex      = errors.New("no object present in this hex")
    ErrBlackHole          = errors.New("black hole — no further generation needed")
    ErrNoHabitableZone    = errors.New("star has no habitable zone")
    ErrPlanetOutOfBounds  = errors.New("planet orbit is outside valid range")
)
```

---

## Dependencies

- `internal/domain/engine/dice` — all random rolls
- Standard library only (no external dependencies)

---

## Testing Strategy

1. **Table tests** for all lookup tables — verify every entry matches the document
2. **Unit tests per step** — mock dice rolls to test deterministic behavior
3. **Integration tests** — full generation runs with known seeds, verify output structure
4. **Edge case tests** — black holes, neutron stars, no habitable zone, migrated gas giants

---

## Implementation Order

1. `types.go` — define all domain types
2. `tables.go` — populate data tables (tedious but straightforward)
3. `errors.go` — define error types
4. `roller.go` — dice adapter
5. `generator.go` — Generator struct and public API
6. `steps.go` — implement steps 1-18 in order
7. Tests — add tests incrementally per step

---

## Notes

- The document references external rules (Clement Sector: The Rules) for government, law level, tech level, and UWP generation. These should be stubbed or implemented as separate future packages
- Zimm Points require knowledge of adjacent systems — this will be an input parameter
- Planet placement has two methods (doubling vs varied distance) — both should be supported via a `PlacementMethod` parameter
- Star zone data tables are very large — consider generating them programmatically or storing in a separate data file
- Temperature calculation formulas (p.136-144) are optional detail — implement as a separate `climate` sub-package if needed
