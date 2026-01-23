package world

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

type World struct {
	id          string
	imported    t5ss.WorldData
	coordinates coordinates.Global
	trade       *tradeConnections
}

func Import(wd t5ss.WorldData) (*World, error) {
	w := World{}
	w.id = fmt.Sprintf("{%v,%v}", wd.WorldX, wd.WorldY)
	w.imported = wd
	w.coordinates = wd.Coordinates()
	w.trade = newTradeConnections(wd.TradeCodes()...)
	return &w, nil
}

func (w *World) Coordinates() coordinates.Cube {
	return w.coordinates.ToCube()
}

func (w *World) UWP() uwp.UWP {
	u, _ := uwp.New(w.imported.UWP)
	return u
}

func (w *World) Hydrate() error {
	// fill world with custom data
	return nil
}

func (w *World) SearchKey() string {
	return w.imported.SearchKey()
}

func (w *World) DatabaseKey() string {
	crd := w.coordinates
	return fmt.Sprintf("{%v,%v}", crd.X(), crd.Y())
}
