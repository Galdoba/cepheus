package world

import "github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"

func (w *World) SetMainworldUWP(u uwp.UWP) {
	w.mainworldUWP = u
}
