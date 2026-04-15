// Package ehex provides encoding for extended hexadecimal (0-9, A-Z excluding I and O).
// It maps integers 0-33 to characters and vice versa.
package ehex

var (
	Unknown     = Ehex{code: "?", value: -101, description: "unknown"}
	Any         = Ehex{code: "*", value: -102, description: "any"}
	Invalid     = Ehex{code: "!", value: -103, description: "invalid"}
	Default     = Ehex{code: "#", value: -104, description: "default"}
	Ignore      = Ehex{code: "-", value: -105, description: "ignore"}
	Reserved    = Ehex{code: "&", value: -106, description: "reserved"}
	Masked      = Ehex{code: "~", value: -107, description: "masked"}
	Extension   = Ehex{code: ">", value: -108, description: "extension"}
	Placeholder = Ehex{code: ".", value: -109, description: "placeholder"}
)

// Ehex is minimal profile and data transfer encoding block.
type Ehex struct {
	code        string
	value       int
	description string
}

// New creates a custom Ehex with the given value, code, and optional description.
// This is intended for application‑specific extended codes that are not part of the
// standard 0–33 mapping. For standard values, use FromValue or FromCode.
func New(value int, code string, description ...string) Ehex {
	e := Ehex{}
	e.value = value
	e.code = code
	for _, d := range description {
		e.description = d
		break
	}
	return e
}

// valueToCode maps integer values (0-33) to their corresponding characters.
// Characters: 0-9, A-H, J-N, P-Z (excluding I and O).
var valueToCode = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9",
	10: "A", 11: "B", 12: "C", 13: "D", 14: "E", 15: "F", 16: "G", 17: "H", 18: "J", 19: "K",
	20: "L", 21: "M", 22: "N", 23: "P", 24: "Q", 25: "R", 26: "S", 27: "T", 28: "U", 29: "V",
	30: "W", 31: "X", 32: "Y", 33: "Z",
}

// codeToValue maps characters back to integer values (basic set).
var codeToValue = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	"A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15, "G": 16, "H": 17, "J": 18, "K": 19,
	"L": 20, "M": 21, "N": 22, "P": 23, "Q": 24, "R": 25, "S": 26, "T": 27, "U": 28, "V": 29,
	"W": 30, "X": 31, "Y": 32, "Z": 33,
}

// codeToValueExtended maps additional characters (lowercase, symbols) to values.
// These are checked only after the basic set.
var codeToValueExtended = map[string]int{
	"?": -101, "*": -102, "!": -103, "#": -104, "-": -105, "&": -106, "~": -107, ">": -108, ".": -109,
	"s": 0, "r": 0,
}

// FromValue returns an Ehex instance for the given integer value.
// If value is out of range (0-33), returns Unknown.
// Special values (Any, Invalid, etc.) cannot be obtained via FromValue.
func FromValue(value int) Ehex {
	if code, ok := valueToCode[value]; ok {
		return New(value, code)
	}
	return Unknown
}

// FromCode returns an Ehex instance for the given code character.
// It first checks the basic codeToValue map, then the extended map.
// If the code is not found in either, returns Unknown.
func FromCode(code string) Ehex {
	if value, ok := codeToValue[code]; ok {
		return New(value, code)
	}
	if value, ok := codeToValueExtended[code]; ok {
		switch value {
		case -101:
			return Unknown
		case -102:
			return Any
		case -103:
			return Invalid
		case -104:
			return Default
		case -105:
			return Ignore
		case -106:
			return Reserved
		case -107:
			return Masked
		case -108:
			return Extension
		case -109:
			return Placeholder
		default:
			return New(value, code)
		}
	}
	return Unknown
}

// WithDescription returns a copy of e with the description set.
func (e Ehex) WithDescription(desc string) Ehex {
	e.description = desc
	return e
}

// Code returns the code character of e (e.g., "A", "Z", "9").
// May return "?" or other string if not properly initialized.
func (e Ehex) Code() string {
	return e.code
}

// Value returns the integer value of e (0-33, or negative for special values).
func (e Ehex) Value() int {
	return e.value
}

// Description returns the human-readable description of e.
func (e Ehex) Description() string {
	return e.description
}

// String implements fmt.Stringer interface.
// Returns the code character (same as Code()).
func (e Ehex) String() string {
	return e.Code()
}
