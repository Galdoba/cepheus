# Cepheus — Project Status

**Date:** 2026-04-15
**Module:** `github.com/Galdoba/cepheus`
**Go version:** 1.25.1
**Build:** ✅ `go build ./...`, `go vet ./...`, `go test ./...` — all clean

---

## Overview

Cepheus is a Go library implementing the **Cepheus Engine** RPG rules — a Traveller-derived sci-fi tabletop system. The project provides:

1. A **dice-rolling engine** with expression parsing, modifier chaining, caching, and special RPG mechanics (D66, Flux, Variance)
2. A **generic table lookup system** with index parsing, validation, cascading rolls, and JSON persistence
3. An **extended hex (Ehex) encoder** for Traveller UWP codes (0–9, A–Z excluding I/O) with special sentinel values
4. A **procedural star-system generator** (domain types and data tables defined; generation logic in progress)

---

## Package Status

### ✅ `internal/domain/engine/dice` — COMPLETE

| Component | Status | Notes |
|-----------|--------|-------|
| Expression parser (`parse.go`) | ✅ | Supports `NdN`, `:dlN`, `:dhN`, `:+Ne`, `:+N>>M`, `/N`, `*N`/`xN`, simple additive `+N`/`-N` |
| Modifier pipeline (`mods.go`) | ✅ | 8 modifiers with priority ordering (20–130) |
| Roller (`roll.go`, `roller.go`) | ✅ | `randRoller` via `math/rand`; seedable from string via `stringToInt64` hash |
| Manager (`manager.go`) | ✅ | `Roll`, `MustRoll`, `D66`, `Flux`, `FluxGood`, `FluxBad`, `Variance`; thread-safe via `sync.Mutex` |
| Cache (`cache.go`) | ✅ | Thread-safe expression cache (`sync.RWMutex`), parsed expressions reused |
| Die / Dicepool (`die.go`) | ✅ | `die` with codes/metadata, `dicepool` builder pattern |
| API (`api.go`) | ✅ | Package-level convenience wrappers + `Manager` methods |
| `ValidateExpression` | ✅ | Public wrapper around `parseExpression` for table validation |
| Tests (`spec_test.go`) | ✅ | 16 test functions — parsing, rolling, D66, Flux, Variance, concurrency, cache, edge cases |

**File inventory (10 files):**

| File | Lines | Purpose |
|------|-------|---------|
| `spec.go` | Core types | `Manager`, `Result`, `Roller` interface, default init |
| `api.go` | Public API | `Roll`, `MustRoll`, `D66`, `Flux`, `FluxGood`, `FluxBad`, `Variance` |
| `parse.go` | Parser | `parseExpression`, modifier parsing, `ValidateExpression` |
| `mods.go` | Modifiers | `addToEach`, `addIndividual`, `dropLowest`, `dropHighest`, `divide`, `multiply`, `addConst`, `summ` |
| `roll.go` | Rolling | `basicRoll` — rolls every die in a dicepool |
| `manager.go` | Orchestration | `Manager.roll`, `stdInterpreter` (applies mods to raw roll) |
| `roller.go` | RNG | `randRoller`, `stringToInt64` seed hashing, `randomSeed` |
| `cache.go` | Cache | `expressionCache` with `sync.RWMutex` |
| `die.go` | Types | `die`, `dicepool` with builder methods |
| `spec_test.go` | Tests | 16 test functions |

**TODO (backlog):**
- `D` (concat) and `DD` (destructive) dice types — parser recognizes them, rolling returns error
- Exploding dice (`d6!`), rerolls, advantage/disadvantage
- Logger interface for roll auditing
- Probability calculator (distribution analysis)

---

### ✅ `internal/domain/engine/tables` — COMPLETE

| Component | Status | Notes |
|-----------|--------|-------|
| `GameTable` | ✅ | Name, expression, data map, auto-D66 detection |
| `GameTable.Validate()` | ✅ | 9 checks: name, data count, expression, index parseability, duplicates, range holes, bounds, sentinel leakage, empty values |
| `Collection` | ✅ | Multi-table orchestration, duplicate name detection in constructor |
| `Collection.Roll()` | ✅ | D66 and standard rolls; error messages include table name, numeric index, D66 string |
| `Collection.RollCascade()` | ✅ | Chains rolls across tables; max depth 1000 (loop detection) |
| `Collection.Reset()` | ✅ | Clears roll history for reuse |
| `Collection.Validate()` | ✅ | Delegates to per-table validation |
| Index parsing | ✅ | Range (`1-10`), open-ended (`11+`, `5-`), single, comma-separated, mixed; deduplication |
| `Save` / `Load` | ✅ | JSON persistence with `json.MarshalIndent` |
| `TableRoller` interface | ✅ | Decoupled from dice engine — `D66() string`, `Roll(string, ...int) (int, error)` |
| Tests (`table_test.go`) | ✅ | 9 test functions — validation, parsing, collection ops, cascade, reset |

**Validation coverage:**

| Check | Status | Notes |
|-------|--------|-------|
| Empty table name | ✅ | |
| Minimum 2 data entries | ✅ | |
| Empty data values | ✅ | |
| Non-D66: expression parseable | ✅ | Regex `^(\d*)d(\d+)([+-]\d+)?$` |
| D66: expression is `d66`/`D66` | ✅ | Auto-detected in `New()` |
| Index string parseability | ✅ | All keys parse via `stringToIndexes` |
| Duplicate indexes across keys | ✅ | Set-based — catches overlapping ranges like `"1-3"` + `"2-4"` |
| Range holes (non-D66) | ✅ | Expanded index range must be contiguous |
| Index bounds [-1000, 1000] | ✅ | |
| Sentinel value leakage (`andAbove`/`andBelow`) | ✅ | |
| Duplicate table names in collection | ✅ | Detected in `NewCollection` |

**Deferred by design:**

| Item | Rationale |
|------|-----------|
| Expression range vs index overlap | Mods can push roll results beyond expression range; unreachable indexes are special cases |
| D66 strict index bounds (11–66) | Tables with 00–99 range + per-die mods are valid use cases |
| Cross-table reference validation | Deferred until value semantics and `RollCascade` design are finalized |
| Configurable bounds per table | ±1000 is the safety overhead; per-table bounds not needed |

**File inventory (4 files):**

| File | Lines | Purpose |
|------|-------|---------|
| `table.go` | GameTable | `New`, `Validate`, `stringToIndexes`, `indexesToString`, `Save`/`Load` |
| `collection.go` | Collection | `NewCollection`, `Roll`, `RollCascade`, `Reset`, `Validate` |
| `roller.go` | Interface | `TableRoller` — dice abstraction |
| `table_test.go` | Tests | 9 test functions, mock roller |

---

### ✅ `internal/domain/engine/ehex` — COMPLETE

| Component | Status | Notes |
|-----------|--------|-------|
| `Ehex` type | ✅ | Immutable struct: `code` (string), `value` (int), `description` (string) |
| Standard encoding (0–33) | ✅ | `0-9A-HJ-NP-Z` — excludes I and O |
| `FromValue` | ✅ | int → Ehex; returns `Unknown` for values outside 0–33 |
| `FromCode` | ✅ | string → Ehex; checks basic map, then extended map; returns `Unknown` |
| `New` | ✅ | Custom Ehex with arbitrary value/code, optional description |
| `WithDescription` | ✅ | Immutable copy — original unchanged |
| 9 special constants | ✅ | `Unknown`, `Any`, `Invalid`, `Default`, `Ignore`, `Reserved`, `Masked`, `Extension`, `Placeholder` |
| Extended aliases | ✅ | `s` → 0, `r` → 0 (lowercase zero aliases) |
| `fmt.Stringer` | ✅ | `String()` returns `Code()` |
| Comparable | ✅ | Struct equality on all three fields |
| Tests (`ehex_test.go`) | ✅ | 10 test functions — completeness, uniqueness, consistency, round-trip, constants, custom creation, immutability, zero value, equality |

**File inventory (2 files):**

| File | Lines | Purpose |
|------|-------|---------|
| `ehex.go` | Core | `Ehex` type, encoding maps, `FromValue`, `FromCode`, `New`, methods, constants |
| `ehex_test.go` | Tests | 10 test functions |

---

### 🚧 `internal/domain/engine/systemgen` — PARTIALLY IMPLEMENTED

#### Types (`types.go`) — ✅ COMPLETE

All domain types and enums defined:
- `StarSystem`, `CelestialBody`, `StarData`, `BrownDwarfData`, `RoguePlanetData`
- `CompanionStar`, `GasGiant`, `AsteroidBelt`, `RockyPlanet`, `Moon`
- `ZimmPoint`, `SystemQuirk`
- All enums: `ObjectType`, `StellarClass`, `LuminosityClass`, `BrownDwarfClass`, `PlanetType`, `CompanionDistance`, `GasGiantSize`, `ZimmPointType`
- Config types: `SubsectorType`, `RealismMode`, `PlacementMethod`
- `Generator` struct with `New()`, `GenerateHex()`, `GenerateSubsector()` stubs
- 18 generation step function signatures declared (all `TODO` bodies)
- `GenerationResult` struct with `System`, `Log`, `Warnings`

#### Tables (`tables.go`) — ⚠️ INCOMPLETE

| Table | Status | Notes |
|-------|--------|-------|
| Object presence thresholds | ✅ | All 6 subsector densities (Empty → Dense) |
| Object type ranges (d100) | ✅ | Star → Black Hole |
| Brown dwarf class ranges | ✅ | L / T / Y |
| Stellar class — Realistic | ✅ | d100 → O–M |
| Stellar class — Semi-realistic | ✅ | d100 → O–M |
| Stellar class — Fantastic | ✅ | d100 → O–M |
| Luminosity class (d100 + d10 sub-roll) | ✅ | V, IV, D, III, II, VI, Ia, Ib |
| Multiple star system (OBO / KGF / M) | ✅ | 3 tables by stellar class |
| Companion distance (Contact–Distant) | ✅ | d100 + 4 AU sub-tables |
| Main Sequence zones (V) | ✅ | O0–O9, B0–B9, A0–A9, F0–F9, G0–G7, K0–K9, M0–M9 |
| Brown Dwarf zones | ✅ | L0–L9, T0–T9, Y0–Y9 — 30 entries |
| Gas giant size tables | ✅ | Neptunian (2d6), Jovian (2d10) |
| System quirk table (d66) | ✅ | 11–66 — all 36 entries |
| **Luminous Supergiant zones (Ia)** | ❌ | O0–M9 — ~90 entries |
| **Less Luminous Supergiant zones (Ib)** | ❌ | O0–M8 — ~89 entries |
| **Bright Giant zones (II)** | ❌ | O0–M9 — ~90 entries |
| **Giant zones (III)** | ❌ | O0–M9 — ~90 entries |
| **Subgiant zones (IV)** | ❌ | O0–M9 — ~90 entries |
| **Subdwarf zones (VI)** | ❌ | O0–M9 — ~90 entries |
| **White Dwarf zones (D)** | ❌ | D0–D9 — 10 entries |

**Missing zone tables summary:**
- **~539 star zone entries** across 7 luminosity classes (Ia, Ib, II, III, IV, VI, D)
- Data exists in `assets/generic_starzone_*.json` files — can be loaded dynamically instead of hardcoded
- Each entry: `Name`, `Temperature`, `Mass`, `Luminosity`, `InnerLimit`, `HabitableMin`, `HabitableMax`, `SnowLine`, `OuterLimit`

**Tables still needed (not in codebase or assets):**
- Rocky planet type by orbital zone (2d6 tables for inner/habitable/outer)
- Planet size/atmosphere/hydrographics tables
- UWP (Universal World Profile) generation — deferred to separate package per plan
- Moon generation tables
- Eccentricity table (2d10 → 0.002–0.250)
- Tech level table (1d6 → 10–12)

#### Generation Steps — ❌ NOT IMPLEMENTED

All 18 step functions declared with `TODO` bodies in `types.go`:

| Step | Function | Purpose |
|------|----------|---------|
| 1 | `step1` | Determine object presence in hex |
| 2 | `step2` | Determine object type (Star/BrownDwarf/etc.) |
| 3 | `step3` | Brown dwarf subtype (L/T/Y) |
| 4 | `step4` | Stellar class (Realistic/Semi-realistic/Fantastic) |
| 5 | `step5` | Numeric classification (0–9) |
| 6 | `step6` | Luminosity class |
| 7 | `step7` | Multiple star system determination |
| 8 | `step8` | Companion star distances |
| 9 | `step9` | Gas giant count |
| 10 | `step10` | Asteroid belt count |
| 11 | `step11` | Habitable zone lookup |
| 12 | `step12` | Main world placement |
| 13 | `step13` | Gas giant placement + migration |
| 14 | `step14` | Orbital placement (doubling or varied) |
| 15 | `step15` | Rocky planet details |
| 15A–C | `step15A/B/C` | Planet detail by zone (inner/habitable/outer) |
| 16 | `step16` | Moon generation |
| 17 | `step17` | Zimm Point placement |
| 18 | `step18` | System quirks |

**Missing files (per plan):**
- `generator.go` — main orchestration logic (currently stubs in `types.go`)
- `steps.go` — step implementations
- `roller.go` — dice adapter (would wrap `dice.Roll`/`dice.MustRoll`)
- `errors.go` — custom error types

---

## Documentation

| File | Status | Covers |
|------|--------|--------|
| `doc/dice.md` | ✅ Updated 2026-04-15 | Expression syntax, modifier priority, full API, architecture, error handling, future plans |
| `doc/tables.md` | ✅ New 2026-04-15 | GameTable, Collection, RollCascade, index parsing, TableRoller interface, JSON persistence |
| `doc/ehex.md` | ✅ New 2026-04-15 | Standard encoding (0–33), special constants, API, usage patterns, test coverage |
| `doc/project_status.md` | This file | Overall project status |
| `doc/System_Generation_Extended.md` | Reference | Source rules for system generation algorithm |
| `doc/system_generator_plan.md` | Reference | Implementation plan for systemgen |

---

## Assets

**127 JSON files** in `assets/` covering:

| Category | Files | Content |
|----------|-------|---------|
| Star zones | 8 | `generic_starzone_V.json`, `Ia`, `Ib`, `II`, `III`, `IV`, `VI`, `D`, `BD` |
| Planet types | 5 | Dwarf, Mercurian, Subterran, Terran, Superterran size/atmosphere tables |
| Hydrographics | 6 | Modifiers, zone-specific hydro tables per planet type |
| Government/Tech | 10+ | Government codes, tech levels, population modifiers |
| Moons | 6 | Moon count, orbit, size tables (by primary type) |
| Asteroids | 5 | Belt count, width, composition, diameter, rings |
| Gas giants | 3 | Count, migration distance, moon orbit |
| Formulas | 8 | Blackbody temp, greenhouse, gravity, mass, orbital period, rotation, tidal lock, polar/equatorial/avg temp |
| Special | 10+ | Biology modifiers, atmosphere types (corrosive/insidious), albedo tables, axial tilt, quirks, UWP-related |

**Design decision pending:** Should `systemgen` load tables from JSON assets dynamically, or hardcode them in Go (as currently done for Main Sequence V + Brown Dwarfs)?

---

## Dependencies

| Package | Usage |
|---------|-------|
| `internal/domain/engine/dice` | All random rolls in `tables` and `systemgen` |
| `internal/domain/engine/ehex` | Standalone — no intra-project dependencies |
| `internal/domain/engine/tables` | Depends on `dice` only via `TableRoller` interface |
| Standard library | `fmt`, `math/rand`, `strconv`, `strings`, `slices`, `sort`, `sync`, `time`, `regexp`, `errors`, `encoding/json`, `os`, `path/filepath` |

No external dependencies.

---

## Build & Test Status

```
go build ./...     ✅ clean
go vet ./...       ✅ clean
go test ./...      ✅ all pass (35 tests total)
  dice             ✅ 16 tests — 0.004s
  ehex             ✅ 10 tests — 0.002s
  tables           ✅  9 tests — 0.029s
  systemgen        [no test files]
```

---

## Next Steps (Priority Order)

1. **Complete star zone tables** — Load from `assets/generic_starzone_*.json` or hardcode Ia, Ib, II, III, IV, VI, D (~539 entries)
2. **Add planet detail tables** — Rocky planet types, size/atmosphere/hydrographics, eccentricity, tech level
3. **Implement generation steps 1–11** — Deterministic, table-lookup driven
4. **Implement steps 12–18** — Placement logic, migration, Zimm Points, quirks
5. **Extract `generator.go` and `steps.go`** from current `types.go` stubs
6. **Add `errors.go`** with domain-specific error types
7. **Add `roller.go`** dice adapter for systemgen
8. **Write tests** — table verification, step unit tests (mock dice rolls), integration tests with known seeds
9. **Decide on asset loading strategy** — dynamic JSON load vs. hardcoded Go tables
