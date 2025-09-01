package goods

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/ehex"
	"github.com/Galdoba/cepheus/pkg/trade/tradecodes"
	"github.com/Galdoba/cepheus/pkg/uwp"
)

type Goods struct {
	uwp             string
	codes           []string
	imbalance       string
	population      int
	techLevel       int
	basePrice       int
	productionPrice int
}

func New(source string) *Goods {
	g := Goods{}
	g.basePrice = 3000
	u := uwp.New(uwp.FromString(source))
	g.uwp = source
	g.codes = tradecodes.GenerateFromUWP(u)
	g.population = u.ValueOf(uwp.Pops)
	g.techLevel = u.ValueOf(uwp.TL)
	g.calculateProductionCost()
	return &g
}

func (g *Goods) calculateProductionCost() {
	price := g.basePrice
	for _, tc := range g.codes {
		switch tc {
		case "Ag", "As", "Hi", "In", "Po":
			price -= 1000
		case "Ba", "De", "Fl", "Lo", "Ni", "Ri", "Va":
			price += 1000
		}
	}
	g.productionPrice = price + (g.techLevel * 100)
}

func (g *Goods) ID() string {
	s := ehex.FromInt(g.techLevel).Code()
	s += "-"
	if len(g.codes) == 0 && g.imbalance == "" {
		s += "[n/a]"
	}
	for _, code := range g.codes {
		s += code + " "
	}
	if g.imbalance != "" {
		s += "(" + g.imbalance + ") "
	}
	s += fmt.Sprintf("Cr%v", g.productionPrice)
	return s
}
