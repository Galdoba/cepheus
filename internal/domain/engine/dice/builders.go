package dice

import (
	"math/rand"
	"time"
)

// newRoller initializes a Roller with a specific seed.
// It converts the seed string to an int64 or uses a random seed if empty.
func newRoller(seed string) *Roller {
	seedInt := randomSeed()
	if seed != "" {
		seedInt = stringToInt64(seed)
	}
	r := Roller{}
	r.rng = rand.New(rand.NewSource(int64(seedInt)))
	return &r
}

// stringToInt64 converts a string seed to an int64.
// This ensures consistent random sequences from the same seed string.
func stringToInt64(seed string) int64 {
	if seed == "" {
		return 0
	}

	var h1, h2 int64 = 0, 0

	for i := 0; i < len(seed); i++ {
		c := int64(seed[i])

		// First polynomial with a large prime factor
		h1 = h1*131 + c + int64(i)*31

		// Second polynomial with a different factor for better mixing
		h2 = h2*257 + c - int64(i)*17

		// Mix the hashes together
		h1 ^= h2
		h2 ^= h1
	}

	// Combine both hashes
	return h1 ^ (h2 << 32) ^ int64(len(seed))
}

// randomSeed generates a time-based random seed using the current time in nanoseconds.
func randomSeed() int64 {
	return time.Now().UnixNano()
}

//Die

func NewDice(faces int) Die {
	return Die{Faces: faces}
}

func (d Die) WithCodes(codes map[int]string) Die {
	d.Codes = codes
	return d
}

func (d Die) WithMeta(meta map[string]string) Die {
	d.Metadata = meta
	return d
}

// Dicepool

func NewDicepool(dice ...Die) Dicepool {
	dp := Dicepool{
		Dice:      dice,
		Modifiers: []Mod{None{}},
		Metadata:  map[string]string{},
	}
	return dp
}

func (dp Dicepool) WithMods(mods ...Mod) Dicepool {
	dp.Modifiers = mods
	return dp
}

func (dp Dicepool) WithMeta(meta map[string]string) Dicepool {
	dp.Metadata = meta
	return dp
}
