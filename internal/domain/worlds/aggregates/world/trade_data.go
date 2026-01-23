package world

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/tradegoods"
)

type tradeConnections struct {
	available   []tradegoods.TradeGood
	tradeRoutes int
	imported    map[coordinates.Global][]tradegoods.TradeGood
	expoted     map[coordinates.Global][]tradegoods.TradeGood
}

func newTradeConnections(tc ...classifications.Classification) *tradeConnections {
	td := tradeConnections{}
	// panic(tc)
	td.available = append(td.available, tradegoods.Available(tc...)...)
	td.imported = make(map[coordinates.Global][]tradegoods.TradeGood)
	td.expoted = make(map[coordinates.Global][]tradegoods.TradeGood)
	return &td
}

func (w *World) CreateTradeConnection(gc coordinates.Global, imported, exported []tradegoods.TradeGood) error {
	if w.trade == nil {
		fmt.Println("=========", w.imported)
		w.trade = newTradeConnections(w.imported.TradeCodes()...)
	}
	w.trade.create(gc, imported, exported)
	return nil
}

func (t *tradeConnections) create(gc coordinates.Global, imported, exported []tradegoods.TradeGood) {
	if len(imported) > 0 {
		t.imported[gc] = imported
	}
	if len(exported) > 0 {
		t.expoted[gc] = exported
	}
	t.tradeRoutes++
}

func (t *tradeConnections) read(gc coordinates.Global) ([]tradegoods.TradeGood, []tradegoods.TradeGood) {
	i := []tradegoods.TradeGood{}
	e := []tradegoods.TradeGood{}
	if tg, ok := t.imported[gc]; ok {
		i = tg
	} else {
		i = nil
	}
	if tg, ok := t.expoted[gc]; ok {
		e = tg
	} else {
		e = nil
	}
	return i, e
}
