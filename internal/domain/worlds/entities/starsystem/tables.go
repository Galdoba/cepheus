package starsystem

type TableRoller interface {
	DeterminePrimaryStar() *starPrecursor
}

func primaryStarTables() (TableRoller, error) {
	return nil, nil
}
