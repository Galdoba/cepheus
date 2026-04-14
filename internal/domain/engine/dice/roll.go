package dice

func basicRoll(roller Roller, dp Dicepool) *Result {
	res := Result{}
	raw := make(map[*Die]int)
	for _, d := range dp.Dice {
		raw[&d] = roller.Roll(d)
	}
	res.Raw = raw
	return &res

}
