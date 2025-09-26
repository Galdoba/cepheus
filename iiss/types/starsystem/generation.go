package starsystem

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/Galdoba/cepheus/iiss/types/orbit"
	"github.com/Galdoba/cepheus/iiss/types/star"
	"github.com/Galdoba/cepheus/internal/interpolate"
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/float"
)

func (ssg *StarSystemGenerator) GenerateStarOrbits(ss *StarSystem) {
	ss.Orbits = make(map[float64]*orbit.Orbit)
	primary := ss.Stars["Aa"]
	for _, code := range allCodes {
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
		ss.Stars[code].Designation = code
		ss.Stars[code].OrbitN = orb.FromParent
		ss.Stars[code].Eccentricity = orb.Eccentricity
	}
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
	if ss.primary().Class == "BD" {
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
		if ss.primary().Class == "BD" {
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

func companionCode(code string) string {
	switch code {
	case "Aa":
		return "Ab"
	case "Ba":
		return "Bb"
	case "Ca":
		return "Cb"
	case "Da":
		return "Db"
	}
	return ""
}

var allCodes = []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Cb", "Da", "Db"}

func (ssg *StarSystemGenerator) GenerateAllowedOrbits(ss *StarSystem) {
	//rules 1-2
	ss.setMao()
	ss.joinWithCompanionStar()
	//rules 3-11
	ss.setUnavailabilityZones()

}

func (ss *StarSystem) setMao() error {
	for _, st := range ss.Stars {
		st.MinAO = getMAO(st)
	}
	return nil
}

func (ss *StarSystem) joinWithCompanionStar() error {
	for code, st := range ss.Stars {
		if !ss.hasCompanion(code) {
			continue
		}
		compCode := strings.ReplaceAll(code, "a", "b")
		compStar := ss.Stars[compCode]
		maxVal := maxOf(st.MinAO, compStar.MinAO, compStar.Eccentricity+0.5)
		st.MinAO = maxVal
		compStar.MinAO = maxVal
	}
	return nil
}

func maxOf(v ...float64) float64 {
	slices.Sort(v)
	return v[len(v)-1]
}

var secondaryCodes = []string{"Ba", "Ca", "Da"}

func (ss *StarSystem) setUnavailabilityZones() error {
	//rule 3
	ss.primary().MaxAO = 20.0
	ss.primary().AllowedOrbits = orbit.InitialSequance(ss.primary().MinAO, ss.primary().MaxAO)
	//rule 4
	codes := secondaryCodes
	//rule 5
	for _, code := range codes {
		if secondary, ok := ss.Stars[code]; ok {
			width := 1.0
			//rule 6
			if secondary.Eccentricity > 0.2 {
				width = width + 1.0
			}
			//rule 7
			if secondary.Eccentricity > 0.5 && code != "Da" {
				width = width + 1.0
			}
			newAo := orbit.SubtractSubSequence(ss.primary().AllowedOrbits, secondary.OrbitN, width)
			if len(newAo.Segments) > 1 {
				fmt.Println(newAo)

			}
			ss.Stars["Aa"].AllowedOrbits = newAo

			//rule 8
			secondary.MaxAO = secondary.OrbitN - 3.0
			//rule 9
			neibhours := neibhourStars(ss, code)
			neibMod := float64(len(neibhours))
			if neibMod > 0.0 {
				secondary.MaxAO -= neibMod
			}
			//rules 10-11
			if rule1011_apply(neibhours, 0.2) {
				secondary.MaxAO -= 1.0
			}
			if rule1011_apply(neibhours, 0.5) {
				secondary.MaxAO -= 1.0
			}

		}
	}
	return nil
}

func rule1011_apply(nb []*star.Star, eccValue float64) bool {
	apply := false
	for _, st := range nb {
		if st.Eccentricity > eccValue {
			apply = true
		}
	}
	return apply
}

func neibhourStars(ss *StarSystem, code string) []*star.Star {
	strs := []*star.Star{}
	if _, ok := ss.Stars[code]; ok {

		switch code {
		case "Ba":
			if nb, ok := ss.Stars["Ca"]; ok {
				strs = append(strs, nb)
			}
		case "Ca":
			if nb, ok := ss.Stars["Ba"]; ok {
				strs = append(strs, nb)
			}
			if nb, ok := ss.Stars["Da"]; ok {
				strs = append(strs, nb)
			}
		case "Da":
			if nb, ok := ss.Stars["Ca"]; ok {
				strs = append(strs, nb)
			}
		}
	}
	return strs
}

func (ss *StarSystem) hasCompanion(code string) bool {
	_, ok := ss.Stars[companionCode(code)]
	return ok
}

func getMAO(st *star.Star) float64 {
	index := st.Index()
	return interpolate.MAO_ByIndex(index)
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
	ssg.GenerateAllowedOrbits(ss)
	ssg.CalculateHZCO(ss)
	return ss
}

func (ssg *StarSystemGenerator) CalculateHZCO(ss *StarSystem) error {
	for _, code := range []string{"Aa", "Ba", "Ca", "Da"} {
		if _, ok := ss.Stars[code]; !ok {
			continue
		}
		st := ss.Stars[code]
		luma := st.Luminocity
		if comp, ok := ss.Stars[companionCode(code)]; ok {
			luma += comp.Luminocity
		}
		hzco := float.Round(math.Sqrt(luma))
		st.HZCO = hzco
		if comp, ok := ss.Stars[companionCode(code)]; ok {
			comp.HZCO = hzco
		}

	}
	return nil
}
