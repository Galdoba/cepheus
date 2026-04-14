# Cepheus — Project Status

**Date:** 2026-04-14
**Module:** `github.com/Galdoba/cepheus`
**Go version:** 1.25.1
**Build:** ✅ `go build ./...`, `go vet ./...`, `go test ./...` — all clean

---

## Overview

Cepheus is a Go library implementing the **Cepheus Engine** RPG rules — a Traveller-derived sci-fi tabletop system. The project provides a dice-rolling engine with expression parsing, a generic game-table lookup system, and a procedural star-system generator.

---

## Package Status

### ✅ `internal/domain/engine/dice` — COMPLETE

| Component | Status | Notes |
|-----------|--------|-------|
| Expression parser (`parse.go`) | ✅ | Supports `NdN`, `:dlN`, `:dhN`, `:+Ne`, `:+N>>M`, `/N`, `*N`/`xN` |
| Modifier pipeline (`mods.go`) | ✅ | 8 modifiers with priority ordering (0–130) |
| Roller (`roll.go`, `spec.go`) | ✅ | Seedable RNG via `rand.New(rand.NewSource(...))` |
| Manager (`manager.go`) | ✅ | `RollSafe`, `Roll`, `D66` methods; `Result` struct with per-die tracking |
| Cache (`cache.go`) | ✅ | Thread-safe expression cache (sync.RWMutex) |
| Builders (`builders.go`) | ✅ | Fluent constructors, `stringToInt64` seed hashing |
| `ValidateExpression` | ✅ | Public wrapper around `parseExpression` for table validation |
| Tests | ✅ | `manager_test.go` — all pass |

**TODO (backlog):**
- `D` (concat) and `DD` (destructive) dice types — parser recognizes them, rolling not implemented
- Exploding dice, rerolls, advantage/disadvantage (documented in future plans)
- Logger interface for roll auditing

---

### ✅ `internal/domain/engine/tables` — VALIDATION IMPROVED

| Component | Status | Notes |
|-----------|--------|-------|
| `GameTable` | ✅ | Name, expression, data map, D66 support |
| `GameTable.Validate()` | ✅ | Comprehensive: name, data count, expression/D66 consistency, duplicate indexes, range holes, bounds, empty values |
| `Collection` | ✅ | Multi-table orchestration with duplicate name detection in constructor |
| `NewCollection()` | ✅ | Validates during construction — rejects duplicate table names |
| `Collection.Validate()` | ✅ | Delegates to per-table validation |
| `Collection.Roll()` | ✅ | Improved error message: reports table name, numeric index, and D66 index string |
| Index parsing | ✅ | Range (`1-10`), open-ended (`11+`, `5-`), single values, deduplication |
| Tests | ✅ | `table_test.go` — all 29 tests pass |

**Validation coverage (current):**

| Check | Status | Notes |
|-------|--------|-------|
| Empty table name | ✅ | |
| Minimum 2 data entries | ✅ | |
| Empty data values | ✅ | |
| Non-D66: expression parseable | ✅ | Uses `dice.ValidateExpression` |
| D66: expression must be empty | ✅ | `WithD66Roll(true)` clears expression |
| Index string parseability | ✅ | All keys parse via `stringToIndexes` |
| Duplicate indexes across keys | ✅ | Set-based detection — catches overlapping ranges like `"1-3"` + `"2-4"` |
| Range holes (non-D66) | ✅ | Expanded index range must be contiguous |
| Index bounds [-1000, 1000] | ✅ | |
| Sentinel value leakage (`andAbove`/`andBelow`) | ✅ | |
| Duplicate table names in collection | ✅ | Detected in `NewCollection` |

**Deferred by design:**

| Item | Rationale |
|------|-----------|
| Expression range vs index overlap | Mods can push roll results beyond expression range; unreachable indexes are special cases for manual `Roll` |
| D66 strict index bounds (11–66) | Tables with 00–99 range + per-die mods are valid use cases |
| Cross-table reference validation | Deferred until value semantics and `RollCascade` design are finalized |
| Configurable bounds per table | ±1000 is the safety overhead; per-table bounds not needed |

**TODO:**
- `RollCascade` implementation — chain rolls: if a table's result value matches another table name, re-roll on that table until reaching a terminal value

---

### 🚧 `internal/domain/engine/systemgen` — PARTIALLY IMPLEMENTED

#### Types (`types.go`) — ✅ COMPLETE

All domain types defined:
- `StarSystem`, `CelestialBody`, `StarData`, `BrownDwarfData`, `RoguePlanetData`
- `CompanionStar`, `GasGiant`, `AsteroidBelt`, `RockyPlanet`, `Moon`
- `ZimmPoint`, `SystemQuirk`
- All enums: `ObjectType`, `StellarClass`, `LuminosityClass`, `BrownDwarfClass`, `PlanetType`, `CompanionDistance`, `GasGiantSize`, `ZimmPointType`
- Config types: `SubsectorType`, `RealismMode`, `PlacementMethod`
- `Generator` struct with `New()`, `GenerateHex()`, `GenerateSubsector()` stubs
- 18 generation step function signatures declared (all `TODO` bodies)

#### Tables (`tables.go`) — ⚠️ INCOMPLETE

| Table | Status | Notes |
|-------|--------|-------|
| Object presence thresholds | ✅ | All 6 subsector densities |
| Object type ranges (d100) | ✅ | Star → Black Hole |
| Brown dwarf class ranges | ✅ | L / T / Y |
| Stellar class — Realistic | ✅ | d100 → O–M |
| Stellar class — Semi-realistic | ✅ | d100 → O–M |
| Stellar class — Fantastic | ✅ | d100 → O–M |
| Luminosity class (d100 + d10 sub-roll) | ✅ | V, IV, D, III, II, VI, Ia, Ib |
| Multiple star system (OBO / KGF / M) | ✅ | 3 tables by stellar class |
| Companion distance (Contact–Distant) | ✅ | d100 + 4 AU sub-tables |
| **Main Sequence zones (V)** | ✅ | O0–O9, B0–B9, A0–A9, F0–F9, G0–G7, K0–K9, M0–M9 |
| **Luminous Supergiant zones (Ia)** | ❌ MISSING | O0–M9 — ~90 entries in source doc |
| **Less Luminous Supergiant zones (Ib)** | ❌ MISSING | O0–M8 — ~89 entries |
| **Bright Giant zones (II)** | ❌ MISSING | O0–M9 — ~90 entries |
| **Giant zones (III)** | ❌ MISSING | O0–M9 — ~90 entries |
| **Subgiant zones (IV)** | ❌ MISSING | O0–M9 — ~90 entries |
| **Subdwarf zones (VI)** | ❌ MISSING | O0–M9 — ~90 entries |
| **White Dwarf zones (D)** | ❌ MISSING | D0–D9 — 10 entries |
| Brown Dwarf zones | ✅ | L0–L9, T0–T9, Y0–Y9 — 30 entries |
| Gas giant size tables | ✅ | Neptunian (2d6), Jovian (2d10) |
| System quirk table (d66) | ✅ | 11–66 — all 36 entries |

**Missing zone tables summary:**
- **~539 star zone entries** across 7 luminosity classes (Ia, Ib, II, III, IV, VI, D)
- All are present in `System_Generation_Extended.md` (pp. 337–858)
- Each entry: `Name`, `Temperature`, `Mass`, `Luminosity`, `InnerLimit`, `HabitableMin`, `HabitableMax`, `SnowLine`, `OuterLimit`

**Tables still needed (not yet in source doc review):**
- Rocky planet type by orbital zone (2d6 tables for inner/habitable/outer)
- Planet size/atmosphere/hydrographics tables
- UWP (Universal World Profile) generation — deferred to separate package per plan
- Moon generation tables
- Eccentricity table (2d10 → 0.002–0.250)
- Tech level table (1d6 → 10–12)

#### Generation Steps — ❌ NOT IMPLEMENTED

All 18 step functions declared with `TODO` bodies in `types.go`.

**Missing files (per plan):**
- `generator.go` — main orchestration logic (currently stubs in `types.go`)
- `steps.go` — step implementations
- `roller.go` — dice adapter (would wrap `dice.Roll`/`dice.RollSafe`)
- `errors.go` — custom error types

---

## Dependencies

| Package | Usage |
|---------|-------|
| `internal/domain/engine/dice` | All random rolls in systemgen + tables |
| Standard library only | `fmt`, `math/rand`, `strconv`, `strings`, `slices`, `sort`, `sync`, `time`, `regexp`, `errors` |

No external dependencies.

---

## Build & Test Status

```
go build ./...     ✅ clean
go vet ./...       ✅ clean
go test ./...      ✅ all pass
  dice             ✅ 0.002s
  tables           ✅ 0.005s (29 tests)
  systemgen        [no test files]
```

---

## Next Steps (Priority Order)

1. **Complete star zone tables** — Ia, Ib, II, III, IV, VI, D (~539 entries)
2. **Add planet detail tables** — rocky planet types, size/atmosphere/hydrographics, eccentricity
3. **Implement generation steps 1–11** — deterministic, table-lookup driven
4. **Implement steps 12–18** — placement logic + complex migrations
5. **Extract `steps.go` and `generator.go`** from current `types.go` stubs
6. **Add `errors.go`** with domain-specific error types
7. **Add `roller.go`** dice adapter for systemgen
8. **Implement `RollCascade`** in tables package
9. **Write tests** — table verification, step unit tests, integration tests with known seeds
