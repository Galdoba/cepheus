package world

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/trade/cargo"
	"github.com/Galdoba/cepheus/pkg/trade/tradecodes"
	"github.com/Galdoba/cepheus/pkg/uwp"
)

func TestNew(t *testing.T) {
	u := uwp.New(uwp.FromString("A555154-E"))
	w := New()
	w.UWP = u
	w.TradeCodes = tradecodes.GenerateFromUWP(w.UWP)
	lot, err := cargo.NewCargo(dice.NewDicepool(), w.TradeCodes)

	fmt.Println(w)
	fmt.Println(err)
	fmt.Println(lot)

}
