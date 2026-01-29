package orbital

import (
	"errors"
	"strconv"
	"strings"
)

type StarDesignation string
type SatelliteOrbit string

// Valid star designations in priority order
var validStars = []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Cb", "Da", "Db"}
var starToPair = map[string]string{
	"Aa": "A", "Ab": "A",
	"Ba": "B", "Bb": "B",
	"Ca": "C", "Cb": "C",
	"Da": "D", "Db": "D",
}

// Pair order for validation
var pairOrder = map[string]int{"A": 0, "B": 1, "C": 2, "D": 3}

// Valid satellite orbits
var validSatellites = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

var (
	ErrInvalidStar        = errors.New("invalid star designation")
	ErrDuplicateStar      = errors.New("duplicate star")
	ErrInvalidPairOrder   = errors.New("invalid star pair order")
	ErrMissingPrimary     = errors.New("primary star Aa is missing")
	ErrMissingMainStar    = errors.New("main star is missing for companion")
	ErrInvalidPlanetOrbit = errors.New("invalid planet orbit (allowed 0-20)")
	ErrInvalidSatellite   = errors.New("invalid satellite orbit (allowed a-l)")
	ErrSatelliteNoPlanet  = errors.New("satellite specified without planet")
	ErrInvalidCodeFormat  = errors.New("invalid code format")
	ErrInvalidPairCode    = errors.New("invalid star pair code")
	ErrPairNotSorted      = errors.New("star pairs are not sorted by distance")
	ErrEmptyCode          = errors.New("empty code")
)

// Encode converts orbital position into full string code (without compression)
func Encode(stars []string, planet int, satellite string) (string, error) {
	// Validate stars
	if err := validateStars(stars); err != nil {
		return "", err
	}

	// Sort stars in the correct order: Aa, Ab, Ba, Bb, Ca, Cb, Da, Db
	sortedStars := sortStars(stars)

	// Build star center code (full designations, no spaces)
	code := strings.Join(sortedStars, "")

	// Add planet if specified
	if planet >= 0 {
		if planet > 20 {
			return "", ErrInvalidPlanetOrbit
		}
		code += " " + strconv.Itoa(planet)

		// Add satellite if specified
		if satellite != "" {
			if !isValidSatellite(satellite) {
				return "", ErrInvalidSatellite
			}
			code += " " + satellite
		}
	} else {
		// Check that satellite is not specified without planet
		if satellite != "" {
			return "", ErrSatelliteNoPlanet
		}
	}

	return code, nil
}

// Decode parses full string code into components (expects full star designations)
func Decode(code string) (stars []string, planet int, satellite string, err error) {
	if code == "" {
		return nil, -1, "", ErrEmptyCode
	}

	// Split code into parts
	parts := strings.Fields(code)
	if len(parts) < 1 || len(parts) > 3 {
		return nil, -1, "", ErrInvalidCodeFormat
	}

	// Parse star part - should contain full star designations
	starCode := parts[0]

	// Check if star code length is multiple of 2 (each star is 2 chars)
	if len(starCode)%2 != 0 {
		return nil, -1, "", ErrInvalidCodeFormat
	}

	// Split into individual stars (2 chars each)
	stars = make([]string, 0)
	for i := 0; i < len(starCode); i += 2 {
		star := starCode[i : i+2]
		if !IsValidStar(star) {
			return nil, -1, "", ErrInvalidStar
		}
		stars = append(stars, star)
	}

	// Validate star order and composition
	if err := validateStars(stars); err != nil {
		return nil, -1, "", err
	}

	// Initialize default values
	planet = -1
	satellite = ""

	// Parse planet if present
	if len(parts) >= 2 {
		p, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, -1, "", ErrInvalidPlanetOrbit
		}
		if p < 0 || p > 20 {
			return nil, -1, "", ErrInvalidPlanetOrbit
		}
		planet = p

		// Parse satellite if present
		if len(parts) == 3 {
			sat := parts[2]
			if !isValidSatellite(sat) {
				return nil, -1, "", ErrInvalidSatellite
			}
			satellite = sat
		}
	}

	return stars, planet, satellite, nil
}

// Compress converts full star designations to compressed pair codes
func Compress(code string) (string, error) {
	// Decode to get full star list
	stars, planet, satellite, err := Decode(code)
	if err != nil {
		return "", err
	}

	// Collect unique pairs in correct order
	pairs := make([]string, 0)
	seenPairs := make(map[string]bool)

	for _, star := range stars {
		pair := starToPair[star]
		if !seenPairs[pair] {
			pairs = append(pairs, pair)
			seenPairs[pair] = true
		}
	}

	// Validate pair order
	for i := 1; i < len(pairs); i++ {
		if pairOrder[pairs[i]] < pairOrder[pairs[i-1]] {
			return "", ErrInvalidPairOrder
		}
	}

	// Build compressed code
	compressedCode := strings.Join(pairs, "")

	// Add planet if present
	if planet >= 0 {
		compressedCode += " " + strconv.Itoa(planet)

		// Add satellite if present
		if satellite != "" {
			compressedCode += " " + satellite
		}
	}

	return compressedCode, nil
}

// sortStars sorts stars in the correct order: Aa, Ab, Ba, Bb, Ca, Cb, Da, Db
func sortStars(stars []string) []string {
	// Create a map for quick lookup
	starMap := make(map[string]bool)
	for _, star := range stars {
		starMap[star] = true
	}

	// Build sorted list using the predefined order
	sorted := make([]string, 0)
	for _, validStar := range validStars {
		if starMap[validStar] {
			sorted = append(sorted, validStar)
		}
	}

	return sorted
}

// validateStars validates star list
func validateStars(stars []string) error {
	if len(stars) == 0 {
		return ErrInvalidStar
	}

	// Check for Aa presence
	hasAa := false

	// Validate uniqueness and validity
	seen := make(map[string]bool)
	mainStars := make(map[string]bool)  // Main stars (ending with 'a')
	companions := make(map[string]bool) // Companions (ending with 'b')

	for _, star := range stars {
		if star == "Aa" {
			hasAa = true
		}

		// Validate designation
		if !IsValidStar(star) {
			return ErrInvalidStar
		}

		// Check for duplicates
		if seen[star] {
			return ErrDuplicateStar
		}
		seen[star] = true

		// Classify stars
		if strings.HasSuffix(star, "a") {
			mainStars[star[:1]] = true
		} else {
			companions[star[:1]] = true
		}
	}
	if !hasAa {
		return ErrMissingPrimary
	}

	// Check that each companion has its main star
	for pair := range companions {
		if !mainStars[pair] {
			return ErrMissingMainStar
		}
	}

	// Validate order - stars should be in validStars order
	lastIndex := -1
	for _, star := range stars {
		// Find index in validStars
		idx := -1
		for i, validStar := range validStars {
			if star == validStar {
				idx = i
				break
			}
		}

		if idx <= lastIndex {
			return ErrInvalidPairOrder
		}
		lastIndex = idx
	}

	return nil
}

// isValidSatellite validates satellite orbit
func isValidSatellite(s string) bool {
	if len(s) != 1 {
		return false
	}

	for _, valid := range validSatellites {
		if s == valid {
			return true
		}
	}
	return false
}

// Helper functions for external use

// GetStarComponents returns components of a star pair
func GetStarComponents(pair string) []string {
	components := []string{}

	switch pair {
	case "A":
		components = []string{"Aa", "Ab"}
	case "B":
		components = []string{"Ba", "Bb"}
	case "C":
		components = []string{"Ca", "Cb"}
	case "D":
		components = []string{"Da", "Db"}
	}

	return components
}

// IsValidStar checks if string is a valid star designation
func IsValidStar(star string) bool {
	for _, validStar := range validStars {
		if star == validStar {
			return true
		}
	}
	return false
}
