package tttable

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// RangeKey represents a parsed table key with boundaries
type RangeKey struct {
	Min          int    // Minimum value
	Max          int    // Maximum value
	MinInclusive bool   // Whether min boundary is inclusive
	MaxInclusive bool   // Whether max boundary is inclusive
	Original     string // Original key string
}

// MustIndex generates a table key string from given integers
// Panics on error. Use IndexSafe for error handling.
//
// Rules:
//   - If no integers provided: panic
//   - If one integer: return that integer as string (e.g., 5 -> "5")
//   - If two or more integers:
//   - Find min and max values
//   - If min == max: panic
//   - If min < MinRollBound and max > MaxRollBound: panic
//   - If min > MaxRollBound or max < MinRollBound: panic
//   - If min < MinRollBound: return "max-" (values <= max)
//   - If max > MaxRollBound: return "min+" (values >= min)
//   - Otherwise: return "min-max" range
func MustIndex(numbers ...int) string {
	key, err := IndexSafe(numbers...)
	if err != nil {
		panic(fmt.Sprintf("MustIndex error: %v", err))
	}
	return key
}

// IndexSafe generates a table key string from given integers
// Returns key and error based on validation rules.
func IndexSafe(numbers ...int) (string, error) {
	// Check if no arguments provided
	if len(numbers) == 0 {
		return "", ErrNoArguments
	}

	// Single number case
	if len(numbers) == 1 {
		return strconv.Itoa(numbers[0]), nil
	}

	// Multiple numbers case
	min, max := findMinMax(numbers)

	// Check if min equals max (shouldn't happen with multiple distinct numbers)
	if min == max {
		return "", ErrMinMaxEqual
	}

	// Check if entire range is out of bounds
	if min > MaxRollBound {
		return "", fmt.Errorf("%w: minimum %d exceeds upper bound %d",
			ErrOutOfBounds, min, MaxRollBound)
	}
	if max < MinRollBound {
		return "", fmt.Errorf("%w: maximum %d is below lower bound %d",
			ErrOutOfBounds, max, MinRollBound)
	}

	// Check if range exceeds both bounds
	if min < MinRollBound && max > MaxRollBound {
		return "", fmt.Errorf("%w: range %d-%d exceeds bounds [%d, %d]",
			ErrBothBoundsExceeded, min, max, MinRollBound, MaxRollBound)
	}

	// Handle lower bound violation
	if min < MinRollBound {
		// We can only represent values up to max
		return fmt.Sprintf("%d-", max), nil
	}

	// Handle upper bound violation
	if max > MaxRollBound {
		// We can only represent values from min upward
		return fmt.Sprintf("%d+", min), nil
	}

	// Normal range within bounds
	return fmt.Sprintf("%d-%d", min, max), nil
}

// findMinMax finds minimum and maximum values in a slice of integers
func findMinMax(numbers []int) (min, max int) {
	if len(numbers) == 0 {
		return 0, 0
	}

	min = numbers[0]
	max = numbers[0]

	for _, num := range numbers[1:] {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}
	return min, max
}

// ParseKey parses a table key string into a RangeKey
// Supported formats:
//
//	"5"      - single number 5
//	"5-10"   - range 5 to 10 inclusive
//	"5-"     - 5 and below (down to MinRollBound)
//	"5+"     - 5 and above (up to MaxRollBound)
//	"-5"     - single negative number -5
//	"-5-"    - -5 and below
//	"-5+"    - -5 and above
//	"-5--2"  - range from -5 to -2
//	"-5-2"   - range from -5 to 2
func ParseKey(key string) (*RangeKey, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("empty key")
	}

	// Try parsing as a single number first
	if num, err := strconv.Atoi(key); err == nil {
		return &RangeKey{
			Min:          num,
			Max:          num,
			MinInclusive: true,
			MaxInclusive: true,
			Original:     key,
		}, nil
	}

	// Parse using regex patterns
	return parseKeyWithRegex(key)
}

// parseKeyWithRegex parses complex keys using regular expressions
func parseKeyWithRegex(key string) (*RangeKey, error) {
	// Pattern for range with two numbers (can be negative)
	rangeRegex := regexp.MustCompile(`^(-?\d+)\s*-\s*(-?\d+)$`)
	// Pattern for lower bound only (X- or -X-)
	lowerBoundRegex := regexp.MustCompile(`^(-?\d+)-$`)
	// Pattern for upper bound only (X+ or -X+)
	upperBoundRegex := regexp.MustCompile(`^(-?\d+)\+$`)

	// Try range pattern first
	if matches := rangeRegex.FindStringSubmatch(key); matches != nil {
		min, err1 := strconv.Atoi(matches[1])
		max, err2 := strconv.Atoi(matches[2])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid range numbers in key: %s", key)
		}
		if min > max {
			return nil, fmt.Errorf("range minimum %d cannot be greater than maximum %d", min, max)
		}
		return &RangeKey{
			Min:          min,
			Max:          max,
			MinInclusive: true,
			MaxInclusive: true,
			Original:     key,
		}, nil
	}

	// Try lower bound pattern (X- or -X-)
	if matches := lowerBoundRegex.FindStringSubmatch(key); matches != nil {
		bound, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid bound in key: %s", key)
		}
		return &RangeKey{
			Min:          MinRollBound,
			Max:          bound,
			MinInclusive: true,
			MaxInclusive: true,
			Original:     key,
		}, nil
	}

	// Try upper bound pattern (X+ or -X+)
	if matches := upperBoundRegex.FindStringSubmatch(key); matches != nil {
		bound, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid bound in key: %s", key)
		}
		return &RangeKey{
			Min:          bound,
			Max:          MaxRollBound,
			MinInclusive: true,
			MaxInclusive: true,
			Original:     key,
		}, nil
	}

	return nil, fmt.Errorf("invalid key format: %s. Expected formats: 'X', 'X-Y', 'X-', 'X+', '-X', '-X-Y', '-X-', '-X+'", key)
}

// rangesOverlap checks if two RangeKey objects overlap
func rangesOverlap(r1, r2 *RangeKey) bool {
	// If one range ends before the other begins
	if r1.Max < r2.Min || r2.Max < r1.Min {
		return false
	}

	// If boundaries touch but are exclusive
	if r1.Max == r2.Min && !(r1.MaxInclusive && r2.MinInclusive) {
		return false
	}
	if r1.Min == r2.Max && !(r1.MinInclusive && r2.MaxInclusive) {
		return false
	}

	return true
}
