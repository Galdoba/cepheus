package dice

import "time"

type setupSettings struct {
	seed   int64
	dices  []dice
	locked bool
}

func defaultSettings() setupSettings {
	return setupSettings{
		seed: time.Now().UnixNano(),
	}
}

type SetupOption func(*setupSettings)

func WithSeed(seed int64) SetupOption {
	return func(ss *setupSettings) {
		ss.seed = seed
	}
}

func Locked() SetupOption {
	return func(ss *setupSettings) {
		ss.locked = true
	}
}

func WithDices(dices ...dice) SetupOption {
	return func(ss *setupSettings) {
		for _, d := range dices {
			ss.dices = append(ss.dices, d)
		}
	}
}
