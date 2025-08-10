package lscconfig

import "encoding/json"

type Config struct {
	ShipCrew       int `json:"ship crew total"`
	PlayersPresent int `json:"players present"`
	OfficerRatio   int `json:"officers ratio per crew members"`
	MasteryMod     int `json:"mastery mod"`
}

func Default() Config {
	return Config{
		ShipCrew:       11,
		PlayersPresent: 5,
		OfficerRatio:   10,
		MasteryMod:     0,
	}
}

func (cfg *Config) Marshal() ([]byte, error) {
	return json.Marshal(cfg)
}

func (cfg *Config) Unmarshal(data []byte) error {
	return json.Unmarshal(data, cfg)
}
