package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions for parsing concatenated dice expressions
var (
	concatBaseRe = regexp.MustCompile(`^d(\d+)`)            // Matches base pattern like d66, d505
	concatModRe  = regexp.MustCompile(`cm(\d+):([+-]?\d+)`) // Matches modifier patterns like cm1:5
)

// parseConcatString parses a concatenated dice expression string into ConcatDirectives.
// The expression format is dN where N is a sequence of digits, each representing a die.
// Optional modifiers can be added with cmA:X format where A is the die position (1-based)
// and X is the modifier value to apply to that die.
func parseConcatString(s string) (ConcatDirectives, error) {
	cd := ConcatDirectives{}
	remaining := strings.ToLower(strings.TrimSpace(s))

	// Step 1: Parse the base part (required) - dN format
	baseMatch := concatBaseRe.FindStringSubmatch(remaining)
	if baseMatch == nil {
		return cd, fmt.Errorf("failed to parse base Concat: expected format dN")
	}

	// Extract and remove the base part from the string
	diceStr := baseMatch[1]

	// Split the number into individual digits - each digit represents a die's faces
	for _, char := range diceStr {
		face, err := strconv.Atoi(string(char))
		if err != nil {
			// This shouldn't happen since the regex guarantees digits
			return cd, fmt.Errorf("invalid digit in dice specification: %v", char)
		}
		cd.Faces = append(cd.Faces, face)
	}

	// Initialize modifiers array with zeros (no modification by default)
	cd.Mods = make([]int, len(cd.Faces))

	// Remove the parsed base part from the remaining string
	remaining = strings.TrimPrefix(remaining, baseMatch[0])

	// Helper function to remove found matches from the remaining string
	removeFound := func(match string) string {
		if idx := strings.Index(remaining, match); idx != -1 {
			remaining = remaining[:idx] + remaining[idx+len(match):]
		}
		return remaining
	}

	// Step 2: Parse modifiers (cmA:X format)
	modMap := make(map[int]int) // Track modifiers to detect duplicates

	for {
		match := concatModRe.FindStringSubmatch(remaining)
		if match == nil {
			break // No more modifiers found
		}

		// Parse the die index (A) - 1-based position
		diceIndex, err := strconv.Atoi(match[1])
		if err != nil {
			return cd, fmt.Errorf("failed to parse dice index in modifier: %v", match[1])
		}

		// Validate that the index is within range
		if diceIndex < 1 || diceIndex > len(cd.Faces) {
			return cd, fmt.Errorf("dice index %d out of range (1-%d)", diceIndex, len(cd.Faces))
		}

		// Check for duplicate modifiers for the same die
		if _, exists := modMap[diceIndex]; exists {
			return cd, fmt.Errorf("duplicate modifier for dice %d", diceIndex)
		}

		// Parse the modifier value (X)
		modValue, err := strconv.Atoi(match[2])
		if err != nil {
			return cd, fmt.Errorf("failed to parse modifier value: %v", match[2])
		}

		// Store the modifier (convert to 0-based index)
		cd.Mods[diceIndex-1] = modValue
		modMap[diceIndex] = modValue

		// Remove this modifier from the remaining string
		remaining = removeFound(match[0])
	}

	// Step 3: Verify that the entire string was processed
	if strings.TrimSpace(remaining) != "" {
		return cd, fmt.Errorf("unrecognized tokens in expression: %s", remaining)
	}

	return cd, nil
}
