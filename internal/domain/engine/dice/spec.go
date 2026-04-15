// Package dice provides a flexible dice rolling system with expression parsing,
// modifier chaining, and support for common RPG mechanics (D66, Flux, etc.).
package dice

import "sync"

var defaultRoller roller
var defaultManager *Manager

func init() {
	defaultRoller = newRoller("")
	defaultManager = newManager(defaultRoller)
}

// Roller is the interface that wraps the Roll method.
// It defines how a single die is rolled.
type Roller interface {
	roll(die) int
}

// Manager coordinates dice rolling, expression caching, and result interpretation.
type Manager struct {
	roller    roller
	rollState *rollState
	mu        sync.Mutex
}

// New creates a new Manager with the given Roller.
// If roller is nil, the default (math/rand with random seed) is used.
func New(seed string) (*Manager, error) {
	return newManager(newRoller(seed)), nil
}

// Result contains the dice that were rolled and their raw values (order preserved).
type Result struct {
	dice []die
	raw  []int
}

// Dice returns a copy of the dice that were rolled.
func (r Result) Dice() []die {
	return append([]die(nil), r.dice...)
}

// Raw returns a copy of the raw roll results.
func (r Result) Raw() []int {
	return append([]int(nil), r.raw...)
}
