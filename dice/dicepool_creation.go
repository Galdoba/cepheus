package dice

type Dicepool struct {
	roller *roller
	dices  map[int]dice
	locked bool
}

func NewDicepool(opts ...SetupOption) *Dicepool {
	settings := defaultSettings()
	dp := Dicepool{}
	for _, modify := range opts {
		modify(&settings)
	}
	dp.roller = newRoller(settings.seed)
	dp.dices = make(map[int]dice)
	for i, d := range settings.dices {
		dp.dices[i] = d
	}
	dp.locked = settings.locked
	return &dp
}
