package starsystem

import (
	"fmt"
	"math"
	"strings"

	"github.com/Galdoba/cepheus/iiss/types/orbit"
	"github.com/Galdoba/cepheus/iiss/types/star"
	"github.com/Galdoba/cepheus/internal/interpolate"
	"github.com/Galdoba/cepheus/pkg/dice"
)

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
		ss.Stars[code].Designation = code
		ss.Stars[code].OrbitN = orb.FromParent
		if ss.Stars[code].OrbitN == 0 {
			ss.Stars[code].OrbitN = 0.01
		}
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

func (ssg *StarSystemGenerator) GenerateAllowedOrbits(ss *StarSystem) {
	codes := []string{"Aa", "Ba", "Ca", "Da"}

	for _, code := range codes {
		if _, ok := ss.Stars[code]; !ok {
			continue
		}
		st := ss.Stars[code]
		mao := st.MinimumAllowedOrbit
		if comp, ok := ss.Stars[companionCode(code)]; ok {
			mao += comp.Eccentricity
		}
		max := 20.0
		if st.OrbitN > 0.5 {
			max = st.OrbitN - 3.0
		}
		if max < mao {
			max = mao
		}
		st.AllowedOrbits = orbit.InitialSequance(mao, max)

	}
	for _, code := range []string{"Ba", "Ca", "Da"} {
		if _, ok := ss.Stars[code]; !ok {
			continue
		}
		st := ss.Stars[code]
		mao := st.MinimumAllowedOrbit
		if comp, ok := ss.Stars[companionCode(code)]; ok {
			mao += comp.Eccentricity
		}
		max := 20.0
		if st.OrbitN > 0.5 {
			max = st.OrbitN - 3.0
		}
		if max < mao {
			max = mao
		}
		st.AllowedOrbits = orbit.InitialSequance(mao, max)

	}

	// primary := ss.Primary
	// primMao := primary.MinimumAllowedOrbit
	// if primComp, ok := ss.Stars["Ab"]; ok {
	// 	primMao += primComp.Eccentricity
	// }
	// primary.AllowedOrbits = orbit.InitialSequance(primary.MinimumAllowedOrbit, 20.0)
	// for _, code := range codes {
	// 	if _, ok := ss.Stars[code]; !ok {
	// 		continue
	// 	}
	// 	st := ss.Stars[code]
	// 	st.MinimumAllowedOrbit = getMAO(st)
	// 	ss.Stars["Aa"].AllowedOrbits = calculatePrimaryOrbits(ss)
	// 	if code == "Aa" {
	// 		for _, code2 := range []string{"Ba", "Ca", "Da"} {
	// 			if _, ok := ss.Stars[code2]; !ok {
	// 				continue
	// 			}

	// 			secondary := ss.Stars[code2]

	// 			fmt.Println(code2, "===>", secondary)
	// 		}
	// 	}

	// for _, code := range codes {
	// 	if ss.Stars[code] == nil {
	// 		continue
	// 	}
	// 	switch code {
	// 	case "Aa":
	// 		ss.Stars[code].AllowedOrbits = calculatePrimaryOrbits(ss)
	// 	default:
	// 		allowed := calculateSecondaryOrbits(ss.Stars[code], ss)
	// 		fmt.Println(allowed)
	// 		ss.Stars[code].AllowedOrbits = allowed
	// 	}
	// 	// }
	// }
}
func getMAO(st *star.Star) float64 {
	index := st.Index()
	return interpolate.MAO_ByIndex(index)
}

// calculatePrimaryOrbits вычисляет доступные орбиты вокруг главной звезды
func calculatePrimaryOrbits(ss *StarSystem) []map[string]float64 {
	// Начальные доступные орбиты: от MAO первичной звезды до 20.0
	primaryMAO := getMAO(ss.Primary)
	available := []map[string]float64{{fmt.Sprintf("%v start", "Aa"): primaryMAO, fmt.Sprintf("%v end", "Aa"): 20.0}}

	// Применить правила для компаньонов

	for code, comp := range ss.Stars {
		if strings.Contains(code, "a") {
			continue
		}
		minOrbit := 0.5 + comp.Eccentricity
		available = subtractRange(available, 0, minOrbit)
	}

	// Применить правила для вторичных звёзд
	for code, sec := range ss.Stars {
		if strings.Contains(code, "a") {
			continue
		}
		if strings.Contains(code, "b") {
			continue
		}
		// Базовая запретная зона
		zoneStart := sec.OrbitN - 1.0
		zoneEnd := sec.OrbitN + 1.0

		// Расширить зону если ecc > 0.2
		if sec.Eccentricity > 0.2 {
			zoneStart -= 1.0
			zoneEnd += 1.0
		}

		// Дополнительное расширение если ecc > 0.5 и не Far звезда
		// Note: В реализации нужно определить, является ли звезда Far
		// Для простоты предположим, что у нас есть метод определения
		if sec.Eccentricity > 0.5 && !isFarStar(sec) {
			zoneStart -= 1.0
			zoneEnd += 1.0
		}

		available = subtractRange(available, zoneStart, zoneEnd)
	}

	return available
}

func calculateSecondaryOrbits(sec *star.Star, sys *StarSystem) []map[string]float64 {
	// Базовый максимум
	orb := sec.OrbitN

	maxOrbit := orb - 3.0

	// Применить модификаторы
	hasCloseNeighbor := hasAdjacentStar(sys, "Close")
	hasFarNeighbor := hasAdjacentStar(sys, "Far")

	// Правило 9: соседи в смежных зонах
	if (hasCloseNeighbor && hasFarNeighbor) ||
		(hasCloseNeighbor && orb > 5.0) || // Примерная логика для Near
		(hasFarNeighbor && orb < 10.0) { // Примерная логика для Near
		maxOrbit -= 1.0
	}

	// Правило 10: высокий эксцентриситет у соседей
	if hasHighEccNeighbor(sys, sec, 0.2) {
		maxOrbit -= 1.0
	}

	// Правило 11: очень высокий собственный эксцентриситет
	if sec.Eccentricity > 0.5 {
		maxOrbit -= 1.0
	}

	// Вернуть доступный диапазон
	secMAO := getMAO(sec)
	if maxOrbit <= secMAO {
		return []map[string]float64{} // Нет доступных орбит
	}

	return []map[string]float64{{"start": secMAO, "end": maxOrbit}}
}

func hasAdjacentStar(sys *StarSystem, zone string) bool {
	// Проверка наличия звезды в указанной зоне
	for _, sec := range sys.Stars {
		orb := sec.OrbitN
		if (zone == "Close" && orb < 5.0) ||
			(zone == "Near" && orb >= 5.0 && orb <= 11.0) ||
			(zone == "Far" && orb > 11.0) {
			return true
		}
	}
	return false
}

func hasHighEccNeighbor(sys *StarSystem, current *star.Star, threshold float64) bool {
	// Проверка наличия соседа с высоким эксцентриситетом
	for code, st := range sys.Stars {
		if strings.Contains(code, "b") || strings.Contains(code, "A") {
			continue
		}
		orb1, orb2 := st.OrbitN, current.OrbitN

		if st.Designation != current.Designation && math.Abs(orb1-orb2) < 5.0 {
			if st.Eccentricity > threshold {
				return true
			}
		}
	}
	return false
}

func subtractRange(ranges []map[string]float64, start, end float64) []map[string]float64 {
	var result []map[string]float64

	for _, r := range ranges {
		rStart := r["start"]
		rEnd := r["end"]

		// Проверяем, есть ли перекрытие диапазонов
		if end <= rStart || start >= rEnd {
			// Нет перекрытия - добавляем весь диапазон
			result = append(result, r)
			continue
		}

		// Есть перекрытие - разделяем на части
		if rStart < start {
			result = append(result, map[string]float64{"start": rStart, "end": start})
		}
		if rEnd > end {
			result = append(result, map[string]float64{"start": end, "end": rEnd})
		}
	}

	return result
}

func isFarStar(sec *star.Star) bool {
	// Реализация определения, является ли звезда Far
	// Это упрощённая версия - в реальности нужно учитывать систему классификации
	return sec.OrbitN > 11.0
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
	return ss
}
