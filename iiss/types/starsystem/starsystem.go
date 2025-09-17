package starsystem

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Galdoba/cepheus/iiss/types/orbit"
	"github.com/Galdoba/cepheus/iiss/types/star"
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

func (ssg *StarSystemGenerator) GenerateSystem() *StarSystem {
	ss := &StarSystem{}
	switch ssg.method {
	case method_Extended:
		ss = ssg.extendedStarGeneration()
	case method_Continuation:
		ss = ssg.continuationStarGeneration()
	}
	ssg.GenerateStarOrbits(ss)
	return ss
}

func (ssg *StarSystemGenerator) GenerateStarOrbits(ss *StarSystem) {
	ss.Orbits = make(map[float64]*orbit.Orbit)
	primary := ss.Stars["Aa"]
	codes := []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Cb", "Da", "Db"}
	for _, code := range codes {
		if ss.Stars[code] == nil {
			continue
		}
		parentCode := orbit.StarParent(code)
		parentMass := 0.0
		starMass := 0.0
		if parentCode != "" {
			parentMass = ss.Stars[parentCode].Mass
			starMass = ss.Stars[code].Mass
		}
		orb := orbit.NewStar(ssg.dp, code, orbit.StarMass(1, starMass), orbit.StarMass(2, parentMass), orbit.SystemAge(primary.Age), orbit.IsStar(true))

		ss.Orbits[orb.FromParent] = orb
	}
}

func (ssg *StarSystemGenerator) extendedStarGeneration() *StarSystem {
	ss := NewStarSystem()
	primary, _ := star.Generate(ssg.dp)
	codes := positions(ssg.dp, secondaryPlacementDM(primary))
	for _, code := range codes {
		switch code {
		case "Ab", "Bb", "Cb", "Db":
			primCode := strings.ReplaceAll(code, "b", "a")
			ss.Stars[code] = setStar(ssg.dp, code, *ss.Stars[primCode])
		default:
			ss.Stars[code] = setStar(ssg.dp, code, primary)
		}
	}
	ss.sortStarsByMass()
	return &ss
}

func (ssg *StarSystemGenerator) continuationStarGeneration() *StarSystem {
	starData := star.FromStellar(ssg.injectedStellar)
	ss := NewStarSystem()
	switch len(starData) {
	case 0:
		return ssg.extendedStarGeneration()
	default:
		ss.Stars = make(map[string]*star.Star)
		_, _, _, stellarData := starData[0].StarData()
		primary, _ := star.Generate(ssg.dp, star.KnownStellar(stellarData))
		codes := []string{}
		for len(codes) != len(starData) {
			codes = positions(ssg.dp, secondaryPlacementDM(primary))
		}
		for i, code := range codes {
			_, _, _, stellarData := starData[i].StarData()
			newStar, _ := star.Generate(ssg.dp, star.KnownStellar(stellarData))
			ss.Stars[code] = &newStar
		}
	}
	return &ss
}

func (ss *StarSystem) sortStarsByMass() bool {
	done := false
	hadSwitch := false
	if ss.Stars == nil {
		return true
	}
mainLoop:
	for !done {
		for _, code := range []string{"A", "B", "C", "D"} {
			if ss.Stars[code+"a"] == nil || ss.Stars[code+"b"] == nil {
				continue
			}
			if ss.Stars[code+"a"].Mass < ss.Stars[code+"b"].Mass {
				ss.Stars[code+"a"], ss.Stars[code+"b"] = ss.Stars[code+"b"], ss.Stars[code+"a"]
				hadSwitch = true
				continue mainLoop
			}
		}
		for _, code := range []string{"B", "C", "D"} {
			if ss.Stars[code+"a"] == nil {
				continue
			}
			if ss.Stars["Aa"].Mass < ss.Stars[code+"a"].Mass {
				ss.Stars["Aa"], ss.Stars[code+"a"] = ss.Stars[code+"a"], ss.Stars["Aa"]
				hadSwitch = true
				continue mainLoop
			}
		}
		done = true
	}

	return hadSwitch
}

func secondaryPlacementDM(primary star.Star) int {
	dm := 0
	switch primary.Class {
	case "Ia", "Ib", "II", "III", "IV":
		dm++
	case "V", "VI":
		switch primary.Type {
		case "O", "B", "A", "F":
			dm++
		case "M":
			dm--
		}
	case "D", "BD", "PSR", "BH", "NS":
		dm--
	}
	return dm
}

func setStar(dp *dice.Dicepool, code string, primary star.Star) *star.Star {
	column := ""
	switch code {
	case "Aa":
		return &primary
	case "Ba", "Ca", "Da":
		column = "secondary"
	case "Ab", "Bb", "Cb", "Db":
		column = "companion"
	}
	switch primary.Class {
	case "D", "BD", "NS", "PSR", "BH":
		column = "postStellar"
	}
	starRollType, err := NonPrimaryStarDetection(dp, column)
	if err != nil {
		panic(err)
	}
	if primary.ProtoStar || primary.Class == "BD" {
		starRollType = "twin"
	}
	newStar := star.Star{}
	switch starRollType {
	case "random":
		newStar = primary.Random(dp)
	case "lesser":
		newStar = primary.Lesser(dp)
	case "sibling":
		newStar = primary.Sibling(dp)
	case "twin":
		newStar = primary.Twin(dp)
	case "D", "BD", "NS":
		newStar, _ = star.Generate(dp, star.KnownStellar(starRollType))
	default:
		panic(starRollType)
	}

	return &newStar
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

type StarSystem struct {
	Primary     *star.Star
	Stars       map[string]*star.Star
	Orbits      map[float64]*orbit.Orbit
	presenceGG  int
	presenceBT  int
	presenceTP  int
	totalWorlds int
	MAO         float64
}

func NewStarSystem() StarSystem {
	ss := StarSystem{
		Primary:    &star.Star{},
		Stars:      make(map[string]*star.Star),
		Orbits:     make(map[float64]*orbit.Orbit),
		presenceGG: -1,
		presenceBT: -1,
		presenceTP: -1,
	}
	return ss
}

func (ssg *StarSystemGenerator) CreateSystemWorldsAndOrbits(ss *StarSystem) error {
	for _, setParameter := range []func() error{
		func() error { return ssg.setGGquantity(ss) },
		func() error { return ssg.setBTquantity(ss) },
		func() error { return ssg.setTPquantity(ss) },
	} {
		if err := setParameter(); err != nil {
			return fmt.Errorf("failed to Create System Worlds and Orbits: %v", err)
		}
	}
	ss.totalWorlds = ss.presenceGG + ss.presenceBT + ss.presenceTP

	return nil
}

func (ssg *StarSystemGenerator) setGGquantity(ss *StarSystem) error {
	if ss == nil {
		return fmt.Errorf("no Star System provided")
	}
	ss.presenceGG = ssg.injectedGGquantity
	targetNumber := 10
	if ss.Primary.Class == "BD" {
		targetNumber = 8
	}
	if ss.presenceGG < 0 {
		switch ssg.dp.Sum("2d6") >= targetNumber {
		case true:
			switch bounded(ssg.dp.Sum("2d6")+ggQuantityDM(ss), 4, 13) {
			case 4:
				ss.presenceGG = 1
			case 5, 6:
				ss.presenceGG = 2
			case 7, 8:
				ss.presenceGG = 3
			case 9, 10, 11:
				ss.presenceGG = 4
			case 12:
				ss.presenceGG = 5
			case 13:
				ss.presenceGG = 6
			}
		case false:
			ss.presenceGG = 0
		}
	}
	return nil
}

func (ssg *StarSystemGenerator) setBTquantity(ss *StarSystem) error {
	if ss == nil {
		return fmt.Errorf("no Star System provided")
	}
	ss.presenceBT = ssg.injectedBeltsQuantity
	if ss.presenceBT < 0 {
		targetNumber := 10
		if ss.Primary.Class == "BD" {
			targetNumber = 8
		}
		if ss.presenceBT < 0 {
			switch ssg.dp.Sum("2d6") >= targetNumber {
			case true:
				switch bounded(ssg.dp.Sum("2d6")+ggQuantityDM(ss), 6, 12) {
				case 6:
					ss.presenceBT = 1
				case 7, 8, 9, 10, 11:
					ss.presenceBT = 2
				case 12:
					ss.presenceBT = 3
				}
			case false:
				ss.presenceBT = 0
			}
		}
	}
	return nil
}

func (ssg *StarSystemGenerator) setTPquantity(ss *StarSystem) error {
	if ss == nil {
		return fmt.Errorf("no Star System provided")
	}
	ss.presenceBT = ssg.injectedBeltsQuantity
	if ss.presenceBT < 2 {
		ss.presenceTP = ssg.dp.Sum("2d6") - 2 - DeadStars(ss)
		switch ss.presenceTP >= 3 {
		case true:
			ss.presenceTP += ssg.dp.Sum("1d3") - 1
		case false:
			ss.presenceTP = ssg.dp.Sum("1d3") + 2
		}
	}
	return nil
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
