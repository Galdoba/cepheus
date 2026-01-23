package traderoute

import (
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

// func CalculateImport(local, distant Market) []tradegoods.TradeGood {
// 	importing := []tradegoods.TradeGood{}
// 	if coordinates.Equal(local.Coordinates(), distant.Coordinates()) {
// 		return importing
// 	}
// 	localTC := local.TradeCodes()
// 	if local.TravelZone() == travelzone.Amber {
// 		localTC = append(localTC, classifications.Amber)
// 	}
// 	if local.TravelZone() == travelzone.Red {
// 		localTC = append(localTC, classifications.Red)
// 	}
// 	for _, goods := range tradegoods.Available(distant.TradeCodes()...) {
// 		sf, ok := tradegoods.SaleFactor(goods, localTC...)
// 		if ok && sf >= 0 {
// 			importing = append(importing, goods)
// 		}
// 	}
// 	return importing
// }

// func CalculateExport(local, distant Market) []tradegoods.TradeGood {
// 	exporting := []tradegoods.TradeGood{}
// 	if coordinates.Equal(local.Coordinates(), distant.Coordinates()) {
// 		return exporting
// 	}
// 	distantTC := distant.TradeCodes()
// 	if distant.TravelZone() == travelzone.Amber {
// 		distantTC = append(distantTC, classifications.Amber)
// 	}
// 	if distant.TravelZone() == travelzone.Red {
// 		distantTC = append(distantTC, classifications.Red)
// 	}
// 	for _, goods := range tradegoods.Available(distant.TradeCodes()...) {
// 		// sf, wantSale := tradegoods.SaleFactor(goods, distantTC...)
// 		// pf, wantPurchase := tradegoods.PurchseFactor(goods, local.TradeCodes()...)
// 		// if ok && sf >= 0 {
// 		// 	exporting = append(exporting, goods)
// 		// }
// 	}
// 	return exporting
// }

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
