package stellar

import (
	"fmt"
	"strings"
)

// StarDesignation represents a designation for a star in a multiple star system,
// such as primary, close companion, etc.
type StarDesignation string

// Predefined star designations.
const (
	Primary     StarDesignation = "Aa"
	Close       StarDesignation = "Ba"
	Near        StarDesignation = "Ca"
	Far         StarDesignation = "Da"
	PrimaryComp StarDesignation = "Ab"
	CloseComp   StarDesignation = "Bb"
	NearComp    StarDesignation = "Cb"
	FarComp     StarDesignation = "Db"
)

// Stellar represents a description of a star system, consisting of one or more
// star codes concatenated together.
type Stellar string

// StarCode represents a spectral classification code for a star, e.g., "G2 V".
type StarCode string

// Roller provides a dice rolling function for random generation.
type Roller interface {
	Roll(string) int
}

// ------------------------------------------------------------------------
// Exported constants for spectral classification
// ------------------------------------------------------------------------

// Spectral type constants.
const (
	SpectralTypeO = "O"
	SpectralTypeB = "B"
	SpectralTypeA = "A"
	SpectralTypeF = "F"
	SpectralTypeG = "G"
	SpectralTypeK = "K"
	SpectralTypeM = "M"
	SpectralTypeL = "L"
	SpectralTypeT = "T"
	SpectralTypeY = "Y"
)

// SpectralTypes contains all recognized spectral types.
var SpectralTypes = []string{
	SpectralTypeO, SpectralTypeB, SpectralTypeA, SpectralTypeF,
	SpectralTypeG, SpectralTypeK, SpectralTypeM,
	SpectralTypeL, SpectralTypeT, SpectralTypeY,
}

// MainSequenceTypes contains spectral types O through M (main sequence stars).
var MainSequenceTypes = []string{
	SpectralTypeO, SpectralTypeB, SpectralTypeA, SpectralTypeF,
	SpectralTypeG, SpectralTypeK, SpectralTypeM,
}

// BrownDwarfTypes contains spectral types L, T, Y (brown dwarfs).
var BrownDwarfTypes = []string{
	SpectralTypeL, SpectralTypeT, SpectralTypeY,
}

// Luminosity class constants.
const (
	LuminosityClassIa  = "Ia"
	LuminosityClassIb  = "Ib"
	LuminosityClassII  = "II"
	LuminosityClassIII = "III"
	LuminosityClassIV  = "IV"
	LuminosityClassV   = "V"
	LuminosityClassVI  = "VI"
)

// LuminosityClasses contains all recognized luminosity classes.
var LuminosityClasses = []string{
	LuminosityClassIa, LuminosityClassIb,
	LuminosityClassIII, LuminosityClassII,
	LuminosityClassIV, LuminosityClassVI, LuminosityClassV,
}

// Subtype constants (numerals 0-9).
const (
	Subtype0 = "0"
	Subtype1 = "1"
	Subtype2 = "2"
	Subtype3 = "3"
	Subtype4 = "4"
	Subtype5 = "5"
	Subtype6 = "6"
	Subtype7 = "7"
	Subtype8 = "8"
	Subtype9 = "9"
)

// StarSubtypes contains all possible numeric subtypes.
var StarSubtypes = []string{
	Subtype0, Subtype1, Subtype2, Subtype3, Subtype4,
	Subtype5, Subtype6, Subtype7, Subtype8, Subtype9,
}

// Special star code constants.
const (
	SpecialStarBD  = "BD"
	SpecialStarD   = "D"
	SpecialStarBH  = "BH"
	SpecialStarNS  = "NS"
	SpecialStarPNb = "pNb"
	SpecialStarPSR = "PSR"
)

// SpecialStarCodes contains all special star codes (non‑standard spectral types).
var SpecialStarCodes = []string{
	SpecialStarBD, SpecialStarD, SpecialStarBH,
	SpecialStarNS, SpecialStarPNb, SpecialStarPSR,
}

// DesignationOrder defines the order in which designations are considered when
// generating a multiple star system.
var DesignationOrder = []StarDesignation{Primary, Close, Near, Far}

// ------------------------------------------------------------------------
// Exported functions
// ------------------------------------------------------------------------

// Decode parses the StarCode into its spectral type, subtype, and luminosity class.
// Returns an error if the code is invalid.
func (sc StarCode) Decode() (string, string, string, error) {
	t, s, c := "", "", ""
	str := string(sc)
	parts := strings.Split(str, " ")
	switch len(parts) {
	case 1:
	case 2:
		for _, class := range LuminosityClasses {
			if parts[1] == class {
				c = class
			}
		}
		if c == "" {
			return "", "", "", fmt.Errorf("invalid star code: %v", sc)
		}
	default:
		return "", "", "", fmt.Errorf("invalid star code: %v", sc)
	}

	s = extractSubtype(parts[0])
	if s != "" {
		ok := false
		for _, subtype := range StarSubtypes {
			if subtype == s {
				ok = true
				break
			}
		}
		if !ok {
			return "", "", "", fmt.Errorf("invalid star code: %v", sc)
		}
	}

	t = strings.TrimSuffix(parts[0], s)
	return t, s, c, nil
}

// String returns the string representation of the StarCode.
func (sc StarCode) String() string {
	return string(sc)
}

// New creates a Stellar from the given string after validating it.
// Returns an error if the string contains invalid star codes.
func New(s string) (Stellar, error) {
	if err := Stellar(s).Validate(); err != nil {
		return "", err
	}
	return Stellar(s), nil
}

// Validate checks that all star codes in the Stellar string are valid.
func (s Stellar) Validate() error {
	for _, str := range ExtractStars(string(s)) {
		if !validateStar(str) {
			return fmt.Errorf("invalid code '%v' in '%v'", str, s)
		}
	}
	return nil
}

// String returns a normalized representation of the stellar system with
// star codes joined by colons.
func (s Stellar) String() string {
	str := ""
	for _, star := range ExtractStars(string(s)) {
		str += string(star) + ":"
	}
	return strings.TrimSuffix(str, ":")
}

// Split returns the individual star codes as a slice of strings.
func (s Stellar) Split() []string {
	return strings.Split(s.String(), ":")
}

// PrimaryCode returns the first star code in the system, which is typically the primary star.
func (s Stellar) PrimaryCode() StarCode {
	codes := ExtractStars(string(s))
	if len(codes) < 1 {
		return ""
	}
	return codes[0]
}

// ExtractStars parses a stellar system string and returns a slice of StarCode
// for each recognized star.
func ExtractStars(stellar string) []StarCode {
	scodes := []StarCode{}
	for _, star := range parseStellar(stellar) {
		scodes = append(scodes, StarCode(star))
	}
	return scodes
}

// RollStellarDesignations assigns random designations (like Aa, Ba, etc.) to each star
// in the given stellar system. The number of designations must match the number of stars;
// it calls RollDesignations repeatedly until a match is found.
func RollStellarDesignations(r Roller, s Stellar) []StarDesignation {
	stars := ExtractStars(string(s))
	listing := []StarDesignation{}
	for len(listing) != len(stars) {
		listing = RollDesignations(r)
	}
	return listing
}

// RollDesignations generates a random set of star designations for a multiple star system
// based on dice rolls. The result may include primary and companions depending on rolls.
func RollDesignations(r Roller) []StarDesignation {
	stars := make(map[StarDesignation]bool)
	stars[Primary] = true
	for _, des := range DesignationOrder {
		if des == Primary {
			stars[des] = true
			continue
		}
		if r.Roll("2d6") >= 10 {
			stars[des] = true
		}
	}
	result := []StarDesignation{}
	for _, des := range DesignationOrder {
		if stars[des] {
			result = append(result, des)
			code := strings.TrimSuffix(string(des), "a")
			if r.Roll("2d6") >= 10 {
				result = append(result, StarDesignation(code+"b"))
			}
		}
	}
	return result
}

// AllDesignations returns a slice of all possible star designations.
func AllDesignations() []StarDesignation {
	return []StarDesignation{Primary, PrimaryComp, Close, CloseComp, Near, NearComp, Far, FarComp}
}

// ------------------------------------------------------------------------
// Unexported helpers
// ------------------------------------------------------------------------

func validateStar(s StarCode) bool {
	switch s {
	case SpecialStarBD, SpecialStarD, SpecialStarBH, SpecialStarNS, SpecialStarPSR, SpecialStarPNb:
		return true
	}
	validType := false
	validClass := false
	asString := string(s)
	for _, tp := range SpectralTypes {
		if strings.HasPrefix(asString, tp) {
			validType = true
		}
	}
	for _, cl := range LuminosityClasses {
		if strings.HasSuffix(asString, cl) {
			validClass = true
		}
	}
	return validType && validClass
}

func listStars() []string {
	stars := make([]string, 0, len(SpecialStarCodes)+
		len(MainSequenceTypes)*len(StarSubtypes)*len(LuminosityClasses)+
		len(BrownDwarfTypes)*len(StarSubtypes))

	// Special codes
	stars = append(stars, SpecialStarCodes...)

	// Main sequence stars (O..M with subtype and class)
	for _, sType := range MainSequenceTypes {
		for _, sSubtype := range StarSubtypes {
			for _, sClass := range LuminosityClasses {
				stars = append(stars, fmt.Sprintf("%v%v %v", sType, sSubtype, sClass))
			}
		}
	}

	// Brown dwarfs (L,T,Y with subtype only)
	for _, sType := range BrownDwarfTypes {
		for _, sSubtype := range StarSubtypes {
			stars = append(stars, fmt.Sprintf("%v%v", sType, sSubtype))
		}
	}
	return stars
}

func parseStellar(stellar string) []string {
	stars := []string{}
main_loop:
	for stellar != "" {
		for _, s := range listStars() {
			if strings.HasPrefix(stellar, s) {
				stars = append(stars, s)
				stellar = strings.TrimPrefix(stellar, s)
				continue main_loop
			}
		}
		stellar = cutFirstLetter(stellar)
	}
	return stars
}

func cutFirstLetter(s string) string {
	newStr := ""
	for i, l := range strings.Split(s, "") {
		if i == 0 {
			continue
		}
		newStr += l
	}
	return newStr
}

func extractSubtype(s string) string {
	ss := ""
	for _, l := range strings.Split(s, "") {
		for _, reference := range StarSubtypes {
			if l == reference {
				ss += l
			}
		}
	}
	return ss
}
