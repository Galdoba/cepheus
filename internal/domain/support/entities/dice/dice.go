package dice

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Code_1D  = "1D"
	Code_2D  = "2D"
	Code_3D  = "3D"
	Code_4D  = "4D"
	Code_5D  = "5D"
	Code_6D  = "6D"
	Code_7D  = "7D"
	Code_8D  = "8D"
	Code_9D  = "9D"
	Code_10D = "10D"
	Code_1DD = "1DD"
	Code_D2  = "D2"
	Code_D3  = "D3"
	Code_D9  = "D9"
	Code_D10 = "D10"
	Code_D66 = "D66"
)

type Dicepool struct {
	seedValue  int64
	lastRolled []int
	lastSum    int
	roller     *rand.Rand
}

type source uint64

func (s source) Uint64() uint64 {
	return uint64(s)
}

func New(seed string) *Dicepool {
	d := newDice()
	if seed != "" {
		d.roller = rand.New(rand.NewSource(seedFromString(seed)))
	}
	return d
}

func newDice() *Dicepool {
	seedVal := time.Now().UnixMilli()
	r := rand.New(rand.NewSource(seedVal))
	d := Dicepool{
		seedValue:  seedVal,
		lastRolled: []int{},
		roller:     r,
	}
	return &d
}

func (dp *Dicepool) Roll(code string, mods ...int) int {
	num, faces, primeMod := parseDiceString(code)
	mods = append(mods, primeMod)
	dp.roll(num, faces, mods...)
	return dp.lastSum
}

func (dp *Dicepool) D66(mods ...int) string {
	for len(mods) < 2 {
		mods = append(mods, 0)
	}
	r1 := minmax(dp.Roll("1D", mods[0]), 1, 9)
	r2 := minmax(dp.Roll("1D", mods[1]), 1, 9)
	return fmt.Sprintf("%v%v", r1, r2)
}

func D66(mods ...int) string {
	for len(mods) < 2 {
		mods = append(mods, 0)
	}
	dp := newDice()
	r1 := minmax(dp.Roll("1D", mods[0]), 1, 9)
	r2 := minmax(dp.Roll("1D", mods[1]), 1, 9)
	return fmt.Sprintf("%v%v", r1, r2)
}

func minmax(i, minimum, maximum int) int {
	i = max(i, minimum)
	i = min(i, maximum)
	return i
}

func seedFromString(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())

}

func (dp *Dicepool) roll(count, faces int, mods ...int) {
	dp.lastRolled = []int{}
	dp.lastSum = 0
	for range count {
		r := dp.roller.Intn(faces) + 1
		dp.lastRolled = append(dp.lastRolled, r)
		dp.lastSum += r
	}
	dp.lastSum += sum(mods...)
}

func sum(mods ...int) int {
	s := 0
	for _, m := range mods {
		s += m
	}
	return s
}

func parseDiceString(s string) (diceCount, faceCount, modifier int) {
	s = strings.ToUpper(s)
	re := regexp.MustCompile(`^(\d+)?D(\d+)?([+-]\d+)?$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		panic(fmt.Sprintf("invalid dice code provided: '%s'", s))
	}

	diceCount = 1
	if matches[1] != "" {
		diceCount, _ = strconv.Atoi(matches[1])
	}

	faceCount = 6
	if matches[2] != "" {
		faceCount, _ = strconv.Atoi(matches[2])
	}

	modifier = 0
	if matches[3] != "" {
		modifier, _ = strconv.Atoi(matches[3])
	}

	return diceCount, faceCount, modifier
}

func Roll(code string, mods ...int) int {
	d := newDice()
	num, faces, primeMod := parseDiceString(code)
	mods = append(mods, primeMod)
	d.roll(num, faces, mods...)
	return d.lastSum
}

func CharacteristicDM(i int) int {
	i = max(i, 0)
	switch i {
	case 0:
		return -3
	case 1, 2:
		return -2
	default:
		return (i / 3) - 2
	}
}
