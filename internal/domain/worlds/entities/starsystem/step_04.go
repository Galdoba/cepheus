package starsystem

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
)

const (
	RocheLimit = 0.01 //au
)

// - [ ] 4. **Determine allowable planetary Orbit#s**
//   - [+] a. Determine a star's minimum allowable Orbit# from table
//   - [ ] b. In multiple star systems, follow process to exclude orbital ranges
func (builder *Builder) runStep04(systemPrecursor *starSystemPrecursor) error {
	slots := starsOrbitalSlots(orbit.Orbit{}, systemPrecursor.PrimaryStar)
	closeOrb := orbit.Orbit{}
	nearOrb := orbit.Orbit{}
	farOrb := orbit.Orbit{}
	for orb, star := range systemPrecursor.Stars {
		switch star.Designation {
		default:
			continue
		case "Aa":
			slots.removeCompanionInfluence(orb)
		case "Ba", "Ca", "Da":
			slots.removeSecondaryInfluence(orb)
			systemPrecursor.Stars[orb].OrbitalSlots = starsOrbitalSlots(orb, star)
			fmt.Println("add")
			fmt.Println(systemPrecursor.Stars[orb])
			switch orb.Designation {
			case "Ba":
				closeOrb = orb
			case "Ca":
				nearOrb = orb
			case "Da":
				farOrb = orb
			}
		}
	}
	closeSrink := 0.0
	nearASrink := 0.0
	nearBSrink := 0.0
	farSrink := 0.0
	if nearOrb.Distance != 0 {
		if closeOrb.Distance != 0.0 {
			for _, ecc := range []float64{0.0, 0.2, 0.5} {
				if nearOrb.Eccentricity >= ecc || closeOrb.Eccentricity >= ecc {
					closeSrink += 1
					nearASrink += 1
				}
			}
		}
		if farOrb.Distance != 0.0 {
			for _, ecc := range []float64{0.0, 0.2, 0.5} {
				if nearOrb.Eccentricity >= ecc || farOrb.Eccentricity >= ecc {
					farSrink += 1
					nearBSrink += 1
				}
			}
		}
	}
	si := newStarIterator(systemPrecursor.Stars)
	for si.next() {
		orb, star, _ := si.getValues()
		switch star.Designation {
		case "Aa":
			systemPrecursor.Stars[orb].OrbitalSlots = slots
		case "Ba":
			systemPrecursor.Stars[orb].OrbitalSlots.shrinkMaximum(closeSrink)
		case "Ca":
			systemPrecursor.Stars[orb].OrbitalSlots.shrinkMaximum(max(nearASrink, nearBSrink))
		case "Da":
			systemPrecursor.Stars[orb].OrbitalSlots.shrinkMaximum(farSrink)
		}
	}

	return validateStep04Results(systemPrecursor)
}

func validateStep04Results(systemPrecursor *starSystemPrecursor) error {
	si := newStarIterator(systemPrecursor.Stars)
	for si.next() {
		_, star, _ := si.getValues()
		if strings.Contains(string(star.Designation), "b") {
			continue
		}
		// fmt.Println(star.OrbitalSlots)

		fmt.Println("validation slots:", len(star.OrbitalSlots.points))
	}
	return nil
}

func ownStarMinimumAllowableOrbit(star *starPrecursor) float64 {
	return float.RoundN(star.Diameter*RocheLimit, 2)
}

func ownStarMaximumAllowableOrbit(starDistance float64) float64 {
	if starDistance == 0 {
		return 20
	}
	return max(0, starDistance-3)
}

func starsOrbitalSlots(orb orbit.Orbit, star *starPrecursor) *orbitalSlots {
	minAO := ownStarMinimumAllowableOrbit(star)
	maxAO := ownStarMaximumAllowableOrbit(orb.Distance)
	if star.Class == "D" {
		maxAO = maxAO * 0.25
	}
	return newSlots(minAO, maxAO)
}

type orbitalSlots struct {
	points map[float64]string
}

func newSlots(min, max float64) *orbitalSlots {
	os := orbitalSlots{make(map[float64]string)}
	for i := min; i <= max; i += 0.01 {
		val := float.RoundN(i, 2)
		os.points[val] = "allowed"
	}
	return &os
}

func (os *orbitalSlots) punch(center, radius float64) {
	for point := range os.points {
		if point > center-radius && point < center+radius {
			delete(os.points, point)
		}
	}
}

func (os *orbitalSlots) shrinkMaximum(val float64) {
	if val == 0 {
		return
	}
	max := 0.0
	for point := range os.points {
		if point > max {
			max = point
		}
	}
	for point := range os.points {
		if point > max-val {
			delete(os.points, point)
		}
	}
}

// 2
func (os *orbitalSlots) removeCompanionInfluence(companionOrbit orbit.Orbit) {
	max := 0.5 + companionOrbit.Eccentricity
	for point := range os.points {
		if point <= max {
			delete(os.points, point)
		}
	}
}

// 3,4,5
func (os *orbitalSlots) removeSecondaryInfluence(secondaryOrbit orbit.Orbit) {
	radius := 1.0
	if secondaryOrbit.Eccentricity > 0.2 {
		radius += 1.0
	}
	if secondaryOrbit.Eccentricity > 0.5 && secondaryOrbit.Designation != "Da" {
		radius += 1.0
	}
	os.punch(secondaryOrbit.Distance, radius)
}
