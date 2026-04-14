package dice

import "math/rand"

var defaultRoller Roller
var defaultManager *Manager

func init() {
	defaultRoller = newRoller("")
	defaultManager = newManager(defaultRoller)
}

// Die is a elementary component of Dicepool
type Die struct {
	Faces    int
	Codes    map[int]string    // текстовые обозначения для UI/логов (например, 20→"crit")
	Metadata map[string]string // цвет, имя, теги
}

// Dicepool describes composition of dice and modifications of the result.
type Dicepool struct {
	Type      string
	Dice      []Die
	Modifiers []Mod // в порядке применения
	Metadata  map[string]string
	roller    *Roller
}

// Roller

type Roller interface {
	Roll(Die) int
}

type randRoller struct {
	rng *rand.Rand
}

func (r *randRoller) Roll(d Die) int {
	return r.rng.Intn(d.Faces) + 1
}
