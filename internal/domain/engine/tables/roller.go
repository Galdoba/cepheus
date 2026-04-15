package tables

// TableRoller is dice roller interface.
// Implemented by /domain/engine/dice package. Can be inlemented elsewhere.
type TableRoller interface {
	// D66 - is a 2d6 roll that concatenate first and second die values with applied mods (if any)
	// expercted output is a string containing only digits: "00" - "99"
	D66(...int) string
	// Roll- is a standard roll that returns sum of dice rolled
	// Provided mods might be handled separatly.
	Roll(string, ...int) (int, error)
}
