package dice

import "math/rand"

// DieSpec описывает один кубик (неизменяемый шаблон).
type Die struct {
	Faces    int
	Codes    map[int]string    // текстовые обозначения для UI/логов (например, 20→"crit")
	Metadata map[string]string // цвет, имя, теги
}

// PoolSpec описывает пул кубиков и модификаторы уровня пула (неизменяемый).
type Dicepool struct {
	Type      string
	Dice      []Die
	Modifiers []Mod // в порядке применения
	Metadata  map[string]string
}

type Roller struct {
	rng *rand.Rand
}

type Result struct {
	Rolled         Dicepool
	Raw            []int
	Mods           []Mod
	Final          []int
	FinalAsStrings []string
}
