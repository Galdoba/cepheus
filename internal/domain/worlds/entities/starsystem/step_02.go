package starsystem

import (
	"fmt"
	"strconv"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/au"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
)

func (b *Builder) runStep2(ss *StarSystem) error {
	// - [+] 2. **Determine if system has multiple stars, if yes, then:**
	//   - [+] a. Determine Orbit#s of secondary and companion stars
	b.step2.starSchema = stellar.RollDesignations(b.rng)
	if b.imported.Allegiance != "" && b.imported.Stellar != "" {
		b.step2.starSchema = stellar.RollStellarDesignations(b.rng, stellar.Stellar(b.imported.Stellar))
	}
	ss.Stars = make(map[orbit.Orbit]*Star)
	switch len(b.step2.starSchema) {
	case 0:
		return fmt.Errorf("failed to roll stellar designations")
	case 1:
		ss.Stars[orbit.RollStellarOrbit(b.rng, b.step2.starSchema[0])] = ss.PrimaryStar
		b.step2.completed = true
		return nil
	default:
		ss.Stars[orbit.RollStellarOrbit(b.rng, b.step2.starSchema[0])] = ss.PrimaryStar
		for _, designation := range b.step2.starSchema {
			switch designation {
			case stellar.Primary:
				ss.Stars[orbit.RollStellarOrbit(b.rng, stellar.StarDesignation(designation))] = ss.PrimaryStar
			default:
				ss.Stars[orbit.RollStellarOrbit(b.rng, stellar.StarDesignation(designation))] = &Star{Designation: designation}
			}
		}
	}
	//   - [+] b. Determine eccentricity of secondary stars and check for overlaps
	if err := b.determineEccentricityOfSecondaryStars(ss); err != nil {
		return err
	}
	//   - [+] c. Determine secondary and companion star types
	if err := b.determineSecondaryTypeAndClass(ss); err != nil {
		return err
	}
	//   - [ ] d. Adjust system age to account for post-stellar objects (if any)
	//   - [ ] e. Determine star orbital periods
	if err := b.determineMassOfSecondaryStars(ss); err != nil {
		return err
	}
	if err := b.determineStarsOrbitalPeriods(ss); err != nil {
		return err
	}
	if err := step2Validation(ss); err != nil {
		return err
	}
	b.step2.completed = true
	return nil
}

func step2Validation(ss *StarSystem) error {
	switch len(ss.Stars) {
	case 0:
		return fmt.Errorf("no stars found")
	case 1:
		return nil
	default:
		for o, star := range ss.Stars {
			if o.Distance < 0.02 && star.Designation != stellar.Primary {
				return fmt.Errorf("star distance is lower than expected: %v : %v", o, star)
			}
			if o.Distance > 20 && star.Designation != stellar.Primary {
				return fmt.Errorf("star distance is higher than expected: %v : %v", o, star)
			}
			if o.Eccentricity < 0 && star.Designation != stellar.Primary {
				return fmt.Errorf("star eccentrisity is not set: %v : %v", o, star)
			}
			if star.Mass == 0 {
				return fmt.Errorf("no mass for: %v", star)
			}
			if star.Designation != stellar.Primary && star.Period == 0 {
				return fmt.Errorf("no period for: %v", star)
			}
			// fmt.Printf("star: %v; orbit:= %v\n", star, o)
		}
	}
	PrintStarPositions(ss)

	return nil
}

func (b *Builder) determineMassOfSecondaryStars(ss *StarSystem) error {
	si := newStarIterator(ss.Stars)
	for si.next() {
		o, star, err := si.getValues()
		if err != nil {
			return err
		}
		if star.Mass == 0 {
			if err := determineMassDiameterAgeTemperature(b.rng, star); err != nil {
				return err
			}
		}
		if star.Mass == 0 {
			return fmt.Errorf("failed to calcaulate mass for %v", star)
		} else {
			ss.Stars[o] = star
		}

	}
	return nil
}

func (b *Builder) determineStarsOrbitalPeriods(ss *StarSystem) error {
	si := newStarIterator(ss.Stars)
	for si.next() {
		o, star, err := si.getValues()
		if err != nil {
			return err
		}
		if star.Designation == stellar.Primary {
			continue
		}
		parent := getParent(ss, o)
		dist := o.Distance
		distAU := au.FromOrbitNumber(dist)
		star.Period = float.RoundN(orbit.StarOrbitPeriod(parent.Mass, star.Mass, distAU), 3)
		ss.Stars[o] = star

	}
	return nil
}

func (b *Builder) determineEccentricityOfSecondaryStars(ss *StarSystem) error {
	si := newStarIterator(ss.Stars)
	for si.next() {
		o, star, err := si.getValues()
		if err != nil {
			return err
		}
		dm := orbit.DM_ObjectIsStar
		if o.Distance <= 0 {
			continue
		}
		switch star.Designation {
		case stellar.PrimaryComp, stellar.CloseComp, stellar.NearComp, stellar.FarComp:
			dm += orbit.DM_PerParentObject
		case stellar.Close, stellar.Near, stellar.Far:
			for i, des := range b.step2.starSchema {
				if des == star.Designation {
					dm += (orbit.DM_PerParentObject * i)
				}
			}
		}
		updatedOrbit := o
		updatedOrbit.Eccentricity = orbit.RollStarEccentricity(b.rng, dm)
		ss.Stars[updatedOrbit] = star
		delete(ss.Stars, o)
	}

	return nil
}

func (b *Builder) determineSecondaryTypeAndClass(ss *StarSystem) error {
	si := newStarIterator(ss.Stars)
	for si.next() {
		o, star, err := si.getValues()
		if err != nil {
			return nil
		}
		if star.Designation == stellar.Primary {
			continue
		}
		secondaryType, err := b.step2.tables.Roll(rtg.TableSecondary, ss.PrimaryStar.Class)
		if err != nil {
			return err
		}
		parent := getParent(ss, o)
		// fmt.Println("secondary:", star.Designation, secondaryType)
		switch secondaryType {
		default:
			panic(secondaryType + " not implemented")
		case "Sibling":
			if err := makeSibling(b.rng, parent, star); err != nil {
				return fmt.Errorf("failed to make sibling: %v", err)
			}
		case "Twin":
			if err := makeTwin(b.rng, parent, star); err != nil {
				return fmt.Errorf("failed to make sibling: %v", err)
			}
		case "Lesser", "Random":
			if err := makeLesser(b.rng, parent, star); err != nil {
				return fmt.Errorf("failed to make sibling: %v", err)
			}
		case "BD", "D", "NS":
			star.Class = ""
			star.Type = secondaryType
			star.SubType = ""

		}
		switch star.Type {
		case "O", "B", "A", "F", "G", "K", "M":
		default:
			star.SubType = ""
			star.Class = ""
		}
		fmt.Println(star)
		if err := validateTSC(star); err != nil {
			return err
		}

		ss.Stars[o] = star
	}

	return nil
}

func getParent(ss *StarSystem, o orbit.Orbit) *Star {
	parentDes := stellar.StarDesignation(o.Designation)
	return ss.getStarByDesignation(parentDes)

}

func (ss *StarSystem) getStarByDesignation(d stellar.StarDesignation) *Star {
	for _, star := range ss.Stars {
		if star.Designation == d {
			return star
		}
	}
	return nil
}

func makeSibling(r stellar.Roller, parent, child *Star) error {
	if parent == nil {
		return fmt.Errorf("no parent provided")
	}
	if child == nil {
		return fmt.Errorf("no child provided")
	}
	switch parent.SubType {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		n, _ := strconv.Atoi(parent.SubType)
		if n == 0 {
			n = 10
		}
		n = n - r.Roll("1d6")
		switch n > 0 {
		case true:
			child.Type = parent.Type
			child.SubType = strconv.Itoa(n)
		case false:
			child.Type = coolerType(parent.Type)
			child.SubType = parent.SubType
		}
	default:
		child.Type = parent.Type
	}

	return nil
}

func makeTwin(r stellar.Roller, parent, child *Star) error {
	if parent == nil {
		return fmt.Errorf("no parent provided")
	}
	if child == nil {
		return fmt.Errorf("no child provided")
	}
	child.Class = parent.Class
	child.Type = parent.Type
	child.SubType = parent.SubType
	child.Mass = float.RoundN(parent.Mass*(1-(0.01*float64(r.Roll("1d6")))), 3)
	return nil
}

func makeLesser(r stellar.Roller, parent, child *Star) error {
	if parent == nil {
		return fmt.Errorf("no parent provided")
	}
	if child == nil {
		return fmt.Errorf("no child provided")
	}
	fmt.Println(parent)
	child.Class = parent.Class
	child.Type = coolerType(parent.Type)
	n := r.Roll("1d10-1")
	child.SubType = strconv.Itoa(n)
	m, _ := strconv.Atoi(parent.SubType)
	if child.Type == "M" && parent.Type == "M" && (m+1) < n {
		child.SubType = ""
		child.Type = "BD"
		child.Class = ""
	}
	return nil
}

func makeRandom(r stellar.Roller, parent, child *Star) error {
	if parent == nil {
		return fmt.Errorf("no parent provided")
	}
	if child == nil {
		return fmt.Errorf("no child provided")
	}
	child.Class = parent.Class
	child.Type = coolerType(parent.Type)
	n := r.Roll("1d10")
	child.SubType = strconv.Itoa(n)
	m, _ := strconv.Atoi(parent.SubType)
	if child.Type == "M" && parent.Type == "M" && (m+1) < n {
		child.SubType = ""
		child.Type = ""
		child.Class = "BD"
	}
	return nil
}

func coolerType(stype string) string {
	switch stype {
	case "O":
		return "B"
	case "B":
		return "A"
	case "A":
		return "F"
	case "F":
		return "G"
	case "G":
		return "K"
	case "K":
		return "M"
	}
	return stype
}

func orderedDesignations(stars map[orbit.Orbit]*Star) []stellar.StarDesignation {
	designations := []stellar.StarDesignation{}
	for _, d := range stellar.AllDesignations() {
		for _, star := range stars {
			if star.Designation == d {
				designations = append(designations, d)
			}
		}
	}
	return designations
}

type starIterator struct {
	pool         map[orbit.Orbit]*Star
	designations []stellar.StarDesignation
	index        int
	hasMore      bool
}

func newStarIterator(stars map[orbit.Orbit]*Star) *starIterator {
	si := starIterator{}
	si.pool = stars
	si.designations = stellar.AllDesignations()
	si.index = -1
	return &si
}

func (si *starIterator) next() bool {
	for i := si.index + 1; i < len(si.designations); i++ {
		for _, star := range si.pool {
			if star.Designation == si.designations[i] {
				si.index = i
				si.hasMore = true
				return true
			}
		}
	}
	si.hasMore = false
	return false
}

func (si *starIterator) getValues() (orbit.Orbit, *Star, error) {
	if !si.hasMore {
		return orbit.Orbit{}, nil, fmt.Errorf("no more stars")
	}
	for o, star := range si.pool {
		if star.Designation == si.designations[si.index] {
			return o, star, nil
		}
	}
	return orbit.Orbit{}, nil, fmt.Errorf("unexpected end")
}

func PrintStarPositions(ss *StarSystem) {
	si := newStarIterator(ss.Stars)
	for si.next() {
		o, star, err := si.getValues()
		if err != nil {
			panic(err)
		}
		fmt.Printf("orbit %v around %v (%v) star: '%v%v %v'\n", star.Designation, o.Designation, o.Distance, star.Type, star.SubType, star.Class)
	}
}
