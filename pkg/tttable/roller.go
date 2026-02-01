package tttable

// Roller interface for random number generation
type Roller interface {
	RollSafe(string) (int, error)
	ConcatRollSafe(string) (string, error)
}
