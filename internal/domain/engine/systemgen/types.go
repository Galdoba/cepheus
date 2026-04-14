// Package systemgen implements the Cepheus Engine System Generation algorithm.
//
// The generation process creates procedurally generated star systems following
// the 18-step ruleset defined in System_Generation_Extended.md.
package systemgen

import (
	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

// ---------------------------------------------------------------------------
// Subsector and hex configuration
// ---------------------------------------------------------------------------

// SubsectorType defines the density classification for a subsector,
// which determines the probability of an object existing in a given hex.
type SubsectorType int

const (
	SubEmpty SubsectorType = iota
	SubScattered
	SubDispersed
	SubAverage
	SubCrowded
	SubDense
)

// RealismMode controls which stellar class probability table is used.
type RealismMode int

const (
	Realistic RealismMode = iota
	SemiRealistic
	Fantastic
)

// PlacementMethod controls the algorithm used for orbital planet placement.
type PlacementMethod int

const (
	PlacementDoubling PlacementMethod = iota
	PlacementVaried
)

// ---------------------------------------------------------------------------
// Object and body type enumerations
// ---------------------------------------------------------------------------

// ObjectType is the classification of the primary object in a hex.
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

// StellarClass represents the spectral classification of a star (O through M).
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

// String returns the single-letter stellar class designation.
func (s StellarClass) String() string {
	return []string{"O", "B", "A", "F", "G", "K", "M"}[s]
}

// LuminosityClass represents the MKK luminosity classification of a star.
type LuminosityClass int

const (
	LumV  LuminosityClass = iota // Main Sequence
	LumIV                         // Subgiant
	LumIII                        // Giant
	LumII                         // Bright Giant
	LumIa                         // Luminous Supergiant
	LumIb                         // Supergiant
	LumVI                         // Subdwarf
	LumD                          // Dwarf (white dwarf)
)

// String returns the roman numeral luminosity class designation.
func (l LuminosityClass) String() string {
	return []string{"V", "IV", "III", "II", "Ia", "Ib", "VI", "D"}[l]
}

// BrownDwarfClass represents the spectral class of a brown dwarf.
type BrownDwarfClass int

const (
	ClassL BrownDwarfClass = iota
	ClassT
	ClassY
)

// String returns the single-letter brown dwarf class designation.
func (b BrownDwarfClass) String() string {
	return []string{"L", "T", "Y"}[b]
}

// ---------------------------------------------------------------------------
// Zone and range types
// ---------------------------------------------------------------------------

// ZoneRange represents an inner/outer boundary pair in astronomical units.
type ZoneRange struct {
	Inner float64
	Outer float64
}

// StarZoneEntry contains the complete physical and orbital data for a star type.
type StarZoneEntry struct {
	Name         string  // e.g. "O0 V", "G2 V", "L3"
	Temperature  int     // Kelvin
	Mass         float64 // Solar masses
	Luminosity   float64 // Solar units
	InnerLimit   float64 // AU — minimum viable orbit
	HabitableMin float64 // AU — inner edge of habitable zone
	HabitableMax float64 // AU — outer edge of habitable zone
	SnowLine     float64 // AU — 0 means N/A (no snow line)
	OuterLimit   float64 // AU — maximum viable orbit
}

// ---------------------------------------------------------------------------
// Core domain types
// ---------------------------------------------------------------------------

// HexCoords identifies a hex within a subsector.
type HexCoords struct {
	Row int
	Col int
}

// StarSystem is the complete generated model of a single hex's contents.
type StarSystem struct {
	HexCoords     HexCoords
	ObjectType    ObjectType
	PrimaryBody   CelestialBody
	Companions    []CompanionStar
	GasGiants     []GasGiant
	AsteroidBelts []AsteroidBelt
	RockyPlanets  []RockyPlanet
	ZimmPoints    []ZimmPoint
	Quirks        []SystemQuirk
}

// CelestialBody holds whichever type of primary object occupies the hex.
// Exactly one field will be non-nil.
type CelestialBody struct {
	Star       *StarData
	BrownDwarf *BrownDwarfData
	RoguePlanet *RoguePlanetData
}

// StarData contains the full classification and zone data for a star.
type StarData struct {
	StellarClass    StellarClass
	NumericClass    int            // 0–9
	LuminosityClass LuminosityClass
	Temperature     int            // Kelvin
	Mass            float64        // Solar masses
	Luminosity      float64        // Solar units
	InnerLimit      float64        // AU
	HabitableZone   ZoneRange      // AU
	SnowLine        float64        // AU (0 = N/A)
	OuterLimit      float64        // AU
}

// BrownDwarfData contains the full classification and zone data for a brown dwarf.
type BrownDwarfData struct {
	Class        BrownDwarfClass
	NumericClass int
	Temperature  int
	Mass         float64
	Luminosity   float64
	InnerLimit   float64
	HabitableZone ZoneRange
	SnowLine     float64
	OuterLimit   float64
}

// RoguePlanetData contains minimal classification for a rogue planet.
type RoguePlanetData struct {
	// Rogue planets are placed directly without the full star generation flow.
	// Additional fields can be added as the implementation progresses.
}

// CompanionStar describes a secondary (or tertiary, etc.) star in a multiple system.
type CompanionStar struct {
	StarData StarData
	Distance CompanionDistance
	AUs      float64
}

// CompanionDistance is the categorical separation between stars in a multiple system.
type CompanionDistance int

const (
	DistContact CompanionDistance = iota
	DistClose
	DistNear
	DistFar
	DistDistant
)

// GasGiant represents a single gas giant in the system.
type GasGiant struct {
	SizeType        GasGiantSize // Neptunian or Jovian
	SizeCode        int          // Derived from 2d6 (Neptunian) or 2d10 (Jovian)
	OrbitAU         float64
	Migrated        bool
	OriginalOrbitAU float64 // Pre-migration orbit, 0 if never migrated
}

// GasGiantSize distinguishes small (Neptunian) from large (Jovian) gas giants.
type GasGiantSize int

const (
	SizeNeptunian GasGiantSize = iota
	SizeJovian
)

// AsteroidBelt represents a single asteroid belt in the system.
type AsteroidBelt struct {
	OrbitAU float64
	Width   float64 // Optional detail field
}

// RockyPlanet represents a terrestrial-class planet (dwarf through superterran).
type RockyPlanet struct {
	PlanetType     PlanetType
	SizeCode       int
	AtmosphereCode int
	Hydrographics  int
	Population     int
	Government     int
	LawLevel       int
	TechLevel      int
	OrbitAU        float64
	Eccentricity   float64
	IsMainWorld    bool
	Moons          []Moon
}

// PlanetType is the size/mass category of a rocky planet.
type PlanetType int

const (
	TypeDwarf PlanetType = iota
	TypeMercurian
	TypeSubterran
	TypeTerran
	TypeSuperterran
)

// Moon is a natural satellite of a rocky planet.
type Moon struct {
	SizeCode       int
	AtmosphereCode int
	OrbitDistance  float64
}

// ZimmPoint is a faster-than-light jump point in the system.
type ZimmPoint struct {
	Type       ZimmPointType
	OrbitAU    float64
	TargetName string
}

// ZimmPointType distinguishes incoming from outgoing jump points.
type ZimmPointType int

const (
	ZPIncoming ZimmPointType = iota
	ZPOutgoing
)

// SystemQuirk records a special event or anomaly rolled during generation.
type SystemQuirk struct {
	RollResult  int
	Description string
}

// ---------------------------------------------------------------------------
// Generator public API
// ---------------------------------------------------------------------------

// Generator orchestrates system generation for a subsector of a given density.
type Generator struct {
	subsectorType SubsectorType
	seed          string
}

// GenerationResult is the output of a single hex generation run.
type GenerationResult struct {
	System   StarSystem
	Log      []string
	Warnings []string
}

// New creates a Generator configured with the given subsector density and RNG seed.
func New(subsectorType SubsectorType, seed string) *Generator {
	return &Generator{
		subsectorType: subsectorType,
		seed:          seed,
	}
}

// GenerateHex generates a complete star system for the hex at (row, col).
//
// The generation follows the 18-step process defined in System_Generation_Extended.md.
// Returns ErrNoObjectInHex if step 1 determines no object exists, or
// ErrBlackHole if the object is a black hole (no further generation).
func (g *Generator) GenerateHex(row, col int) (*GenerationResult, error) {
	// TODO: implement full generation
	return nil, nil
}

// GenerateSubsector generates all hexes in a subsector.
// A standard subsector is 8 columns × 5 rows = 40 hexes.
func (g *Generator) GenerateSubsector() ([]*GenerationResult, error) {
	// TODO: iterate over 8×5 grid, call GenerateHex for each
	return nil, nil
}

// ---------------------------------------------------------------------------
// Generation step function signatures
//
// Each function corresponds to a numbered step in the document.
// They are declared here as package-level functions for testability.
// ---------------------------------------------------------------------------

// step1 determines whether an object exists in the current hex.
// Returns true if an object is present.
func step1(dice dice.Roller, subsectorType SubsectorType) (bool, error) {
	// TODO: roll d100, compare against subsector threshold
	return false, nil
}

// step2 determines the type of object present.
func step2(dice dice.Roller) (ObjectType, error) {
	// TODO: roll d100 on object type table
	return ObjectStar, nil
}

// step3 determines the brown dwarf subtype (L, T, or Y).
func step3(dice dice.Roller) (BrownDwarfClass, error) {
	// TODO: roll d100 on brown dwarf class table
	return ClassL, nil
}

// step4 determines the stellar class using the chosen realism mode.
func step4(dice dice.Roller, realism RealismMode) (StellarClass, error) {
	// TODO: roll d100 on appropriate stellar class table
	return ClassM, nil
}

// step5 determines the numeric classification (0–9) of the star or brown dwarf.
func step5(dice dice.Roller) (int, error) {
	// TODO: roll 1d10 - 1
	return 0, nil
}

// step6 determines the luminosity class of a star.
func step6(dice dice.Roller) (LuminosityClass, error) {
	// TODO: roll d100 on luminosity class table, handle special d10 sub-roll
	return LumV, nil
}

// step7 determines whether the system is single, binary, trinary, etc.
// Returns the total number of stars in the system.
func step7(dice dice.Roller, primary StellarClass) (int, error) {
	// TODO: roll d100 on appropriate multiple-star table
	return 1, nil
}

// step8 determines the orbital distances of companion stars.
func step8(dice dice.Roller, numCompanions int) ([]CompanionStar, error) {
	// TODO: for each companion, roll d100 for distance category, then d100 for AU
	return nil, nil
}

// step9 determines the number of gas giants in the system.
func step9(dice dice.Roller, body CelestialBody) (int, error) {
	// TODO: roll 1d6 with modifier based on primary body type
	return 0, nil
}

// step10 determines the number of asteroid belts.
func step10(dice dice.Roller) (int, error) {
	// TODO: roll 1d6 - 3, minimum 0
	return 0, nil
}

// step11 looks up the zone data (habitable zone, snow line, etc.) for a celestial body.
func step11(body CelestialBody) (StarZoneEntry, error) {
	// TODO: lookup in static tables
	return StarZoneEntry{}, nil
}

// step12 handles pre-existing main world placement.
// If mainWorld is nil this step is a no-op.
func step12(system *StarSystem, mainWorld *RockyPlanet) error {
	// TODO: place mainWorld into system.RockyPlanets at a valid orbit
	return nil
}

// step13 places gas giants into orbits and handles migration logic.
func step13(dice dice.Roller, system *StarSystem, zones StarZoneEntry) error {
	// TODO: place gas giants near or beyond snow line, check for migration
	return nil
}

// step14 places rocky planets and asteroid belts using the chosen placement method.
func step14(dice dice.Roller, system *StarSystem, zones StarZoneEntry, method PlacementMethod) error {
	// TODO: double or varied-distance placement
	return nil
}

// step15 resolves details of all rocky planets (type, size, atmosphere, etc.).
func step15(dice dice.Roller, system *StarSystem, zones StarZoneEntry) error {
	// TODO: iterate planets, determine type by orbital zone, then fill UWP
	return nil
}

// step15A details a planet between the inner limit and the habitable zone.
func step15A(dice dice.Roller, planet *RockyPlanet) error {
	// TODO: roll for type, then size/atmosphere/etc.
	return nil
}

// step15B details a planet within the habitable zone.
func step15B(dice dice.Roller, planet *RockyPlanet, zones StarZoneEntry) error {
	// TODO: roll for type (terran/superterran), then full UWP
	return nil
}

// step15C details a planet beyond the habitable zone.
func step15C(dice dice.Roller, planet *RockyPlanet) error {
	// TODO: roll for type, then cold-world rules
	return nil
}

// step16 generates moons for a rocky planet.
func step16(dice dice.Roller, planet *RockyPlanet) error {
	// TODO: roll for moon count, then generate each moon
	return nil
}

// step17 places Zimm Points based on adjacent system names.
func step17(dice dice.Roller, system *StarSystem, adjacentSystems []string) error {
	// TODO: place incoming and outgoing Zimm Points at valid orbits
	return nil
}

// step18 rolls for system quirks.
func step18(dice dice.Roller) ([]SystemQuirk, error) {
	// TODO: roll d66 on quirk table
	return nil, nil
}
