package dice

import (
	"math/rand"
	"time"
)

// roller is the internal interface (matches public Roller but unexported method).
type roller interface {
	roll(die) int
}

// randRoller is the default implementation using math/rand.
type randRoller struct {
	rng *rand.Rand
}

func (r *randRoller) roll(d die) int {
	return r.rng.Intn(d.faces) + 1
}

// newRoller creates a roller seeded either by the given string or a random seed.
func newRoller(seed string) roller {
	seedInt := randomSeed()
	if seed != "" {
		seedInt = stringToInt64(seed)
	}
	r := randRoller{}
	r.rng = rand.New(rand.NewSource(seedInt))
	return &r
}

// stringToInt64 converts a string into an int64 seed for deterministic randomness.
func stringToInt64(seed string) int64 {
	if seed == "" {
		return 0
	}
	var h1, h2 int64 = 0, 0
	for i := 0; i < len(seed); i++ {
		c := int64(seed[i])
		h1 = h1*131 + c + int64(i)*31
		h2 = h2*257 + c - int64(i)*17
		h1 ^= h2
		h2 ^= h1
	}
	return h1 ^ (h2 << 32) ^ int64(len(seed))
}

// randomSeed returns a seed derived from the current nanosecond timestamp.
func randomSeed() int64 {
	return time.Now().UnixNano()
}