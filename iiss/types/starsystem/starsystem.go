package starsystem

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/gametable"
)

const (
	method_Extended = iota
	method_Continuation
)

type StarSystemGenerator struct {
	method                int
	injectedStellar       string
	injectedUWP           string
	injectedGGquantity    int
	injectedBeltsQuantity int
	dp                    *dice.Dicepool
	//Star            map[int]*star.Star `json:"stars,omitempty"`
	//Orbits map[float64]orbit.Orbit
}

func NewGenerator(options ...SSG_Option) *StarSystemGenerator {
	ssg := StarSystemGenerator{
		injectedGGquantity:    -1,
		injectedBeltsQuantity: -1,
	}

	for _, modify := range options {
		modify(&ssg)
	}
	if ssg.dp == nil {
		ssg.dp = dice.NewDicepool()
	}
	return &ssg
}

type SSG_Option func(*StarSystemGenerator)

func WithDice(dp *dice.Dicepool) SSG_Option {
	return func(ssg *StarSystemGenerator) {
		ssg.dp = dp
	}
}

func WithStellar(stellar string) SSG_Option {
	return func(ssg *StarSystemGenerator) {
		ssg.injectedStellar = stellar
		ssg.method = method_Continuation
	}
}

func WithUWP(profile string) SSG_Option {
	return func(ssg *StarSystemGenerator) {
		ssg.injectedUWP = profile
	}
}

func NonPrimaryStarDetection(dp *dice.Dicepool, column string) (string, error) {
	err := fmt.Errorf("table not created")
	otherTable, err := gametable.NewTable("other", "2d6",
		gametable.NewRollResult("2-", "NS", nil),
		gametable.NewRollResult("3..7", "D", nil),
		gametable.NewRollResult("8+", "BD", nil),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create otherTable: %v", err)
	}
	postStellarTable, err := gametable.NewTable("postStellar", "2d6",
		gametable.NewRollResult("3-", "other", otherTable),
		gametable.NewRollResult("4..8", "random", nil),
		gametable.NewRollResult("9..10", "lesser", nil),
		gametable.NewRollResult("11+", "twin", nil),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create postStellarTable: %v", err)
	}
	companionTable, err := gametable.NewTable("companion", "2d6",
		gametable.NewRollResult("3-", "other", otherTable),
		gametable.NewRollResult("4..5", "random", nil),
		gametable.NewRollResult("6..7", "lesser", nil),
		gametable.NewRollResult("8..9", "sibling", nil),
		gametable.NewRollResult("10+", "twin", nil),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create companionTable: %v", err)
	}
	secondaryTable, err := gametable.NewTable("secondary", "2d6",
		gametable.NewRollResult("3-", "other", otherTable),
		gametable.NewRollResult("4..6", "random", nil),
		gametable.NewRollResult("7..8", "lesser", nil),
		gametable.NewRollResult("9..10", "sibling", nil),
		gametable.NewRollResult("11+", "twin", nil),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create secondaryTable: %v", err)
	}
	r := ""
	switch column {
	case "secondary":
		r, err = secondaryTable.Roll(dp)
		if err != nil {
			return "", fmt.Errorf("secondary table roll: %v", err)
		}
	case "companion":
		r, err = companionTable.Roll(dp)
		if err != nil {
			return "", fmt.Errorf("companion table roll: %v", err)
		}
	case "postStellar":
		r, err = postStellarTable.Roll(dp)
		if err != nil {
			return "", fmt.Errorf("postStellar table roll: %v", err)
		}
	}
	return r, nil
}

func positions(dp *dice.Dicepool, dm int) []string {
	pos := []string{"Aa"}
	for _, newPos := range []string{"Ba", "Ca", "Da"} {
		if addStar(dp, dm) {
			pos = append(pos, newPos)
		}
	}
	for _, starPos := range pos {
		if addStar(dp, dm) {
			pos = append(pos, starPos+"b")
		}
	}
	for i := range pos {
		pos[i] = strings.ReplaceAll(pos[i], "ab", "b")
	}
	slices.Sort(pos)
	return pos
}

func addStar(dp *dice.Dicepool, dm int) bool {
	r := dp.Sum("2d6") + dm
	return r >= 10
}

func bounded(i, min, max int) int {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}
