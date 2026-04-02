package starsystem

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/services/interpolate"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
	"github.com/Galdoba/cepheus/pkg/dice"
)

// runStep1 executes Step 1 of star system generation: determining the primary star's properties.
// This includes:
// 1. Creating or importing primary star data
// 2. Rolling for star Type, Subtype, and Class (TSC)
// 3. Calculating physical properties (mass, diameter, temperature, age)
// 4. Calculating luminosity
//
// The function sets the system age based on the primary star's age.
func (builder *Builder) runStep01(systemPrecursor *starSystemPrecursor) error {
	systemPrecursor.PrimaryStar = &starPrecursor{Designation: stellar.Primary}

	// TODO: Refactor inverted switch logic to use normal if-else pattern
	switch err := builder.importPrimaryStarData(systemPrecursor); err == nil {
	case true:
	default:
		if err := builder.determinePrimaryStarTypeAndClass(systemPrecursor); err != nil {
			return fmt.Errorf("failed to determine primary star type and class: %v", err)
		}
		err := builder.determineStarTSC(systemPrecursor.PrimaryStar, true)
		if err != nil {
			return err
		}
	}
	if err := validateStarTypeSubtypeClassCombination(systemPrecursor.PrimaryStar); err != nil {
		return err
	}

	if err := calculateStarPhysicalProperties(builder.rng, systemPrecursor.PrimaryStar); err != nil {
		return err
	}

	systemPrecursor.PrimaryStar.Luminosity = float.Round(calculateLuminosity(systemPrecursor.PrimaryStar.Diameter, systemPrecursor.PrimaryStar.Temperature))

	systemPrecursor.Age = systemPrecursor.PrimaryStar.Age
	systemPrecursor.PrimaryStar.Designation = stellar.Primary

	return nil
}

// determinePrimaryStarTypeAndClass rolls for the primary star's spectral type and luminosity class.
// It handles special cases like brown dwarfs, white dwarfs, neutron stars, black holes,
// protostars, and star clusters. The function uses a loop with labeled break to ensure
// a valid TSC combination is generated.
//
// Returns error if the random table rolling fails.
// TODO: The activeMods1 variable is declared but never used - either implement modifier system or remove
func (builder *Builder) determinePrimaryStarTypeAndClass(systemPrecursor *starSystemPrecursor) error {
	activeMods1 := []string{} // TODO: Remove unused variable
primary_star_class_generation:
	for {
		rolledResult, err := builder.step1.tablesStarType.Roll("Type")
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		switch rolledResult {
		case "O", "B", "A", "F", "G", "K", "M":
			switch systemPrecursor.PrimaryStar.Class {
			case "":
				systemPrecursor.PrimaryStar.Class = "V"
			case "IV":
				switch rolledResult {
				case "O":
					rolledResult = "B"
				case "M":
					continue
				}
			case "VI":
				if rolledResult == "F" {
					rolledResult = "G"
				}
				if rolledResult == "A" {
					rolledResult = "B"
				}
			}
			if systemPrecursor.Primordial {
				switch rolledResult {
				case "O", "B":
					continue
				}
			}
			systemPrecursor.PrimaryStar.Type = rolledResult
		case "Ia", "Ib", "II", "III", "IV", "VI":
			if rolledResult != "IV" && rolledResult != "VI" {
				activeMods1 = append(activeMods1, rtg.MOD_NonMainSequenceClass)
			}
			systemPrecursor.PrimaryStar.Class = rolledResult
		case "BD":
			systemPrecursor.PrimaryStar.Type = rolledResult
			systemPrecursor.PrimaryStar.Class = ""
			systemPrecursor.Empty = true
			break primary_star_class_generation
		case "D", "NS", "BH":
			systemPrecursor.PrimaryStar.Type = rolledResult
			systemPrecursor.PrimaryStar.Dead = true
			systemPrecursor.Empty = true
			break primary_star_class_generation
		case "Nb":
			if systemPrecursor.NebulaType == 0 {
				systemPrecursor.NebulaType = rollNebula(builder.rng)
			}
		case "Star Cluster":
			systemPrecursor.Clustered = true
		case "Protostar":
			systemPrecursor.PrimaryStar.Protostar = true
			activeMods1 = append(activeMods1, rtg.MOD_ProtostarSystem)
		case "PSR":
			systemPrecursor.PrimaryStar.Type = rolledResult
			systemPrecursor.Dead = true
			break primary_star_class_generation
		case "Anomaly":
			systemPrecursor.PrimaryStar.Type = rolledResult
			break primary_star_class_generation
		default:
			// TODO: Return error instead of panic for invalid rolled values
			panic(fmt.Sprintf("dev error: invalid value rolled: %v", rolledResult))
		}
		if systemPrecursor.PrimaryStar.Class != "" && systemPrecursor.PrimaryStar.Type != "" {
			break primary_star_class_generation
		}
	}
	return nil
}

// determineStarTSC rolls for a star's Type, Subtype, and Class (TSC).
// This is used for both primary and secondary stars. The function handles
// special types like brown dwarfs, white dwarfs, neutron stars, etc.
//
// The primary flag indicates whether this is the primary star in the system.
// TODO: Refactor to reduce code duplication with determinePrimaryStarTypeAndClass
func (builder *Builder) determineStarTSC(starPrecursor *starPrecursor, primary bool, mods ...string) error {
	activeMods1 := []string{} // TODO: Remove unused variable
primary_star_class_generation:
	for {
		rolledResult, err := builder.step1.tablesStarType.Roll("Type", mods...)
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		switch rolledResult {
		case "O", "B", "A", "F", "G", "K", "M":
			switch starPrecursor.Class {
			case "":
				starPrecursor.Class = "V"
			case "IV":
				switch rolledResult {
				case "O":
					rolledResult = "B"
				case "M":
					continue
				}
			case "VI":
				if rolledResult == "F" {
					rolledResult = "G"
				}
				if rolledResult == "A" {
					rolledResult = "B"
				}
			}
			starPrecursor.Type = rolledResult
			if err := builder.determineStarSubtype(starPrecursor); err != nil {
				return fmt.Errorf("subtype error: %v", err)
			}
		case "Ia", "Ib", "II", "III", "IV", "VI":
			if rolledResult != "IV" && rolledResult != "VI" {
				activeMods1 = append(activeMods1, rtg.MOD_NonMainSequenceClass)
			}
			starPrecursor.Class = rolledResult
		case "BD":
			starPrecursor.Type = rolledResult
			starPrecursor.Class = ""
			starPrecursor.SubType = ""
			break primary_star_class_generation
		case "D", "NS", "BH", "Anomaly", "PSR":
			starPrecursor.Type = rolledResult
			starPrecursor.Class = ""
			starPrecursor.SubType = ""
			break primary_star_class_generation
		case "Nb":
		case "Star Cluster":
		case "Protostar":
			activeMods1 = append(activeMods1, rtg.MOD_ProtostarSystem)
		default:
			// TODO: Return error instead of panic for invalid rolled values
			panic(fmt.Sprintf("dev error: invalid value rolled: %v", rolledResult))
		}
		switch primary {
		case true:
			if starPrecursor.Type != "" && starPrecursor.Class != "" {
				break primary_star_class_generation
			}
		case false:
			if validateStarTypeSubtypeClassCombination(starPrecursor) == nil {
				break primary_star_class_generation
			}
		}

	}
	return nil
}

// determineStarSubtype rolls for the numeric subtype of a star.
// For M-type stars, the subtype is a letter (a, b, c, etc.)
// For O, B, A, F, G, K stars, the subtype is a number (0-9)
// Special rules apply for K-type Class IV stars (subtract 5 from results > 4)
func (builder *Builder) determineStarSubtype(starPrecursor *starPrecursor) error {
	switch starPrecursor.Type {
	case "M":
		rolledResult, err := builder.step1.tablesStarType.Roll("M Type Primary")
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		starPrecursor.SubType = rolledResult
	case "O", "B", "A", "F", "G", "K":
		rolledResult, err := builder.step1.tablesStarType.Roll("M Type Primary")
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		n, err := strconv.Atoi(rolledResult)
		if err != nil {
			return fmt.Errorf("expect number for subtype: '%v'", rolledResult)
		}
		// Apply special rule for K-type Class IV stars
		if starPrecursor.Class == "IV" && starPrecursor.Type == "K" && n > 4 {
			n = n - 5
		}
		starPrecursor.SubType = fmt.Sprintf("%v", n)
	default:
		return nil
	}
	return nil
}

// calculateStarPhysicalProperties calculates a star's mass, diameter, temperature, and age
// based on its TSC (Type/Subtype/Class) using interpolation tables for main sequence stars
// and special formulas for degenerate objects (white dwarfs, neutron stars, black holes).
func calculateStarPhysicalProperties(roller *dice.Roller, starPrecursor *starPrecursor) error {
	interpolationIndex := interpolate.Index(starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
	switch starPrecursor.Type {
	case "D", "BD":
		starPrecursor.Mass = calculateWhiteDwarfMass(roller)
		starPrecursor.Diameter = float.Round(calculateWhiteDwarfDiameter(starPrecursor.Mass))
		starPrecursor.Age = calculateStarAge(roller, starPrecursor)
	case "BH":
		starPrecursor.Mass = calculateBlackHoleMass(roller)
		starPrecursor.Diameter = float.Round(2.95 * starPrecursor.Mass) // TODO: Add unit comment (km)
		starPrecursor.Age = calculateStarAge(roller, starPrecursor)
	case "NS", "PSR":
		starPrecursor.Mass = calculateNeutronStarMass(roller)
		starPrecursor.Diameter = 19 + float64(roller.Roll("1d6")) // TODO: Add unit comment (km)
		starPrecursor.Age = calculateStarAge(roller, starPrecursor)
	case "Anomaly":
		// TODO: Handle anomaly type - currently does nothing

	default:
		starPrecursor.Mass = float.Round(interpolate.MassByIndex(interpolationIndex))
		if starPrecursor.Mass == 0 {
			fmt.Println(interpolationIndex)
			return fmt.Errorf("failed to determine by interpolation: star mass (%v%v %v)", starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
		}
		starPrecursor.Diameter = float.Round(interpolate.DiamByIndex(interpolationIndex))
		if starPrecursor.Diameter == 0 {
			return fmt.Errorf("failed to determine by interpolation: star diameter (%v%v %v)", starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
		}
		starPrecursor.Age = calculateStarAge(roller, starPrecursor)
	}
	starPrecursor.Temperature = float.Round(interpolate.TempByIndex(interpolationIndex))
	return nil
}

// calculateWhiteDwarfMass rolls for the mass of a white dwarf star.
// Uses 2d6 and 1d10 rolls to determine mass in solar masses.
// If mass exceeds 1.44 (Chandrasekhar limit), additional mass is added.
func calculateWhiteDwarfMass(roller *dice.Roller) float64 {
	roll1 := float64(roller.Roll("2d6"))
	roll2 := float64(roller.Roll("1d10"))
	mass := float.Round(((roll1 - 1) / 10) + (roll2 / 100))
	if mass > 1.44 {
		mass = 1.34 + float64(roller.Roll("1d100"))/1000
	}
	return mass
}

// calculateWhiteDwarfDiameter calculates the diameter of a white dwarf based on its mass.
// Uses the inverse relationship between mass and diameter for degenerate matter.
func calculateWhiteDwarfDiameter(mass float64) float64 {
	return (1.0 / mass) * 0.01 // TODO: Add units and formula explanation in comment
}

// calculateBlackHoleMass rolls for the mass of a black hole.
// Uses 1d6 and 1d10 rolls with reroll on 6 for additional mass.
func calculateBlackHoleMass(roller *dice.Roller) float64 {
	roll6 := roller.Roll("1d6")
	roll10 := roller.Roll("1d10")
	mass := 2.1 + float64(roll6) - 1 + (float64(roll10) / 10)
	for roll6 == 6 {
		roll6 = roller.Roll("1d6")
		mass += float64(roll6)
	}
	return float.Round(mass)
}

// calculateNeutronStarMass rolls for the mass of a neutron star.
// Uses 1d6 with reroll on 6 for additional mass.
func calculateNeutronStarMass(roller *dice.Roller) float64 {
	roll1 := roller.Roll("1d6")
	mass := 1 + (float64(roll1) / 10)
	for roll1 == 6 {
		roll1 = roller.Roll("1d6")
		mass += (float64(roll1) - 1.0) / 10.0
	}
	return mass
}

// calculateVariance rolls for a variance multiplier used in age calculations.
// Returns a value between 0 and 1.
func calculateVariance(roller *dice.Roller) float64 {
	roll1 := roller.Roll("1d1001-1")
	return float64(roll1) / 1000.0
}

// calculateMainSequenceLifespan calculates the expected lifespan of a main sequence star
// based on its mass. More massive stars have shorter lifespans.
func calculateMainSequenceLifespan(mass float64) float64 {
	return float.Round(10 / math.Pow(mass, 2.5))
}

// calculateSmallStarAge calculates age for small stars (less than 0.9 solar masses)
// using a different formula than main sequence lifespan.
func calculateSmallStarAge(roller *dice.Roller) float64 {
	return float64(roller.Roll("1d6x2")-roller.Roll("1d3-2")) + float64(roller.Roll("1d10")/10.0)
}

// calculateSubgiantLifespan calculates the expected lifespan of a subgiant star
// based on its main sequence lifespan and mass.
func calculateSubgiantLifespan(mainSequenceLifespan, mass float64) float64 {
	return mainSequenceLifespan / (4 + mass)
}

// calculateGiantLifespan calculates the expected lifespan of a giant star
// based on its main sequence lifespan and mass.
func calculateGiantLifespan(mainSequenceLifespan, mass float64) float64 {
	return mainSequenceLifespan / 10.0 * math.Pow(mass, 3)
}

// calculateDeadStarMass calculates the mass of a dead star (remnant).
// Used for age calculations of white dwarfs, neutron stars, and black holes.
func calculateDeadStarMass(roller *dice.Roller, mass float64) float64 {
	return float.Round(float64(roller.Roll("1d3+2")) * mass)
}

// calculateStarAge calculates the age of a star based on its mass, class, and type.
// Different formulas are used for:
// - Main sequence stars
// - Subgiants (Class IV)
// - Giants (Class III)
// - White dwarfs, neutron stars, black holes
// - Protostars
func calculateStarAge(roller *dice.Roller, starPrecursor *starPrecursor) float64 {
	age := 0.0
	mass := starPrecursor.Mass
	if starPrecursor.Dead {
		mass = calculateDeadStarMass(roller, mass)
	}
	mainSequenceLifespan := calculateMainSequenceLifespan(mass)
	if mass < 0.9 {
		mainSequenceLifespan = calculateSmallStarAge(roller)
	}
	age = mainSequenceLifespan * calculateVariance(roller)
	switch starPrecursor.Class {
	case "BD":
		age = calculateSmallStarAge(roller)
	case "D", "NS", "BH":
		massValue := calculateDeadStarMass(roller, starPrecursor.Mass)
		mainSequenceLifespanValue := calculateMainSequenceLifespan(calculateDeadStarMass(roller, massValue))
		age = calculateSmallStarAge(roller) + mainSequenceLifespanValue + calculateSubgiantLifespan(mainSequenceLifespanValue, massValue) + calculateGiantLifespan(mainSequenceLifespanValue, massValue)
	case "PSR":
		massValue := calculateDeadStarMass(roller, starPrecursor.Mass)
		mainSequenceLifespanValue := calculateMainSequenceLifespan(calculateDeadStarMass(roller, massValue))
		age = (0.1 * float64(roller.Roll("2d10"))) + mainSequenceLifespanValue + calculateSubgiantLifespan(mainSequenceLifespanValue, massValue) + calculateGiantLifespan(mainSequenceLifespanValue, massValue)
	case "IV":
		age = mainSequenceLifespan + (calculateSubgiantLifespan(mainSequenceLifespan, mass) * calculateVariance(roller))
	case "III":
		age = mainSequenceLifespan + calculateSubgiantLifespan(mainSequenceLifespan, mass) + (calculateGiantLifespan(mainSequenceLifespan, mass) * calculateVariance(roller))

	}
	if starPrecursor.Protostar {
		age = float64(roller.Roll("2d10")) * 0.01
	}
	age = min(13.8, float.Round(age))
	for age <= 0 {
		age = calculateVariance(roller) * 13.8
	}
	return age
}

// calculateLuminosity calculates a star's luminosity based on its diameter and temperature.
// Uses the Stefan-Boltzmann law: L = 4 * pi * R^2 * sigma * T^4
// Simplified to relative luminosity assuming solar values.
func calculateLuminosity(diameter, temperature float64) float64 {
	return math.Pow(diameter, 2) * math.Pow(temperature/float64(5772), 4)
}

// validateStarTypeSubtypeClassCombination validates that the star's TSC (Type/Subtype/Class) combination is valid.
// TSC represents the spectral classification of a star:
//   - Type: Spectral class (O, B, A, F, G, K, M) or special types (D, BD, BH, NS, PSR)
//   - SubType: Numeric subtype (0-9 for main sequence, empty for special types)
//   - Class: Luminosity class (V=main sequence, IV=subgiant, III=giant, VI=subdwarf, etc.)
//
// The function also corrects some invalid combinations (e.g., converts to valid types)
// and returns an error for truly invalid combinations.
func validateStarTypeSubtypeClassCombination(starPrecursor *starPrecursor) error {
	switch starPrecursor.Type {
	case "":
		return fmt.Errorf("no star?")
	case "O", "B", "A", "F", "G", "K", "M":
		if starPrecursor.Class == "" {
			starPrecursor.Class = "V"
		}
	case "D", "BD", "BH", "NS", "PSR", "NB":
		starPrecursor.SubType = ""
		starPrecursor.Class = ""
	}
	switch starPrecursor.Class {
	case "Ia", "Ib", "II", "III", "V":
		if !strings.Contains("OBAFGKM", starPrecursor.Type) {
			return fmt.Errorf("invalid combination type=%v subtype=%v class=%v", starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
		}
	case "IV":
		if !strings.Contains("BAFGK", starPrecursor.Type) {
			starPrecursor.Class = "V"
		}
		if starPrecursor.Type == "K" && !strings.Contains("01234", starPrecursor.SubType) {
			starPrecursor.Class = "V"
		}
	case "VI":
		if starPrecursor.Type == "F" {
			starPrecursor.Type = "G"
		}
		if starPrecursor.Type == "A" {
			starPrecursor.Type = "B"
		}
		if !strings.Contains("OBGKM", starPrecursor.Type) {
			return fmt.Errorf("invalid combination type=%v subtype=%v class=%v", starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
		}
		if starPrecursor.Type == "F" && !strings.Contains("56789", starPrecursor.SubType) {
			return fmt.Errorf("invalid combination type=%v subtype=%v class=%v", starPrecursor.Type, starPrecursor.SubType, starPrecursor.Class)
		}
	}
	return nil
}

// importPrimaryStarData attempts to import primary star data from the imported t5ss data.
// If no stellar data is available in the import, returns an error.
// Otherwise, decodes the stellar code to extract Type, SubType, and Class.
func (builder *Builder) importPrimaryStarData(systemPrecursor *starSystemPrecursor) error {
	if builder.imported.Stellar == "" {
		return fmt.Errorf("nothing to import")
	}
	stellarObject, err := stellar.New(builder.imported.Stellar)
	if err != nil {
		return fmt.Errorf("failed to create stellar: %w", err)
	}
	primaryCode := stellarObject.PrimaryCode()

	starType, starSubtype, starClass, err := primaryCode.Decode()
	if err != nil {
		return fmt.Errorf("failed to decode primary code (%v): %w", primaryCode, err)
	}
	systemPrecursor.PrimaryStar.Type = starType
	systemPrecursor.PrimaryStar.SubType = starSubtype
	systemPrecursor.PrimaryStar.Class = starClass
	return nil
}
