package dice

type Dice struct {
	ID     string
	Edges  int
	Value  int
	Glyphs []string
}

func (d *Dice) Score() int {
	return d.Value
}

func (d *Dice) Which() string {
	return d.Which()
}

// dice.Roll("2d6")
// dice.D66() string
// dice.D3()
// dice.RollN(2)
