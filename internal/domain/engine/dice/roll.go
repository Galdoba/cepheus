package dice

// basicRoll rolls every die in the dicepool and returns a Result.
func basicRoll(r roller, dp dicepool) result {
	res := result{}
	for _, d := range dp.dice {
		res.dice = append(res.dice, d)
		res.raw = append(res.raw, r.roll(d))
	}
	return res
}
