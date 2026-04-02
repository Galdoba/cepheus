package starsystem

// TableRoller defines an interface for rolling on star system generation tables.
// This interface is currently not fully implemented and is a placeholder for
// future table-driven generation.
// TODO: Implement proper table rolling interface or remove unused code
type TableRoller interface {
	DeterminePrimaryStar() *starPrecursor
}

// primaryStarTables returns a table roller for primary star generation.
// TODO: Implement or remove - currently returns nil
func primaryStarTables() (TableRoller, error) {
	return nil, nil
}
