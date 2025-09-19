package starsystem

import "github.com/Galdoba/cepheus/iiss/types/star"

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

func ggQuantityDM(ss *StarSystem) int {
	dm := 0
	if len(ss.Stars) == 1 && ss.Primary.Class == "V" {
		dm = dm + 1
	}
	switch ss.Primary.Class {
	case "BD", "D", "L", "T", "Y", "PSR", "NS", "BH":
		dm = dm - 2
	}
	for _, s := range ss.Stars {
		switch s.Class {
		case "D", "L", "T", "Y", "PSR", "NS", "BH":
			dm = dm - 1
		}
	}
	if len(ss.Stars) >= 4 {
		dm = dm - 1
	}
	return dm
}

func beltQuantityDM(ss *StarSystem) int {
	dm := 0
	if ss.presenceGG >= 1 {
		dm = dm + 1
	}
	if IsProtoStarSystem(ss) {
		dm = dm + 3
	}
	if IsPrimordialSystem(ss) {
		dm = dm + 2
	}

	switch ss.Primary.Class {
	case "BH", "D", "PSR", "NS":
		dm = dm + 1
	}

	dm = dm + DeadStars(ss)

	if len(ss.Stars) >= 2 {
		dm = dm + 1
	}

	return dm
}

func IsPrimordialSystem(ss *StarSystem) bool {
	switch ss.Primary.Class {
	case "Ia", "Ib", "II":
		return true
	}
	if ss.Primary.Type == "O" {
		return true
	}
	switch ss.Primary.Mass > 8.0 {
	case true:
		if ss.Primary.Age <= star.MainSequanceLifespan(ss.Primary.Mass) {
			return true
		}
	case false:
	}
	if ss.Primary.Age >= 0.01 && ss.Primary.Age <= 0.1 {
		return true
	}
	return false
}

func IsProtoStarSystem(ss *StarSystem) bool {
	proto := false
	for _, st := range ss.Stars {
		if st.ProtoStar {
			proto = true
		}
	}
	return proto
}

func DeadStars(ss *StarSystem) int {
	dead := 0
	for _, st := range ss.Stars {
		switch st.Class {
		case "D", "NS", "PSR", "BH":
			dead++
		}
	}
	return dead
}
