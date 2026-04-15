package dice

// die represents a single physical die with a number of faces and optional metadata.
type die struct {
	faces    int
	codes    map[int]string
	metadata map[string]string
}

func newDie(faces int) die {
	return die{faces: faces}
}

func (d die) withCodes(codes map[int]string) die {
	d.codes = codes
	return d
}

func (d die) withMeta(meta map[string]string) die {
	d.metadata = meta
	return d
}

// dicepool is a collection of dice that are rolled together.
type dicepool struct {
	dice     []die
	metadata map[string]string
}

func newDicepool(dice ...die) dicepool {
	return dicepool{
		dice:     dice,
		metadata: map[string]string{},
	}
}

func (dp dicepool) withMeta(meta map[string]string) dicepool {
	dp.metadata = meta
	return dp
}