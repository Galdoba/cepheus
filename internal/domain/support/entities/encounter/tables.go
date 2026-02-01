package encounter

import "github.com/Galdoba/cepheus/pkg/tttable"

func EncounterDistance() *tttable.Table {
	tt, err := tttable.NewTable("Encounter Distance",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow("2-", "Close"),
			tttable.NewRow("3", "Short"),
			tttable.NewRow("4-5", "Medium"),
			tttable.NewRow("6-9", "Long"),
			tttable.NewRow("10-11", "Very Long"),
			tttable.NewRow("12+", "Distant"),
		),
		tttable.WithMods(map[string]int{
			"Clean Terrain":                          +3,
			"Forest or Woods":                        -2,
			"Crowded Area":                           -2,
			"In Space":                               +4,
			"Target is a Vechicle":                   +2,
			"Tarvellers actively looking for danger": +1,
		},
		),
	)
	if err != nil {
		panic(err)
	}
	return tt
}

func Person() *tttable.Table {
	tt, err := tttable.NewTable("Encounter Distance",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow("2-", "Close"),
			tttable.NewRow("3", "Short"),
			tttable.NewRow("4-5", "Medium"),
			tttable.NewRow("6-9", "Long"),
			tttable.NewRow("10-11", "Very Long"),
			tttable.NewRow("12+", "Distant"),
		),
		tttable.WithMods(map[string]int{
			"Clean Terrain":                          +3,
			"Forest or Woods":                        -2,
			"Crowded Area":                           -2,
			"In Space":                               +4,
			"Target is a Vechicle":                   +2,
			"Tarvellers actively looking for danger": +1,
		},
		),
	)
	if err != nil {
		panic(err)
	}
	return tt
}
