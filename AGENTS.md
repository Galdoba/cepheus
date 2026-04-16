# Cepheus — Agent Guide

## What This Project Is

**Cepheus** is a Go library implementing the **Cepheus Engine** RPG rules — a Traveller-derived sci-fi tabletop system. It provides four packages:

| Package | Status | Purpose |
|---------|--------|---------|
| `dice` | ✅ Complete | Dice expression parsing, modifier chaining, D66/Flux/Variance |
| `tables` | ✅ Complete | RPG lookup tables with index parsing, validation, cascading rolls |
| `ehex` | ✅ Complete | Traveller extended hex encoding (0-9, A-Z excl I/O) for UWP codes |
| `systemgen` | 🚧 Partial | Procedural star-system generator (types + tables defined, generation logic not implemented) |

**Key facts:** Go 1.25.1 · zero external dependencies · 35 tests · all build/vet/test clean · `assets/` contains 127 JSON data tables.

---

## Documentation Index

Read the relevant file before working on a package.

| File | When to Read |
|------|-------------|
| [`doc/project_status.md`](doc/project_status.md) | **Start here** for current status, what's done, what's next |
| [`doc/dice.md`](doc/dice.md) | Working on dice expressions, modifiers, or the Manager API |
| [`doc/tables.md`](doc/tables.md) | Working on GameTable, Collection, RollCascade, or index parsing |
| [`doc/ehex.md`](doc/ehex.md) | Working on UWP encoding or special sentinel codes |
| [`doc/system_generator_plan.md`](doc/system_generator_plan.md) | Implementing the 18-step system generation algorithm |
| [`doc/System_Generation_Extended.md`](doc/System_Generation_Extended.md) | Reference: the source ruleset for system generation |

---

## Essential Commands

```bash
go build ./...     # always passes
go vet ./...       # always passes
go test ./...      # 35 tests, all pass
```

Run these after any code change.

---

## Architecture at a Glance

```
internal/domain/engine/
├── dice/        ← Foundation: parsing, rolling, caching, special mechanics
├── tables/      ← Builds on dice: GameTable + Collection via TableRoller interface
├── ehex/        ← Independent: UWP character↔value encoding
└── systemgen/   ← Builds on dice + tables: 18-step procedural generation (WIP)

assets/          ← 127 JSON files: star zones, planets, moons, formulas, quirks
doc/             ← This project's documentation
```

**Dependency flow:** `dice` → `tables` → `systemgen`. `ehex` has no dependencies.

---

## Working Principles

- **No external dependencies** — use only the standard library
- **All random rolls go through `dice`** — no direct `math/rand` outside `dice/roller.go`
- **Table data can live in `assets/*.json`** — not all data needs to be hardcoded in Go
- **`systemgen` types are defined** in `types.go` — the 18 step functions have `TODO` bodies
- **Tests exist for `dice`, `tables`, `ehex`** — write tests for new `systemgen` code
- **Follow the plan** in `doc/system_generator_plan.md` for generation implementation
