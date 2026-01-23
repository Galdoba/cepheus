package traderoute

import (
	"slices"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/tradegoods"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/travelzone"
)

type Market interface {
	Coordinates() coordinates.Global
	TradeCodes() []classifications.Classification
	TravelZone() travelzone.Zone
}

func Calculate(source, target Market) ([]tradegoods.TradeGood, []tradegoods.TradeGood) {
	importing, exporting := []tradegoods.TradeGood{}, []tradegoods.TradeGood{}
	if coordinates.Equal(source.Coordinates(), target.Coordinates()) {
		return importing, exporting
	}
	sourceTC := append(source.TradeCodes(), classifications.Classification(string(source.TravelZone())))
	targetTC := append(target.TradeCodes(), classifications.Classification(string(target.TravelZone())))
	for _, goods := range tradegoods.Available(targetTC...) {
		saleDM, wantSale := tradegoods.SaleFactor(goods, targetTC...)
		purchDM, wantPurchse := tradegoods.PurchseFactor(goods, sourceTC...)
		if wantSale && wantPurchse && (saleDM+purchDM) > 0 {
			importing = append(importing, goods)
		}

		saleDM, wantSale = tradegoods.SaleFactor(goods, sourceTC...)
		purchDM, wantPurchse = tradegoods.PurchseFactor(goods, targetTC...)
		if wantSale && wantPurchse && (saleDM+purchDM) > 0 {
			exporting = append(exporting, goods)
		}
	}
	return clearDuplicated(importing, exporting)
}

func HasLink(source, target Market) bool {
	for _, tc1 := range source.TradeCodes() {
		switch tc1 {
		case classifications.In, classifications.Ht:
			for _, tc2 := range target.TradeCodes() {
				switch tc2 {
				case classifications.As, classifications.De, classifications.Ic, classifications.Ni:
					return true
				}
			}
		case classifications.Hi, classifications.Ri:
			for _, tc2 := range target.TradeCodes() {
				switch tc2 {
				case classifications.Ag, classifications.Ga, classifications.Wa:
					return true
				}
			}
		}
	}
	for _, tc1 := range target.TradeCodes() {
		switch tc1 {
		case classifications.In, classifications.Ht:
			for _, tc2 := range source.TradeCodes() {
				switch tc2 {
				case classifications.As, classifications.De, classifications.Ic, classifications.Ni:
					return true
				}
			}
		case classifications.Hi, classifications.Ri:
			for _, tc2 := range source.TradeCodes() {
				switch tc2 {
				case classifications.Ag, classifications.Ga, classifications.Wa:
					return true
				}
			}
		}
	}
	return false
}

func containsAnyTradeCode(pool []classifications.Classification, codes ...classifications.Classification) bool {
	for _, tc := range codes {
		if slices.Contains(pool, tc) {
			return true
		}
	}
	return false
}

func clearDuplicated(imp, exp []tradegoods.TradeGood) ([]tradegoods.TradeGood, []tradegoods.TradeGood) {
	goodsMap := make(map[string]int)
	for _, i := range imp {
		goodsMap[i.TradeGoodType]++
	}
	for _, e := range exp {
		goodsMap[e.TradeGoodType]++
	}
	newImp := []tradegoods.TradeGood{}
	for _, i := range imp {
		if goodsMap[i.TradeGoodType] > 1 {
			continue
		}
		newImp = append(newImp, i)
	}
	newExp := []tradegoods.TradeGood{}
	for _, e := range exp {
		if goodsMap[e.TradeGoodType] > 1 {
			continue
		}
		newExp = append(newExp, e)
	}
	return newImp, newExp
}
