package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
	"github.com/Galdoba/cepheus/pkg/dice"
)

// Builder orchestrates the creation of a star system through a multi-step process.
// It maintains state for each step of star system generation and uses random
// number generation to determine star types, orbits, and system worlds.
type Builder struct {
	rng      *dice.Roller
	options  map[string]bool
	imported t5ss.WorldData
	step1    *primaryStarDeterminator
	step2    *secondaryStarsDeterminator
	step3    *systemWorldsDeterminator
}

// primaryStarDeterminator handles Step 1 of star system creation:
// determining the primary star's type, subtype, and class (TSC).
type primaryStarDeterminator struct {
	tablesStarType rtg.RandomTableGenerator
	activeMods     map[string]bool
	completed      bool
}

// secondaryStarsDeterminator handles Step 2 of star system creation:
// determining secondary stars, their orbits, and physical properties.
type secondaryStarsDeterminator struct {
	tables     rtg.RandomTableGenerator
	starSchema []stellar.StarDesignation
	completed  bool
}

// systemWorldsDeterminator handles Step 3 of star system creation:
// determining gas giants, asteroid belts, and terrestrial planets.
type systemWorldsDeterminator struct {
	tablesPlanets rtg.RandomTableGenerator
	activeMods    map[string]bool
	completed     bool
	ggPresent     bool
	ggNum         int
	bltPresent    bool
	bltNum        int
	tpPresent     bool
	tpNum         int
}

// BuildOption defines a functional option for configuring the Builder.
// TODO: Implement optional configuration options such as:
// - Force specific star types
// - Set system age limits
// - Import/export system data
type BuildOption func(*Builder)

// NewBuilder creates a new Builder instance with the given seed for random number generation.
// It initializes the random table generators for star type and planet determination.
// Returns an error if the random number generators cannot be initialized.
func NewBuilder(seed string, options ...BuildOption) (*Builder, error) {
	builder := Builder{}
	builder.rng = dice.New(seed)
	builder.options = make(map[string]bool)

	primaryStarDeterminatorStep := primaryStarDeterminator{}
	starTypeTable, err := rtg.NewStarTypeDeterminationGenerator(builder.rng)
	if err != nil {
		return nil, fmt.Errorf("failed to create RNG: Star Type Determination: %v", err)
	}
	primaryStarDeterminatorStep.tablesStarType = starTypeTable

	secondaryStarsDeterminatorStep := secondaryStarsDeterminator{}
	secondaryStarsDeterminatorStep.tables = starTypeTable

	systemWorldsDeterminatorStep := systemWorldsDeterminator{}
	planetsTable, err := rtg.NewSystemPlanetsDeterminationGenerator(builder.rng)
	if err != nil {

		return nil, fmt.Errorf("failed to create RNG: System Worlds Determination: %v", err)
	}
	systemWorldsDeterminatorStep.tablesPlanets = planetsTable

	builder.step1 = &primaryStarDeterminatorStep
	builder.step2 = &secondaryStarsDeterminatorStep
	builder.step3 = &systemWorldsDeterminatorStep

	// TODO: Apply configuration options if implemented
	// for _, option := range options {
	// 	option(&builder)
	// }

	return &builder, nil
}

// newStarSystemPrecursor creates a new starSystemPrecursor with default/empty values.
// The GG, Belts, and Planets fields are initialized to -1000 as sentinel values
// indicating they have not yet been determined.
// TODO: Replace sentinel values (-1000) with explicit boolean flags or nullable types
// for better type safety and clarity.
func newStarSystemPrecursor() *starSystemPrecursor {
	systemPrecursor := starSystemPrecursor{}
	systemPrecursor.Stars = make(map[orbit.Orbit]*starPrecursor)
	systemPrecursor.GG = -1000
	systemPrecursor.Belts = -1000
	systemPrecursor.Planets = -1000
	return &systemPrecursor
}

// Build executes the three-step star system generation process:
//   - Step 1: Determine primary star type, subtype, and class (TSC)
//   - Step 2: Determine secondary stars and their orbital characteristics
//   - Step 3: Determine system worlds (gas giants, belts, planets)
//
// The directives parameter is reserved for future use to allow selective
// execution of specific generation steps.
func (builder *Builder) Build(directives ...string) (*starSystemPrecursor, error) {

	systemPrecursor := newStarSystemPrecursor()
	if err := builder.runStep01(systemPrecursor); err != nil {
		return nil, fmt.Errorf("step 1 failed: %v", err)
	}
	if err := builder.runStep02(systemPrecursor); err != nil {
		return nil, fmt.Errorf("step 2 failed: %v", err)
	}
	if err := builder.runStep03(systemPrecursor); err != nil {
		return nil, fmt.Errorf("step 3 failed: %v", err)
	}
	if err := builder.runStep04(systemPrecursor); err != nil {
		return nil, fmt.Errorf("step 4 failed: %v", err)
	}

	return systemPrecursor, nil
}
