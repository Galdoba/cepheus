package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compiled regular expressions for parsing sum dice expressions.
// Compiling once improves performance for repeated parsing.
var (
	baseRe           = regexp.MustCompile(`^(\d+)d(\d+)`)          // Matches base pattern like 2d6, 10d10
	additiveRe       = regexp.MustCompile(`[\+-]\d+`)              // Matches additive modifiers like +2, -3
	multiplicativeRe = regexp.MustCompile(`x[\+-]?\d+`)            // Matches multiplicative modifiers like x2, x-1
	deletiveRe       = regexp.MustCompile(`/[\+-]?\d+`)            // Matches deletive modifiers like /2, /-3
	replaceRe        = regexp.MustCompile(`r(\d+(?:;\d+)*):(\d+)`) // Matches replace patterns like r1:6, r1;2:10
	dropLowRe        = regexp.MustCompile(`dl(\d+)`)               // Matches drop low patterns like dl1, dl2
	dropHighRe       = regexp.MustCompile(`dh(\d+)`)               // Matches drop high patterns like dh1, dh2
	individualRe     = regexp.MustCompile(`i([\+-]?\d+)`)          // Matches individual modifiers like i+1, i-2
	sumMinRe         = regexp.MustCompile(`min([\+-]?\d+)`)        // Matches minimum sum patterns like min5, min-2
	sumMaxRe         = regexp.MustCompile(`max([\+-]?\d+)`)        // Matches maximum sum patterns like max10, max-1
	rerollRe         = regexp.MustCompile(`rr(\d+(?:;\d+)*)`)      // Matches reroll patterns like rr1, rr1;2
)

// parseSumString is a single-pass parser for dice sum expressions.
// It parses the entire expression and returns SumDirectives with all modifiers applied.
// The parsing order is important to avoid conflicts between similar patterns.
func parseSumString(s string) (SumDirectives, error) {
	sd := newSumDirectives()
	remaining := strings.ToLower(strings.TrimSpace(s))

	// Step 1: Parse the base part (required) - NdM format
	baseMatch := baseRe.FindStringSubmatch(remaining)
	if baseMatch == nil {
		return sd, fmt.Errorf("failed to parse base Sum")
	}

	// Extract and parse the number of dice
	num, err := strconv.Atoi(baseMatch[1])
	if err != nil {
		return sd, fmt.Errorf("failed to parse number of dice: %v", err)
	}
	sd.Num = num

	// Extract and parse the number of faces
	faces, err := strconv.Atoi(baseMatch[2])
	if err != nil {
		return sd, fmt.Errorf("failed to parse dice faces: %v", err)
	}
	sd.Faces = faces

	// Remove the parsed base part from the remaining string
	remaining = strings.TrimPrefix(remaining, baseMatch[0])

	// Helper function to remove found matches from the remaining string
	removeFound := func(match string) string {
		if idx := strings.Index(remaining, match); idx != -1 {
			remaining = remaining[:idx] + remaining[idx+len(match):]
		}
		return remaining
	}

	// Step 2: Parse reroll values (rr) - Processed before replacements to avoid conflicts
	for {
		match := rerollRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}
		// Split semicolon-separated values
		valuesStr := match[1]
		values := strings.Split(valuesStr, ";")

		for _, valStr := range values {
			value, err := strconv.Atoi(valStr)
			if err != nil {
				return sd, fmt.Errorf("failed to parse reroll value '%s': %v", valStr, err)
			}
			// Check for duplicate reroll values
			if _, exists := sd.ReRoll[value]; exists {
				return sd, fmt.Errorf("duplicated reroll value %v", value)
			}
			sd.ReRoll[value] = true
		}

		remaining = removeFound(match[0])
	}

	// Step 3: Parse replacements (can be multiple)
	for {
		match := replaceRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}
		// Split semicolon-separated source values
		sourcesStr := match[1]
		sources := strings.Split(sourcesStr, ";")

		// Parse the replacement value
		replacement, err := strconv.Atoi(match[2])
		if err != nil {
			return sd, fmt.Errorf("failed to parse replacement value: %v", err)
		}

		// Map each source value to the replacement
		for _, srcStr := range sources {
			original, err := strconv.Atoi(srcStr)
			if err != nil {
				return sd, fmt.Errorf("failed to parse replacement original: %v", err)
			}
			// Check for duplicate replacements
			if _, exists := sd.Replace[original]; exists {
				return sd, fmt.Errorf("double assignment to replace value %v", original)
			}
			sd.Replace[original] = replacement
		}

		remaining = removeFound(match[0])
	}

	// Step 4: Parse drop low (only one allowed)
	if match := dropLowRe.FindStringSubmatch(remaining); match != nil {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return sd, fmt.Errorf("failed to parse drop low value: %v", err)
		}
		sd.SumMods[DropLow] = value
		remaining = removeFound(match[0])
	}

	// Step 5: Parse drop high (only one allowed)
	if match := dropHighRe.FindStringSubmatch(remaining); match != nil {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return sd, fmt.Errorf("failed to parse drop high value: %v", err)
		}
		sd.SumMods[DropHigh] = value
		remaining = removeFound(match[0])
	}

	// Step 6: Parse individual modifiers (can be multiple, summed together)
	individualSum := 0
	for {
		match := individualRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}
		// Parse the value (can be with or without sign)
		valueStr := match[1]
		// If no sign is present, assume positive
		if valueStr[0] != '+' && valueStr[0] != '-' {
			valueStr = "+" + valueStr
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return sd, fmt.Errorf("failed to parse individual mod: %v", err)
		}
		individualSum += value
		remaining = removeFound(match[0])
	}
	sd.SumMods[Individual] = individualSum

	// Step 7: Parse minimum sum (only one allowed)
	if match := sumMinRe.FindStringSubmatch(remaining); match != nil {
		valueStr := match[1]
		// If no sign is present, assume positive
		if valueStr[0] != '+' && valueStr[0] != '-' {
			valueStr = "+" + valueStr
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return sd, fmt.Errorf("failed to parse minimum sum: %v", err)
		}
		sd.SumMods[SumMininum] = value
		remaining = removeFound(match[0])
	}

	// Step 8: Parse maximum sum (only one allowed)
	if match := sumMaxRe.FindStringSubmatch(remaining); match != nil {
		valueStr := match[1]
		// If no sign is present, assume positive
		if valueStr[0] != '+' && valueStr[0] != '-' {
			valueStr = "+" + valueStr
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return sd, fmt.Errorf("failed to parse maximum sum: %v", err)
		}
		sd.SumMods[SumMaximum] = value
		remaining = removeFound(match[0])
	}

	// Step 9: Parse multiplicative modifiers (can be multiple, multiplied together)
	multiplicativeProd := 1
	for {
		match := multiplicativeRe.FindString(remaining)
		if match == "" {
			break
		}
		value, err := strconv.Atoi(strings.TrimPrefix(match, "x"))
		if err != nil {
			return sd, fmt.Errorf("failed to parse multiplicative mod: %v", err)
		}
		multiplicativeProd *= value
		remaining = removeFound(match)
	}
	sd.SumMods[Multiplicative] = multiplicativeProd

	// Step 10: Parse deletive modifiers (can be multiple, multiplied together)
	deletiveProd := 1
	for {
		match := deletiveRe.FindString(remaining)
		if match == "" {
			break
		}
		value, err := strconv.Atoi(strings.TrimPrefix(match, "/"))
		if err != nil {
			return sd, fmt.Errorf("failed to parse deletive mod: %v", err)
		}
		deletiveProd *= value
		remaining = removeFound(match)
	}
	sd.SumMods[Deletive] = deletiveProd

	// Step 11: Parse additive modifiers (processed last to avoid conflicts with other patterns)
	additiveSum := 0
	for {
		match := additiveRe.FindString(remaining)
		if match == "" {
			break
		}
		value, err := strconv.Atoi(match)
		if err != nil {
			return sd, fmt.Errorf("failed to parse additive mod: %v", err)
		}
		additiveSum += value
		remaining = removeFound(match)
	}
	sd.SumMods[Additive] = additiveSum

	// Step 12: Verify that the entire string was processed
	if strings.TrimSpace(remaining) != "" {
		return sd, fmt.Errorf("unrecognized tokens in expression: %s", remaining)
	}

	// Step 13: Validate the parsed values
	if sd.Num <= 0 {
		return sd, fmt.Errorf("number of dice must be positive, got %d", sd.Num)
	}
	if sd.Faces <= 0 {
		return sd, fmt.Errorf("number of faces must be positive, got %d", sd.Faces)
	}
	if dl, ok := sd.SumMods[DropLow]; ok && dl > 0 {
		if dl >= sd.Num {
			return sd, fmt.Errorf("drop low count (%d) must be less than number of dice (%d)", dl, sd.Num)
		}
	}
	if dh, ok := sd.SumMods[DropHigh]; ok && dh > 0 {
		if dh >= sd.Num {
			return sd, fmt.Errorf("drop high count (%d) must be less than number of dice (%d)", dh, sd.Num)
		}
	}

	// Step 14: Clean up default values (remove identity operations)
	if sd.SumMods[Multiplicative] == 1 {
		delete(sd.SumMods, Multiplicative)
	}
	if sd.SumMods[Deletive] == 1 {
		delete(sd.SumMods, Deletive)
	}
	if sd.SumMods[Individual] == 0 {
		delete(sd.SumMods, Individual)
	}
	if sd.SumMods[Additive] == 0 {
		delete(sd.SumMods, Additive)
	}

	return sd, nil
}
