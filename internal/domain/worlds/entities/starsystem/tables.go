package starsystem

type TableRoller interface {
	DeterminePrimaryStar() *Star
}

func primaryStarTables() (TableRoller, error) {
	return nil, nil
}
