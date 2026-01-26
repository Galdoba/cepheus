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
	tgList := tradegoods.List(false, true, false)
	producedAtSource := Available(source, tgList...)
	producedAtTarget := Available(target, tgList...)
	requestedAtSource := Importing(source, producedAtTarget...)
	requestedAtTarget := Importing(target, producedAtSource...)
	return clearDuplicated(requestedAtSource, requestedAtTarget)
}

func Importing(m Market, tgList ...tradegoods.TradeGood) []tradegoods.TradeGood {
	importing := []tradegoods.TradeGood{}
	for _, goods := range tgList {
		dm := 0
		for _, mcl := range m.TradeCodes() {
			if val, ok := goods.SaleDM[mcl]; ok {
				dm += val
			}
		}
		if dm > 0 {
			importing = append(importing, goods)
		}
	}
	return importing
}

func Available(m Market, tgList ...tradegoods.TradeGood) []tradegoods.TradeGood {
	available := []tradegoods.TradeGood{}
goods:
	for _, tg := range tgList {
		for _, mcl := range m.TradeCodes() {
			if slices.Contains(tg.AvailableAt, mcl) {
				available = append(available, tg)
				continue goods
			}
		}
	}
	return available
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

func filterImport(m Market, importingProposal ...tradegoods.TradeGood) []tradegoods.TradeGood {
	local := make(map[string]bool)
	for _, ltg := range Available(m, importingProposal...) {
		local[ltg.Code] = true
	}
	filtered := []tradegoods.TradeGood{}
	for _, goods := range importingProposal {
		if local[goods.Code] {
			continue
		}
		filtered = append(filtered, goods)
	}
	return filtered
}
