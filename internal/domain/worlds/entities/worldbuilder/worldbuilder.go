package worldbuilder

import "github.com/Galdoba/cepheus/pkg/dice"

const (
	randomSeed = ""
	genericID  = "generic"
)

// WorldBuilder represents a configurable world generation engine.
// It maintains an identifier and a random number generator (dice pool)
// for procedural content generation.
type WorldBuilder struct {
	id   string       // Unique identifier for this builder instance
	seed string       // RNG seed
	dice *dice.Roller // Random number generator for procedural generation
}

// WorldBuilderOption defines a function type for configuring WorldBuilder
// instances through the functional options pattern. Each option modifies
// the WorldBuilder during initialization.
type WorldBuilderOption func(*WorldBuilder)

// WithSeed returns a WorldBuilderOption that injects a custom dice pool seed
// into the WorldBuilder. This allows for controlled randomness during
// world generation, useful for testing or seeded generation.
func WithSeed(seed string) WorldBuilderOption {
	return func(wb *WorldBuilder) {
		wb.seed = seed
	}
}

// WithID returns a WorldBuilderOption that sets a custom identifier
// for the WorldBuilder instance. Useful for tracking or differentiating
// multiple world generators.
func WithID(id string) WorldBuilderOption {
	return func(wb *WorldBuilder) {
		wb.id = id
	}
}

// New initializes a new WorldBuilder instance with default values,
// then applies the provided options in order. Defaults include:
// - id: "generic"
// - dice: a new dice.Dicepool instance with time based seed
func New(options ...WorldBuilderOption) *WorldBuilder {
	wb := WorldBuilder{}
	wb.id = genericID
	wb.seed = randomSeed
	for _, modify := range options {
		modify(&wb)
	}
	wb.dice = dice.New(wb.seed)
	return &wb
}
