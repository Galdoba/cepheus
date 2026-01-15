// Package ehex provides functionality for Extended Hexadecimal (Ehex) encoding.
// Ehex extends standard hexadecimal with additional characters to represent
// values from 0 to 33. It's commonly used in Traveller RPG and other systems.
package ehex

import "fmt"

// Ehex represents an Extended Hexadecimal value with associated metadata.
// It encapsulates a character code, integer value, and optional description.
type Ehex struct {
	code        string // Single character Ehex code (0-9, A-Z)
	value       int    // Integer value represented by the code (0-33)
	description string // Human-readable description of the value
}

// invalid creates and returns an Ehex instance representing an invalid value.
// This is used as a common error response for invalid inputs.
func invalid(s string) Ehex {
	return Ehex{"", -10, fmt.Sprintf("invalid ehex data provided: '%v'", s)}
}

// FromCode creates an Ehex instance from a single-character code.
// Return detected Ehex (special code possible) or an invalid Ehex if the code is not recognized.
func FromCode(code string) Ehex {
	if _, ok := codeToValue[code]; !ok {
		switch code {
		case "*":
			return Any()
		case "?":
			return Unknown()
		case "_":
			return Unassigned()
		case "-":
			return Separator1()
		}
		return invalid(code)
	}
	return Ehex{code, codeToValue[code], ""}
}

// FromValue creates an Ehex instance from an integer value.
// Returns an invalid Ehex if the value is not in the valid range (0-33).
func FromValue(val int) Ehex {
	if _, ok := valueToCode[val]; !ok {
		return invalid(fmt.Sprintf("%v", val))
	}
	return Ehex{valueToCode[val], val, ""}
}

// Any returns a special Ehex instance representing "any value".
// This is typically used as a wildcard in pattern matching.
func Any() Ehex {
	return Ehex{"*", -1, "any value"}
}

// Unknown returns a special Ehex instance representing an unknown value.
// This is typically used to indicate missing or undetermined data.
func Unknown() Ehex {
	return Ehex{"?", -2, "unknown value"}
}

// Unassigned returns a special Ehex instance representing an unassigned value.
// This is typically used as a placeholder for future assignment.
func Unassigned() Ehex {
	return Ehex{"_", -3, "unassigned value"}
}

// Separator1 returns a special Ehex instance representing an unassigned value.
// This is technical value is set to visualy separate Ehex blocks.
func Separator1() Ehex {
	return Ehex{"-", -4, "code separator 1"}
}

// Special creates a custom Ehex instance with arbitrary properties.
// This allows creation of Ehex values outside the standard 0-33 range,
// such as for application-specific special values.
func Special(code string, val int, description string) Ehex {
	return Ehex{code: code, value: val, description: description}
}

// Code returns the single-character Ehex code.
func (e Ehex) Code() string {
	return e.code
}

// Value returns the integer value of the Ehex.
// Special values (Any, Unknown, Unassigned) return negative numbers.
func (e Ehex) Value() int {
	return e.value
}

// Description returns the human-readable description of the Ehex value.
func (e Ehex) Description() string {
	return e.description
}

// SetDescription set description text to Ehex and return it's copy.
func SetDescription(e Ehex, description string) Ehex {
	e.description = description
	return e
}

// Less compares two Ehex values for ordering.
// Special values (with negative numbers) are treated as 0 for comparison.
// Returns true if e1 is less than e2.
func Less(e1, e2 Ehex) bool {
	return max(e1.value, 0) < max(e2.value, 0)
}

// Equal compares two Ehex values for equality.
// Returns true if both Ehex values have the same integer value.
// Returns false if any Ehex value is negative (invalid or uncertain).
func Equal(e1, e2 Ehex) bool {
	if e1.value < 0 || e2.value < 0 {
		return false
	}
	return e1.value == e2.value
}

// codeToValue maps single-character Ehex codes to their integer values.
// Includes standard hexadecimal 0-9,A-F and extended values up to Z.
var codeToValue = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3, "4": 4,
	"5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	"A": 10, "B": 11, "C": 12, "D": 13, "E": 14,
	"F": 15, "G": 16, "H": 17, "J": 18, "K": 19,
	"L": 20, "M": 21, "N": 22, "P": 23, "Q": 24,
	"R": 25, "S": 26, "T": 27, "U": 28, "V": 29,
	"W": 30, "X": 31, "Y": 32, "Z": 33,
}

// valueToCode maps integer values to their single-character Ehex codes.
// This is the inverse of codeToValue for quick value-to-code lookups.
var valueToCode = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3", 4: "4",
	5: "5", 6: "6", 7: "7", 8: "8", 9: "9",
	10: "A", 11: "B", 12: "C", 13: "D", 14: "E",
	15: "F", 16: "G", 17: "H", 18: "J", 19: "K",
	20: "L", 21: "M", 22: "N", 23: "P", 24: "Q",
	25: "R", 26: "S", 27: "T", 28: "U", 29: "V",
	30: "W", 31: "X", 32: "Y", 33: "Z",
}
