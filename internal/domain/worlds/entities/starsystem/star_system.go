package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/pkg/dice"
)

// starSystemPrecursor represents an intermediate state during star system generation.
// It contains all data needed to build a complete star system including:
// - Primary star information
// - Secondary stars and their orbital parameters
// - System worlds (gas giants, belts, planets)
// - System metadata (age, coordinates, special flags)
//
// TODO: Consider renaming to StarSystemBuilder or StarSystemData to avoid "precursor" suffix
// TODO: Split into separate structs: StarSystem, Star, Orbit, SystemWorlds for cleaner separation
type starSystemPrecursor struct {
	ID           string
	coordinates  coordinates.Cube
	importedData t5ss.WorldData
	Empty        bool
	Dead         bool
	Primordial   bool
	Clustered    bool
	NebulaType   int
	Anomaly      bool
	PrimaryStar  *starPrecursor
	Age          float64
	Stars        map[orbit.Orbit]*starPrecursor
	GG           int
	Belts        int
	Planets      int
	TotalWorlds  int
}

// starPrecursor represents a star during the generation process.
// TSC (Type/Subtype/Class) represents the spectral classification:
//   - Type: Spectral class (O, B, A, F, G, K, M) or special types (D, BD, BH, NS, PSR)
//   - SubType: Numeric subtype (0-9 for main sequence, empty for special types)
//   - Class: Luminosity class (V=main sequence, IV=subgiant, III=giant, etc.)
//
// TODO: Add units to physical properties (Mass in solar masses, Diameter in solar diameters, etc.)
// TODO: Consider using typed constants for star types and classes instead of strings
type starPrecursor struct {
	Type         string
	SubType      string
	Class        string
	Designation  stellar.StarDesignation
	Dead         bool
	Protostar    bool
	Mass         float64
	Temperature  float64
	Diameter     float64
	Luminosity   float64
	Age          float64
	Anomaly      string
	Period       float64
	OrbitalSlots *orbitalSlots
}

// rollNebula rolls to determine the type of nebula present in the system.
// Returns 1 for emission nebula, 2 for reflection nebula (or other classification).
// Called when a star system is generated within a nebula.
func rollNebula(roller *dice.Roller) int {
	switch roller.Roll("1d6") {
	case 1:
		return 1
	default:
		return 2
	}
}

// Profile returns a human-readable string representation of the star system.
// Format: "numStars-TypeSubType Class-Mass-Diameter-Luminosity-Age:Designation-Orbit-Eccentricity-TypeSubType..."
// This is useful for debugging and comparing generated star systems.
// TODO: Add unit labels to the output (e.g., "solar masses", "AU")
// TODO: Consider returning a structured type instead of string for programmatic access
func (systemPrecursor *starSystemPrecursor) Profile() string {
	profile := ""
	starIterator := newStarIterator(systemPrecursor.Stars)
	starPosition := 0
	for starIterator.next() {
		starPosition++
		orbit, starPrecursor, err := starIterator.getValues()
		if err != nil {
			panic(err) // TODO: Return error instead of panicking
		}
		switch starPosition {
		case 1:
			profile += fmt.Sprintf("%v-", len(systemPrecursor.Stars))
			profile += fmt.Sprintf("%v%v", starPrecursor.Type, starPrecursor.SubType)
			if starPrecursor.Class != "" {
				profile += fmt.Sprintf(" %v", starPrecursor.Class)
			}
			profile += fmt.Sprintf("-%v", float.RoundN(starPrecursor.Mass, 3))
			profile += fmt.Sprintf("-%v", float.RoundN(starPrecursor.Diameter, 3))
			profile += fmt.Sprintf("-%v", float.RoundN(starPrecursor.Luminosity, 3))
			profile += fmt.Sprintf("-%v", float.RoundN(systemPrecursor.Age, 3))
		default:
			profile += fmt.Sprintf(":%v", starPrecursor.Designation)
			profile += fmt.Sprintf("-%v", orbit.Distance)
			profile += fmt.Sprintf("-%v", orbit.Eccentricity)
			profile += fmt.Sprintf("-%v%v", starPrecursor.Type, starPrecursor.SubType)
			if starPrecursor.Class != "" {
				profile += fmt.Sprintf(" %v", starPrecursor.Class)
			}
			profile += fmt.Sprintf("-%v", float.RoundN(starPrecursor.Mass, 3))
			profile += fmt.Sprintf("-%v", float.RoundN(starPrecursor.Diameter, 3))
			profile += fmt.Sprintf("-%v", starPrecursor.Luminosity)
		}

	}

	return profile
}
