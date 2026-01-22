package dice

import (
	"hash/crc64"
	"time"
)

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

func WithSeedString(s string) SetupOption {
	return func(ss *setupSettings) {
		ss.seed = stringToInt64(s)
	}
}

func stringToInt64(s string) int64 {
	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.Checksum([]byte(s), table)
	return int64(hash)
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
