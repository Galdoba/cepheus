package world

import "github.com/Galdoba/cepheus/pkg/uwp"

type World struct {
	Name        string   `json:"name"`
	IsMainWorld *bool    `json:"is main world,omitempty"`
	UWP         uwp.UWP  `json:"uwp"`
	TradeCodes  []string `json:"classifications"`
	TravelZone  string   `json:"travel zone,omitempty"`
	OrbitN      *float64 `json:"orbit number,omitempty"`
	HZCO        *float64 `json:"habitable zone orbit difference,omitempty"`
}

func New(opts ...WorldOption) *World {
	w := World{}
	for _, modify := range opts {
		modify(&w)
	}

	return &w
}

type WorldOption func(*World)

func WithUWP(profile string) WorldOption {
	return func(w *World) {
		u := uwp.New(uwp.FromString(profile))
		w.UWP = u
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

func (w *World) OrbitNumber() (float64, bool) {
	if w.OrbitN == nil {
		return 0, false
	}
	return *w.OrbitN, true
}
