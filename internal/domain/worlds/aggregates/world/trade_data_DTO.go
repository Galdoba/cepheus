package world

import (
	"fmt"
	"slices"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/tradegoods"
)

type TradeDataDTO struct {
	AvailableGoods []string            `json:"available_goods,omitempty"`
	ActiveRoutes   int                 `json:"active_routes,omitempty"`
	Imported       map[string][]string `json:"imported,omitempty"`
	Exported       map[string][]string `json:"exported,omitempty"`
}

func (t *tradeConnections) ToDTO() TradeDataDTO {
	dto := TradeDataDTO{}
	for _, goods := range t.available {
		dto.AvailableGoods = append(dto.AvailableGoods, goods.Code)
	}
	slices.Sort(dto.AvailableGoods)
	dto.ActiveRoutes = t.tradeRoutes
	dto.Imported = make(map[string][]string)
	dto.Exported = make(map[string][]string)
	for crd, goods := range t.imported {
		for _, specific := range goods {
			dto.Imported[crd.ToCube().ToGlobal().DatabaseKey()] = append(dto.Imported[crd.ToCube().ToGlobal().DatabaseKey()], specific.Code)
		}
	}
	for crd, goods := range t.expoted {
		for _, specific := range goods {
			dto.Exported[crd.ToCube().ToGlobal().DatabaseKey()] = append(dto.Exported[crd.ToCube().ToGlobal().DatabaseKey()], specific.Code)
		}
	}
	return dto
}

func TradeConnectionsFromDTO(dto TradeDataDTO) (*tradeConnections, error) {
	t := newTradeConnections()
	for _, code := range dto.AvailableGoods {
		goods, err := tradegoods.New(code)
		if err != nil {
			return nil, fmt.Errorf("failed to create tradegoods: %v", err)
		}
		t.available = append(t.available, goods)
	}
	t.tradeRoutes = dto.ActiveRoutes
	for key, codes := range dto.Imported {
		crd, err := coordinates.GlobalFromDatabaseKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get coordinates: %v", err)
		}
		for _, code := range codes {
			goods, err := tradegoods.New(code)
			if err != nil {
				return nil, fmt.Errorf("failed to create tradegoods: %v", err)
			}
			t.imported[crd] = append(t.imported[crd], goods)
		}
	}
	for key, codes := range dto.Exported {
		crd, err := coordinates.GlobalFromDatabaseKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get coordinates: %v", err)
		}
		for _, code := range codes {
			goods, err := tradegoods.New(code)
			if err != nil {
				return nil, fmt.Errorf("failed to create tradegoods: %v", err)
			}
			t.expoted[crd] = append(t.expoted[crd], goods)
		}
	}
	return t, nil
}
