package starsystem

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/au"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
)

// runStep2 executes Step 2 of star system generation: determining secondary stars and their properties.
// This includes:
// 1. Rolling for stellar designations (how many stars and their roles)
// 2. Determining orbital eccentricities for companion stars
// 3. Determining star types for secondary stars (sibling, twin, lesser)
// 4. Calculating physical properties for all stars
// 5. Determining orbital periods
// 6. Adjusting system age based on post-stellar objects
func (builder *Builder) runStep02(systemPrecursor *starSystemPrecursor) error {
	builder.step2.starSchema = stellar.RollDesignations(builder.rng)
	if builder.imported.Allegiance != "" && builder.imported.Stellar != "" {
		builder.step2.starSchema = stellar.RollStellarDesignations(builder.rng, stellar.Stellar(builder.imported.Stellar))
	}
	systemPrecursor.Stars = make(map[orbit.Orbit]*starPrecursor)
	switch len(builder.step2.starSchema) {
	case 0:
		return fmt.Errorf("failed to roll stellar designations")
	case 1:
		systemPrecursor.Stars[orbit.RollStellarOrbit(builder.rng, builder.step2.starSchema[0])] = systemPrecursor.PrimaryStar
		builder.step2.completed = true
		return nil
	default:
		systemPrecursor.Stars[orbit.RollStellarOrbit(builder.rng, builder.step2.starSchema[0])] = systemPrecursor.PrimaryStar
		for _, designation := range builder.step2.starSchema {
			switch designation {
			case stellar.Primary:
				systemPrecursor.Stars[orbit.RollStellarOrbit(builder.rng, stellar.StarDesignation(designation))] = systemPrecursor.PrimaryStar
			default:
				systemPrecursor.Stars[orbit.RollStellarOrbit(builder.rng, stellar.StarDesignation(designation))] = &starPrecursor{Designation: designation}
			}
		}
	}
	if err := builder.rollOrbitalEccentricitiesForCompanionStars(systemPrecursor); err != nil {
		return err
	}
	if err := builder.determineCompanionStarClassification(systemPrecursor); err != nil {
		return err
	}
	if err := builder.determineSecondaryStarsDetails(systemPrecursor); err != nil {
		return err
	}
	if err := builder.determineStarsOrbitalPeriods(systemPrecursor); err != nil {
		return err
	}
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		_, starPrecursor, err := starIterator.getValues()
		if err != nil {
			panic(err) // TODO: Return error instead of panic
		}
		// TODO: Implement step 2d - Adjust system age to account for post-stellar objects
		if starPrecursor.Dead && starPrecursor.Age > systemPrecursor.Age {
			fmt.Println("adjust system age:", systemPrecursor.Age, starPrecursor.Age)
			// TODO: Remove debug sleep - this was likely for debugging
			time.Sleep(time.Second)
			systemPrecursor.Age = starPrecursor.Age
		}
	}

	if err := validateStep02Results(systemPrecursor); err != nil {
		return err
	}
	builder.step2.completed = true
	return nil
}

// validateStep2Results validates the results of step 2 generation.
// Checks that:
// - At least one star exists
// - Secondary stars have valid distances (between 0.02 and 20)
// - Eccentricities are set for secondary stars
// - All stars have mass
// - All secondary stars have orbital periods
func validateStep02Results(systemPrecursor *starSystemPrecursor) error {
	switch len(systemPrecursor.Stars) {
	case 0:
		return fmt.Errorf("no stars found")
	case 1:
		return nil
	default:
		for orbit, starPrecursor := range systemPrecursor.Stars {
			if orbit.Distance < 0.02 && starPrecursor.Designation != stellar.Primary {
				return fmt.Errorf("star distance is lower than expected: %v : %v", orbit, starPrecursor)
			}
			if orbit.Distance > 20 && starPrecursor.Designation != stellar.Primary {
				return fmt.Errorf("star distance is higher than expected: %v : %v", orbit, starPrecursor)
			}
			if orbit.Eccentricity < 0 && starPrecursor.Designation != stellar.Primary {
				return fmt.Errorf("star eccentricity is not set: %v : %v", orbit, starPrecursor)
			}
			if starPrecursor.Mass == 0 {
				return fmt.Errorf("no mass for: %v", starPrecursor)
			}
			if starPrecursor.Designation != stellar.Primary && starPrecursor.Period == 0 {
				return fmt.Errorf("no period for: %v", starPrecursor)
			}
		}
	}

	return nil
}

// determineSecondaryStarsDetails calculates physical properties (mass, diameter, temperature, age, luminosity)
// for all secondary stars in the system. The primary star should already have these properties from step 1.
func (builder *Builder) determineSecondaryStarsDetails(systemPrecursor *starSystemPrecursor) error {
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		orbit, starPrecursor, err := starIterator.getValues()
		if err != nil {
			return err
		}
		// Calculate physical properties if not already set
		if starPrecursor.Mass == 0 {
			if err := calculateStarPhysicalProperties(builder.rng, starPrecursor); err != nil {
				return err
			}
		}
		if starPrecursor.Mass == 0 {
			return fmt.Errorf("failed to calculate mass for %v", starPrecursor)
		} else {
			systemPrecursor.Stars[orbit] = starPrecursor
		}
		// Calculate luminosity if not already set
		if starPrecursor.Luminosity == 0 {
			starPrecursor.Luminosity = float.RoundN(calculateLuminosity(starPrecursor.Diameter, starPrecursor.Temperature), 3)
			if starPrecursor.Luminosity < 0.001 {
				starPrecursor.Luminosity = 0.001
			}
		}
		// Duplicate check - could be consolidated with above
		if starPrecursor.Mass == 0 {
			return fmt.Errorf("failed to calculate mass for %v", starPrecursor)
		} else {
			systemPrecursor.Stars[orbit] = starPrecursor
		}

	}
	return nil
}

// determineStarsOrbitalPeriods calculates the orbital period for each secondary star in the system.
// The orbital period is the time it takes for the star to complete one orbit around its parent.
// Uses Kepler's third law modified for the system.
func (builder *Builder) determineStarsOrbitalPeriods(systemPrecursor *starSystemPrecursor) error {
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		orbitInstance, starPrecursor, err := starIterator.getValues()
		if err != nil {
			return err
		}
		if starPrecursor.Designation == stellar.Primary {
			continue
		}
		parentStar := findParentStar(systemPrecursor, orbitInstance)
		dist := orbitInstance.Distance
		distAU := au.FromOrbitNumber(dist)
		starPrecursor.Period = float.RoundN(orbit.StarOrbitPeriod(parentStar.Mass, starPrecursor.Mass, distAU), 3)
		systemPrecursor.Stars[orbitInstance] = starPrecursor

	}
	return nil
}

// rollOrbitalEccentricitiesForCompanionStars rolls for the orbital eccentricity of each secondary star.
// Eccentricity determines how elliptical an orbit is (0 = circular, closer to 1 = more elliptical).
// Different dice modifiers are applied based on whether the companion is a primary companion
// (Close, Near, Far) or a companion of companion.
func (builder *Builder) rollOrbitalEccentricitiesForCompanionStars(systemPrecursor *starSystemPrecursor) error {
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		orbitInstance, starPrecursor, err := starIterator.getValues()
		if err != nil {
			return err
		}
		diceModifier := orbit.DM_ObjectIsStar
		if orbitInstance.Distance <= 0 {
			continue
		}
		switch starPrecursor.Designation {
		case stellar.PrimaryComp, stellar.CloseComp, stellar.NearComp, stellar.FarComp:
			diceModifier += orbit.DM_PerParentObject
		case stellar.Close, stellar.Near, stellar.Far:
			for i, des := range builder.step2.starSchema {
				if des == starPrecursor.Designation {
					diceModifier += (orbit.DM_PerParentObject * i)
				}
			}
		}
		updatedOrbit := orbitInstance
		updatedOrbit.Eccentricity = float.RoundN(orbit.RollStarEccentricity(builder.rng, diceModifier), 3)
		systemPrecursor.Stars[updatedOrbit] = starPrecursor
		delete(systemPrecursor.Stars, orbitInstance)
	}

	return nil
}

// determineCompanionStarClassification determines the classification (type) of each secondary star
// in a multiple star system. Secondary stars can be classified as:
// - Sibling: Similar type to parent but cooler
// - Twin: Same type and subtype as parent
// - Lesser: Cooler and often smaller
// - BD, D, NS: Brown dwarf, white dwarf, or neutron star
func (builder *Builder) determineCompanionStarClassification(systemPrecursor *starSystemPrecursor) error {
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		orbit, starPrecursor, err := starIterator.getValues()
		if err != nil {
			return nil // TODO: Return error instead of nil
		}
		if starPrecursor.Designation == stellar.Primary {
			continue
		}
		secondaryType, err := builder.step2.tables.Roll(rtg.TableSecondary, systemPrecursor.PrimaryStar.Class)
		if err != nil {
			return err
		}
		parentStar := findParentStar(systemPrecursor, orbit)
		switch secondaryType {
		default:
			// TODO: Return error instead of panic for unimplemented types
			panic(secondaryType + " not implemented")
		case "Sibling":
			if err := createSiblingStar(builder.rng, parentStar, starPrecursor); err != nil {
				return fmt.Errorf("failed to create sibling star: %v", err)
			}
		case "Twin":
			if err := createTwinStar(builder.rng, parentStar, starPrecursor); err != nil {
				return fmt.Errorf("failed to create twin star: %v", err)
			}
		case "Lesser", "Random":
			if err := createLesserCompanionStar(builder.rng, parentStar, starPrecursor); err != nil {
				return fmt.Errorf("failed to create lesser companion star: %v", err)
			}
		case "BD", "D", "NS":
			starPrecursor.Class = ""
			starPrecursor.Type = secondaryType
			starPrecursor.SubType = ""

		}
		switch starPrecursor.Type {
		case "O", "B", "A", "F", "G", "K", "M":
		default:
			starPrecursor.SubType = ""
			starPrecursor.Class = ""
		}
		if err := validateStarTypeSubtypeClassCombination(starPrecursor); err != nil {
			return err
		}

		systemPrecursor.Stars[orbit] = starPrecursor
	}

	return nil
}

// findParentStar finds the parent star that a companion star orbits.
// The parent is determined by the orbit's designation.
func findParentStar(systemPrecursor *starSystemPrecursor, orbit orbit.Orbit) *starPrecursor {
	parentDesignation := stellar.StarDesignation(orbit.Designation)
	return systemPrecursor.findStarByDesignation(parentDesignation)

}

// findStarByDesignation searches for a star in the system by its designation.
func (systemPrecursor *starSystemPrecursor) findStarByDesignation(designation stellar.StarDesignation) *starPrecursor {
	for _, starPrecursor := range systemPrecursor.Stars {
		if starPrecursor.Designation == designation {
			return starPrecursor
		}
	}
	return nil
}

// createSiblingStar creates a sibling star - a companion that is similar to the parent
// but typically of a cooler spectral type. The subtype is rolled randomly.
func createSiblingStar(roller stellar.Roller, parentStar, companionStar *starPrecursor) error {
	if parentStar == nil {
		return fmt.Errorf("no parent provided")
	}
	if companionStar == nil {
		return fmt.Errorf("no companion provided")
	}
	switch parentStar.SubType {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		n, _ := strconv.Atoi(parentStar.SubType)
		if n == 0 {
			n = 10
		}
		n = n - roller.Roll("1d6")
		switch n > 0 {
		case true:
			companionStar.Type = parentStar.Type
			companionStar.SubType = strconv.Itoa(n)
		case false:
			companionStar.Type = getCoolerSpectralType(parentStar.Type)
			companionStar.SubType = parentStar.SubType
		}
	default:
		companionStar.Type = parentStar.Type
	}

	return nil
}

// createTwinStar creates a twin star - a companion with the same type, subtype, and class
// as the parent, with slightly less mass.
func createTwinStar(roller stellar.Roller, parentStar, companionStar *starPrecursor) error {
	if parentStar == nil {
		return fmt.Errorf("no parent provided")
	}
	if companionStar == nil {
		return fmt.Errorf("no companion provided")
	}
	companionStar.Class = parentStar.Class
	companionStar.Type = parentStar.Type
	companionStar.SubType = parentStar.SubType
	companionStar.Mass = float.RoundN(parentStar.Mass*(1-(0.01*float64(roller.Roll("1d6")))), 3)
	return nil
}

// createLesserCompanionStar creates a lesser companion star - a cooler and typically smaller
// companion. If the result would be an M-type with higher subtype than parent, converts to brown dwarf.
func createLesserCompanionStar(roller stellar.Roller, parentStar, companionStar *starPrecursor) error {
	if parentStar == nil {
		return fmt.Errorf("no parent provided")
	}
	if companionStar == nil {
		return fmt.Errorf("no companion provided")
	}
	companionStar.Class = parentStar.Class
	companionStar.Type = getCoolerSpectralType(parentStar.Type)
	n := roller.Roll("1d10-1")
	companionStar.SubType = strconv.Itoa(n)
	m, _ := strconv.Atoi(parentStar.SubType)
	if companionStar.Type == "M" && parentStar.Type == "M" && (m+1) < n {
		companionStar.SubType = ""
		companionStar.Type = "BD"
		companionStar.Class = ""
	}
	return nil
}

// createRandomCompanionStar creates a random companion star with a cooler spectral type.
// Similar to lesser but uses different roll for subtype.
// TODO: This function is defined but never used - implement or remove
func createRandomCompanionStar(roller stellar.Roller, parentStar, companionStar *starPrecursor) error {
	if parentStar == nil {
		return fmt.Errorf("no parent provided")
	}
	if companionStar == nil {
		return fmt.Errorf("no companion provided")
	}
	companionStar.Class = parentStar.Class
	companionStar.Type = getCoolerSpectralType(parentStar.Type)
	n := roller.Roll("1d10")
	companionStar.SubType = strconv.Itoa(n)
	m, _ := strconv.Atoi(parentStar.SubType)
	if companionStar.Type == "M" && parentStar.Type == "M" && (m+1) < n {
		companionStar.SubType = ""
		companionStar.Type = ""
		companionStar.Class = "BD"
	}
	return nil
}

// getCoolerSpectralType returns the next cooler spectral type in the sequence:
// O -> B -> A -> F -> G -> K -> M
// Used when creating companion stars that are cooler than their parent.
func getCoolerSpectralType(currentType string) string {
	switch currentType {
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
	return currentType
}

// orderedDesignations returns the designations of all stars in the system ordered
// by their canonical order (Primary, PrimaryComp, Close, etc.)
// TODO: This function is defined but never used - implement or remove
func orderedDesignations(stars map[orbit.Orbit]*starPrecursor) []stellar.StarDesignation {
	designations := []stellar.StarDesignation{}
	for _, designation := range stellar.AllDesignations() {
		for _, starPrecursor := range stars {
			if starPrecursor.Designation == designation {
				designations = append(designations, designation)
			}
		}
	}
	return designations
}

// starIterator provides iteration over stars in a star system, ordered by designation.
type starIterator struct {
	starPool     map[orbit.Orbit]*starPrecursor
	designations []stellar.StarDesignation
	index        int
	hasMore      bool
}

// newStarIterator creates a new star iterator for the given star map.
func newStarIterator(stars map[orbit.Orbit]*starPrecursor) *starIterator {
	iterator := starIterator{}
	iterator.starPool = stars
	iterator.designations = stellar.AllDesignations()
	iterator.index = -1
	return &iterator
}

// next advances the iterator to the next star.
// Returns true if another star exists, false if iteration is complete.
func (iterator *starIterator) next() bool {
	for i := iterator.index + 1; i < len(iterator.designations); i++ {
		for _, starPrecursor := range iterator.starPool {
			if starPrecursor.Designation == iterator.designations[i] {
				iterator.index = i
				iterator.hasMore = true
				return true
			}
		}
	}
	iterator.hasMore = false
	return false
}

// getValues returns the current orbit and star from the iterator.
// TODO: Rename to GetCurrent() for consistency with Go conventions
func (iterator *starIterator) getValues() (orbit.Orbit, *starPrecursor, error) {
	if !iterator.hasMore {
		return orbit.Orbit{}, nil, fmt.Errorf("no more stars")
	}
	for orbit, starPrecursor := range iterator.starPool {
		if starPrecursor.Designation == iterator.designations[iterator.index] {
			return orbit, starPrecursor, nil
		}
	}
	return orbit.Orbit{}, nil, fmt.Errorf("unexpected end")
}

func (iterator *starIterator) callPosition(p int) (orbit.Orbit, *starPrecursor) {
	if !iterator.hasMore {
		return orbit.Orbit{}, nil
	}
	if p < 0 || p > 7 {
		return orbit.Orbit{}, nil
	}
	wantDesignation := stellar.DesignationOrder[p]
	for orbit, starPrecursor := range iterator.starPool {
		if starPrecursor.Designation == wantDesignation {
			return orbit, starPrecursor
		}
	}
	return orbit.Orbit{}, nil

}

func (iterator *starIterator) restart() {
	iterator.hasMore = true
	iterator.index = -1
}

// PrintStarPositions prints a debug view of all stars and their orbits in the system.
// TODO: Remove or convert to String() method for proper debugging
func PrintStarPositions(systemPrecursor *starSystemPrecursor) {
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		orbit, starPrecursor, err := starIterator.getValues()
		if err != nil {
			panic(err) // TODO: Return error instead of panic
		}
		fmt.Printf("orbit %v around %v (%v) star: '%v%v %v'\n", starPrecursor.Designation, orbit.Designation, orbit.Distance, starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
	}
}
