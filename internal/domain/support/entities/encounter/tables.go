package encounter

import "github.com/Galdoba/cepheus/pkg/tttable"

func EncounterDistance() *tttable.Table {
	tt, err := tttable.NewTable("Encounter Distance",
		tttable.WithDiceExpression("2d6"),
		tttable.WithIndexEntries(
			tttable.NewTableEntry("2-", "Close"),
			tttable.NewTableEntry("3", "Short"),
			tttable.NewTableEntry("4-5", "Medium"),
			tttable.NewTableEntry("6-9", "Long"),
			tttable.NewTableEntry("10-11", "Very Long"),
			tttable.NewTableEntry("12+", "Distant"),
		),
		tttable.WithIndexMods(tttable.Flat, map[string]int{
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

func Person() *tttable.D66Table {
	tt, err := tttable.NewD66Table("Person",
		tttable.WithD66Entries(
			tttable.NewTableEntry("11", "Naval Officer"),
			tttable.NewTableEntry("12", "Imperial Diplomat"),
			tttable.NewTableEntry("13", "Crooked Trader"),
			tttable.NewTableEntry("14", "Medical Doctor"),
			tttable.NewTableEntry("15", "Eccentric Scientist"),
			tttable.NewTableEntry("16", "Mercenary"),
			tttable.NewTableEntry("21", "Famous Performer"),
			tttable.NewTableEntry("22", "Alien Thief"),
			tttable.NewTableEntry("23", "Free Trader"),
			tttable.NewTableEntry("24", "Explorer"),
			tttable.NewTableEntry("25", "Marine Captain"),
			tttable.NewTableEntry("26", "Corporate Executive"),
			tttable.NewTableEntry("31", "Researcher"),
			tttable.NewTableEntry("32", "Cultural Attaché"),
			tttable.NewTableEntry("33", "Religious Leader"),
			tttable.NewTableEntry("34", "Conspirator"),
			tttable.NewTableEntry("35", "Rich Noble"),
			tttable.NewTableEntry("36", "Artificial Intelligence"),
			tttable.NewTableEntry("41", "Bored Noble"),
			tttable.NewTableEntry("42", "Planetary Governor"),
			tttable.NewTableEntry("43", "Inveterate Gambler"),
			tttable.NewTableEntry("44", "Crusading Journalist"),
			tttable.NewTableEntry("45", "Doomsday Cultist"),
			tttable.NewTableEntry("46", "Corporate Agent"),
			tttable.NewTableEntry("51", "Criminal Syndicate"),
			tttable.NewTableEntry("52", "Military Governor "),
			tttable.NewTableEntry("53", "Army Quartermaster"),
			tttable.NewTableEntry("54", "Private Investigator"),
			tttable.NewTableEntry("55", "Starport Administrator"),
			tttable.NewTableEntry("56", "Retired Admiral"),
			tttable.NewTableEntry("61", "Alien Ambassador"),
			tttable.NewTableEntry("62", "Smuggler"),
			tttable.NewTableEntry("63", "Weapons Inspector "),
			tttable.NewTableEntry("64", "Elder Statesman"),
			tttable.NewTableEntry("65", "Planetary Warlord"),
			tttable.NewTableEntry("66", "Imperial Agent"),
		),
	)
	if err != nil {
		panic(err)
	}
	return tt
}

func CharacterQuircks() *tttable.D66Table {
	tt, err := tttable.NewD66Table("Character Quircks",
		tttable.WithD66Entries(
			tttable.NewTableEntry("11", "Loyal"),
			tttable.NewTableEntry("12", "Distracted by other worries"),
			tttable.NewTableEntry("13", "In debt to criminals"),
			tttable.NewTableEntry("14", "Makes very bad jokes"),
			tttable.NewTableEntry("15", "Will betray characters"),
			tttable.NewTableEntry("16", "Aggressive"),
			tttable.NewTableEntry("21", "Has secret allies"),
			tttable.NewTableEntry("22", "Secret anagathic user"),
			tttable.NewTableEntry("23", "Looking for something"),
			tttable.NewTableEntry("24", "Helpful"),
			tttable.NewTableEntry("25", "Forgetful"),
			tttable.NewTableEntry("26", "Wants to hire the Travellers"),
			tttable.NewTableEntry("31", "Has useful contacts"),
			tttable.NewTableEntry("32", "Artistic"),
			tttable.NewTableEntry("33", "Easily confused"),
			tttable.NewTableEntry("34", "Unusually ugly"),
			tttable.NewTableEntry("35", "Worried about current situation"),
			tttable.NewTableEntry("36", "Shows pictures of their children"),
			tttable.NewTableEntry("41", "Rumour-monger"),
			tttable.NewTableEntry("42", "Unusually provincial"),
			tttable.NewTableEntry("43", "Drunkard or drug addict"),
			tttable.NewTableEntry("44", "Government informant"),
			tttable.NewTableEntry("45", "Mistakes a Traveller for someone else"),
			tttable.NewTableEntry("46", "Possesses unusually advanced technology"),
			tttable.NewTableEntry("51", "Unusually handsome or beautiful"),
			tttable.NewTableEntry("52", "Spying on the Travellers"),
			tttable.NewTableEntry("53", "Possesses TAS membership"),
			tttable.NewTableEntry("54", "Is secretly hostile towards the Travellers"),
			tttable.NewTableEntry("55", "Wants to borrow money"),
			tttable.NewTableEntry("56", "Is convinced the Travellers are dangerous"),
			tttable.NewTableEntry("61", "Involved in political intrigue"),
			tttable.NewTableEntry("62", "Has a dangerous secret"),
			tttable.NewTableEntry("63", "Wants to get off planet as soon as possible"),
			tttable.NewTableEntry("64", "Attracted to a Traveller"),
			tttable.NewTableEntry("65", "From offworld"),
			tttable.NewTableEntry("66", "Possesses telepathy or other unusual quality"),
		),
	)
	if err != nil {
		panic(err)
	}
	return tt
}
