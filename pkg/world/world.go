package world

import "github.com/Galdoba/cepheus/pkg/uwp"

type World struct {
	Name       string
	UWP        uwp.UWP
	TradeCodes []string
	TravelZone string
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
