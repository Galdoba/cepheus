package dice

func (r *Roller) rollDice(d Die) int {
	return r.rng.Intn(d.Faces) + 1 // Generate random value between 1 and faces
}

func (r *Roller) rollPool(dp Dicepool) Result {
	res := Result{}
	res.Rolled = dp
	for _, d := range dp.Dice {
		res.Raw = append(res.Raw, r.rollDice(d))
	}
	mid := res.Raw
	for _, m := range dp.Modifiers {
		mid = m.Apply(mid)
	}
	res.Final = mid
	return res
}
