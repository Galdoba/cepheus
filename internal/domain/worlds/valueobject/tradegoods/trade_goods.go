package tradegoods

import (
	"fmt"
	"slices"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	tc "github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/pkg/dice"
)

type TradeGood struct {
	Code                    string                    `json:"code"`
	TradeGoodType           string                    `json:"trade_good_type"`
	Descriptions            []string                  `json:"description"`
	AvailableAt             []tc.Classification       `json:"available_at"`
	MaximumSupplyMultiplier int                       `json:"maximum_supply_multiplier"`
	PurchaseDM              map[tc.Classification]int `json:"purchase_dm"`
	SaleDM                  map[tc.Classification]int `json:"sale_dm"`
	MaximumRiskAssesmentDM  int                       `json:"maximum_risk_assesment_dm"`
	DangerousGoodsDM        int                       `json:"dangerous_goods_dm"`
	IncrementBase           []int                     `json:"increment_base"`
	IncrementMultiplier     []int                     `json:"increment_multiplier"`
	IncrementAddition       []int                     `json:"increment_addition"`
	BasePrice               []int                     `json:"base_price"`
	PurchaseSkill           string                    `json:"purchase_skill,omitempty"`
	PurchaseDifficulty      int                       `json:"purchase_difficulty,omitempty"`
	SaleSkill               string                    `json:"sale_skill,omitempty"`
	SaleDificulty           int                       `json:"sale_dificulty,omitempty"`
}

func NewRandom(d *dice.Roller) TradeGood {
	return tradeGoods[d.ConcatRoll("d66")]
}

func New(code string) (TradeGood, error) {
	if !validCode(code) {
		return TradeGood{}, fmt.Errorf("invalid tradegoods code provided: '%v'", code)
	}
	return tradeGoods[code], nil
}

func validCode(code string) bool {
	return slices.Contains([]string{
		"11", "12", "13", "14", "15", "16",
		"21", "22", "23", "24", "25", "26",
		"31", "32", "33", "34", "35", "36",
		"41", "42", "43", "44", "45", "46",
		"51", "52", "53", "54", "55", "56",
		"61", "62", "63", "64", "65", "66",
	}, code)
}

func Available(cls ...tc.Classification) []TradeGood {
	goodsAvailable := []TradeGood{}
goods_loop:
	for _, goods := range List(true, true, true) {
		switch goods.Code {
		case "11", "12", "13", "14", "15", "16", "66":
			goodsAvailable = append(goodsAvailable, goods)
		default:
			for _, cl := range cls {
				if slices.Contains(goods.AvailableAt, cl) {
					goodsAvailable = append(goodsAvailable, goods)
					continue goods_loop
				}
			}

		}
	}
	return goodsAvailable
}

func List(common, special, illegal bool) []TradeGood {
	goodsAvailable := []TradeGood{}
	for k, goods := range tradeGoods {
		switch k {
		case "11", "12", "13", "14", "15", "16":
			if common {
				goodsAvailable = append(goodsAvailable, goods)
			}
		case "61", "62", "63", "64", "65":
			if illegal {
				goodsAvailable = append(goodsAvailable, goods)
			}
		case "66":
			goodsAvailable = append(goodsAvailable, goods)
		default:
			if special {
				goodsAvailable = append(goodsAvailable, goods)
			}
		}
	}
	return goodsAvailable
}

func SaleFactor(tg TradeGood, tcPool ...classifications.Classification) (int, bool) {
	sf := -999
	present := false
	for _, tc := range tcPool {
		if val, ok := tg.SaleDM[tc]; ok {
			sf = max(sf, val)
			present = true
		}
	}
	return sf, present
}

func PurchseFactor(tg TradeGood, tcPool ...classifications.Classification) (int, bool) {
	pf := -999
	present := false
	for _, tc := range tcPool {
		if val, ok := tg.PurchaseDM[tc]; ok {
			pf = max(pf, val)
			present = true
		}
	}
	return pf, present
}

func Types(codes ...string) []string {
	types := []string{}
	for _, code := range codes {
		tg, err := New(code)
		if err != nil {
			types = append(types, fmt.Sprintf("invalid code <%v>", code))
			continue
		}
		types = append(types, tg.TradeGoodType)
	}
	return types
}

var tradeGoods = map[string]TradeGood{
	"11": TG_11,
	"12": TG_12,
	"13": TG_13,
	"14": TG_14,
	"15": TG_15,
	"16": TG_16,
	"21": TG_21,
	"22": TG_22,
	"23": TG_23,
	"24": TG_24,
	"25": TG_25,
	"26": TG_26,
	"31": TG_31,
	"32": TG_32,
	"33": TG_33,
	"34": TG_34,
	"35": TG_35,
	"36": TG_36,
	"41": TG_41,
	"42": TG_42,
	"43": TG_43,
	"44": TG_44,
	"45": TG_45,
	"46": TG_46,
	"51": TG_51,
	"52": TG_52,
	"53": TG_53,
	"54": TG_54,
	"55": TG_55,
	"56": TG_56,
	"61": TG_61,
	"62": TG_62,
	"63": TG_63,
	"64": TG_64,
	"65": TG_65,
}

var TG_11 = TradeGood{
	Code:                    "11",
	TradeGoodType:           goodsType("11"),
	Descriptions:            descriptions("11"),
	AvailableAt:             tc.All(),
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Ht: 3,
		tc.Ri: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ni: 2,
		tc.Lt: 1,
		tc.Po: 1,
	},
	MaximumRiskAssesmentDM: 0,
	DangerousGoodsDM:       -6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 4, 2},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{12000, 16000, 20000, 24000, 28000},
}

var TG_12 = TradeGood{
	Code:                    "12",
	TradeGoodType:           goodsType("12"),
	Descriptions:            descriptions("12"),
	AvailableAt:             tc.All(),
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Na: 2,
		tc.In: 5,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ni: 3,
		tc.Ag: 2,
	},
	MaximumRiskAssesmentDM: 0,
	DangerousGoodsDM:       -6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 4, 2},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{8000, 9000, 10000, 11000, 12000},
}

var TG_13 = TradeGood{
	Code:                    "13",
	TradeGoodType:           goodsType("13"),
	Descriptions:            descriptions("13"),
	AvailableAt:             tc.All(),
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Na: 2,
		tc.In: 5,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ni: 3,
		tc.Hi: 2,
	},
	MaximumRiskAssesmentDM: 0,
	DangerousGoodsDM:       -6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 5, 3},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{16000, 18000, 20000, 22000, 24000},
}
var TG_14 = TradeGood{
	Code:                    "14",
	TradeGoodType:           goodsType("14"),
	Descriptions:            descriptions("14"),
	AvailableAt:             tc.All(),
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 3,
		tc.Ga: 2,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Po: 2,
	},
	MaximumRiskAssesmentDM: 0,
	DangerousGoodsDM:       -6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{28, 24, 20, 10, 6},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{1000, 3000, 5000, 7000, 9000},
}

var TG_15 = TradeGood{
	Code:                    "15",
	TradeGoodType:           goodsType("15"),
	Descriptions:            descriptions("15"),
	AvailableAt:             tc.All(),
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 3,
		tc.Wa: 2,
		tc.Ga: 1,
		tc.As: -4,
	},
	SaleDM: map[tc.Classification]int{
		tc.As: 1,
		tc.Fl: 1,
		tc.Ic: 1,
		tc.Hi: 1,
	},
	MaximumRiskAssesmentDM: 0,
	DangerousGoodsDM:       -6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{24, 20, 20, 16, 8},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{125, 250, 500, 750, 1250},
}

var TG_16 = TradeGood{
	Code:                    "16",
	TradeGoodType:           goodsType("16"),
	Descriptions:            descriptions("16"),
	AvailableAt:             tc.All(),
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 4,
		// tc.Ic: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 3,
		tc.Ni: 1,
	},
	MaximumRiskAssesmentDM: 0,
	DangerousGoodsDM:       -6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{28, 24, 20, 16, 8},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{250, 500, 1000, 1500, 2000},
}

var TG_21 = TradeGood{
	Code:                    "21",
	TradeGoodType:           goodsType("21"),
	Descriptions:            descriptions("21"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Ht: 3,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ni: 1,
		tc.Ri: 2,
		tc.As: 3,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 2, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{25000, 50000, 100000, 125000, 150000},
}

var TG_22 = TradeGood{
	Code:                    "22",
	TradeGoodType:           goodsType("22"),
	Descriptions:            descriptions("22"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Ht: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.As: 2,
		tc.Ni: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 3, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{25000, 50000, 75000, 90000, 100000},
}

var TG_23 = TradeGood{
	Code:                    "23",
	TradeGoodType:           goodsType("23"),
	Descriptions:            descriptions("23"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 1,
		// tc.Ht: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Hi: 1,
		tc.Ri: 2,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 2, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{50000, 75000, 100000, 125000, 150000},
}

var TG_24 = TradeGood{
	Code:                    "24",
	TradeGoodType:           goodsType("24"),
	Descriptions:            descriptions("24"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		// tc.In: 0,
		tc.Ht: 2,
	},
	SaleDM: map[tc.Classification]int{
		tc.Po:    1,
		tc.Amber: 2,
		tc.Red:   4,
	},
	MaximumRiskAssesmentDM: 3,
	DangerousGoodsDM:       0,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{7, 6, 5, 3, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{50000, 100000, 150000, 200000, 250000},
}

var TG_25 = TradeGood{
	Code:                    "25",
	TradeGoodType:           goodsType("25"),
	Descriptions:            descriptions("25"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		// tc.In: 0,
		tc.Ht: 2,
	},
	SaleDM: map[tc.Classification]int{
		tc.As: 2,
		tc.Ri: 2,
	},
	MaximumRiskAssesmentDM: 3,
	DangerousGoodsDM:       0,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{5, 5, 5, 2, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{100000, 140000, 180000, 200000, 250000},
}

var TG_26 = TradeGood{
	Code:                    "26",
	TradeGoodType:           goodsType("26"),
	Descriptions:            descriptions("26"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Wa},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 1,
		tc.Wa: 2,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 2,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 3, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{10000, 25000, 50000, 60000, 80000},
}

var TG_31 = TradeGood{
	Code:                    "31",
	TradeGoodType:           goodsType("31"),
	Descriptions:            descriptions("31"),
	AvailableAt:             []tc.Classification{tc.As, tc.De, tc.Ic},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 2,
		tc.De: 1,
		tc.Ic: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 3,
		tc.Ri: 2,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -1,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{7, 6, 5, 3, 2},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{5000, 10000, 20000, 30000, 45000},
}

var TG_32 = TradeGood{
	Code:                    "32",
	TradeGoodType:           goodsType("32"),
	Descriptions:            descriptions("32"),
	AvailableAt:             []tc.Classification{tc.Ht},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.Ht: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.As: 1,
		tc.Ic: 1,
		tc.Ri: 2,
	},
	MaximumRiskAssesmentDM: 3,
	DangerousGoodsDM:       1,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 1, 0},
	IncrementAddition:      []int{2, 1, 0, 0, 1},
	BasePrice:              []int{100000, 200000, 250000, 350000, 500000},
}

var TG_33 = TradeGood{
	Code:                    "33",
	TradeGoodType:           goodsType("33"),
	Descriptions:            descriptions("33"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Ga},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 2,
		// tc.Ga: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Lo: 3,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 6, 3},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{2500, 5000, 10000, 12500, 15000},
}

var TG_34 = TradeGood{
	Code:                    "34",
	TradeGoodType:           goodsType("34"),
	Descriptions:            descriptions("34"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Ga, tc.Wa},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 2,
		// tc.Ga: 0,
		tc.Wa: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 2,
		tc.Hi: 2,
	},
	MaximumRiskAssesmentDM: 3,
	DangerousGoodsDM:       2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{14, 12, 10, 5, 2},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{5000, 10000, 20000, 30000, 50000},
}

var TG_35 = TradeGood{
	Code:                    "35",
	TradeGoodType:           goodsType("35"),
	Descriptions:            descriptions("35"),
	AvailableAt:             []tc.Classification{tc.Hi},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.Hi: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 4,
	},
	MaximumRiskAssesmentDM: 3,
	DangerousGoodsDM:       2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 1, 0},
	IncrementAddition:      []int{2, 1, 0, 0, 1},
	BasePrice:              []int{50000, 100000, 200000, 250000, 500000},
}

var TG_36 = TradeGood{
	Code:                    "36",
	TradeGoodType:           goodsType("36"),
	Descriptions:            descriptions("36"),
	AvailableAt:             []tc.Classification{tc.Ht, tc.Hi},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.Ht: 2,
		// tc.Hi: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Po: 1,
		tc.Ri: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 2, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{10000, 30000, 50000, 75000, 100000},
}

var TG_41 = TradeGood{
	Code:                    "41",
	TradeGoodType:           goodsType("41"),
	Descriptions:            descriptions("41"),
	AvailableAt:             []tc.Classification{tc.De, tc.Fl, tc.Ic, tc.Wa},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.De: 2,
		// tc.Fl: 0,
		// tc.Ic: 0,
		// tc.Wa: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Ag: 1,
		tc.Lt: 2,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 8, 4},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{2500, 5000, 10000, 20000, 30000},
}

var TG_42 = TradeGood{
	Code:                    "42",
	TradeGoodType:           goodsType("42"),
	Descriptions:            descriptions("42"),
	AvailableAt:             []tc.Classification{tc.As, tc.De, tc.Hi, tc.Wa},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 2,
		// tc.De: 0,
		tc.Hi: 1,
		// tc.Wa: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 2,
		tc.Lt: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       3,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 0, 0},
	IncrementAddition:      []int{3, 2, 0, 2, 1},
	BasePrice:              []int{25000, 50000, 100000, 200000, 500000},
}

var TG_43 = TradeGood{
	Code:                    "43",
	TradeGoodType:           goodsType("43"),
	Descriptions:            descriptions("43"),
	AvailableAt:             []tc.Classification{tc.In},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 2,
		tc.Ni: 1,
	},
	MaximumRiskAssesmentDM: 1,
	DangerousGoodsDM:       0,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 3, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{1000, 3000, 7000, 9000, 10000},
}

var TG_44 = TradeGood{
	Code:                    "44",
	TradeGoodType:           goodsType("44"),
	Descriptions:            descriptions("44"),
	AvailableAt:             []tc.Classification{tc.As, tc.De, tc.Ic, tc.Fl},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 3,
		tc.De: 1,
		tc.Ic: 2,
		// tc.Fl: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 3,
		tc.In: 2,
		tc.Ht: 1,
	},
	MaximumRiskAssesmentDM: 3,
	DangerousGoodsDM:       4,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 1, 0},
	IncrementAddition:      []int{2, 1, 0, 0, 1},
	BasePrice:              []int{10000, 25000, 50000, 75000, 100000},
}

var TG_45 = TradeGood{
	Code:                    "45",
	TradeGoodType:           goodsType("45"),
	Descriptions:            descriptions("45"),
	AvailableAt:             []tc.Classification{tc.As, tc.De, tc.Lo},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 2,
		// tc.De: 0,
		tc.Lo: 2,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 3,
		tc.Ht: 1,
		tc.Ni: -2,
		tc.Ag: -3,
	},
	MaximumRiskAssesmentDM: 4,
	DangerousGoodsDM:       3,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 0, 0},
	IncrementAddition:      []int{3, 2, 0, 1, 1},
	BasePrice:              []int{500000, 750000, 1000000, 1250000, 1500000},
}

var TG_46 = TradeGood{
	Code:                    "46",
	TradeGoodType:           goodsType("46"),
	Descriptions:            descriptions("46"),
	AvailableAt:             []tc.Classification{tc.In},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ag: 2,
		tc.Ht: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       1,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{7, 6, 5, 2, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{150000, 300000, 400000, 500000, 650000},
}

var TG_51 = TradeGood{
	Code:                    "51",
	TradeGoodType:           goodsType("51"),
	Descriptions:            descriptions("51"),
	AvailableAt:             []tc.Classification{tc.Ga, tc.De, tc.Wa},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		// tc.Ga: 0,
		tc.De: 2,
		// tc.Wa: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Hi: 2,
		tc.Ri: 3,
		tc.Po: 3,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -1,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{12, 10, 10, 6, 2},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{1000, 3000, 6000, 9000, 12000},
}

var TG_52 = TradeGood{
	Code:                    "52",
	TradeGoodType:           goodsType("52"),
	Descriptions:            descriptions("52"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Ni},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 7,
		// tc.Ni: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Hi: 3,
		tc.Na: 2,
	},
	MaximumRiskAssesmentDM: 1,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{24, 20, 20, 12, 6},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{1000, 2000, 3000, 4000, 5000},
}

var TG_53 = TradeGood{
	Code:                    "53",
	TradeGoodType:           goodsType("53"),
	Descriptions:            descriptions("53"),
	AvailableAt:             []tc.Classification{tc.As, tc.Ic},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 4,
		// tc.Ic: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 3,
		tc.Ni: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{20, 20, 20, 10, 6},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{1000, 2500, 5000, 7500, 10000},
}

var TG_54 = TradeGood{
	Code:                    "54",
	TradeGoodType:           goodsType("54"),
	Descriptions:            descriptions("54"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.De, tc.Wa},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 2,
		// tc.De: 0,
		tc.Wa: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Ht: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{14, 12, 10, 8, 3},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{5000, 10000, 20000, 35000, 50000},
}

var TG_55 = TradeGood{
	Code:                    "55",
	TradeGoodType:           goodsType("55"),
	Descriptions:            descriptions("55"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Ga},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 6,
		// tc.Ga: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 2,
		tc.In: 1,
	},
	MaximumRiskAssesmentDM: 1,
	DangerousGoodsDM:       -4,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{24, 20, 20, 12, 4},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{100, 500, 1000, 2000, 4000},
}

var TG_56 = TradeGood{
	Code:                    "56",
	TradeGoodType:           goodsType("56"),
	Descriptions:            descriptions("56"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 10,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 2,
		tc.Ht: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ni: 2,
		tc.Hi: 1,
	},
	MaximumRiskAssesmentDM: 2,
	DangerousGoodsDM:       -2,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{14, 12, 10, 6, 2},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{5000, 10000, 15000, 20000, 30000},
}

var TG_61 = TradeGood{
	Code:                    "61",
	TradeGoodType:           goodsType("61"),
	Descriptions:            descriptions("61"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Wa},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		// tc.Ag: 0,
		tc.Wa: 2,
	},
	SaleDM: map[tc.Classification]int{
		tc.In: 6,
	},
	MaximumRiskAssesmentDM: 4,
	DangerousGoodsDM:       4,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 0, 0},
	IncrementAddition:      []int{0, 0, 0, 2, 1},
	BasePrice:              []int{10000, 25000, 50000, 100000, 200000},
}

var TG_62 = TradeGood{
	Code:                    "62",
	TradeGoodType:           goodsType("62"),
	Descriptions:            descriptions("62"),
	AvailableAt:             []tc.Classification{tc.Ht},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.Ht: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.As:    4,
		tc.Ic:    4,
		tc.Ri:    8,
		tc.Amber: 6,
		tc.Red:   6,
	},
	MaximumRiskAssesmentDM: 5,
	DangerousGoodsDM:       5,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{2, 2, 1, 0, 0},
	IncrementAddition:      []int{0, 0, 0, 2, 1},
	BasePrice:              []int{100000, 150000, 250000, 400000, 650000},
}

var TG_63 = TradeGood{
	Code:                    "63",
	TradeGoodType:           goodsType("63"),
	Descriptions:            descriptions("63"),
	AvailableAt:             []tc.Classification{tc.As, tc.De, tc.Ga, tc.Wa},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.As: 1,
		tc.De: 1,
		tc.Ga: 1,
		tc.Wa: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 6,
		tc.Hi: 6,
	},
	MaximumRiskAssesmentDM: 4,
	DangerousGoodsDM:       6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 0, 0},
	IncrementAddition:      []int{2, 1, 0, 2, 1},
	BasePrice:              []int{25000, 50000, 100000, 200000, 300000},
}

var TG_64 = TradeGood{
	Code:                    "64",
	TradeGoodType:           goodsType("64"),
	Descriptions:            descriptions("64"),
	AvailableAt:             []tc.Classification{tc.Ag, tc.Ga, tc.Wa},
	MaximumSupplyMultiplier: 1,
	PurchaseDM: map[tc.Classification]int{
		tc.Ag: 2,
		// tc.Ga: 0,
		tc.Wa: 1,
	},
	SaleDM: map[tc.Classification]int{
		tc.Ri: 6,
		tc.Hi: 4,
	},
	MaximumRiskAssesmentDM: 4,
	DangerousGoodsDM:       4,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{1, 1, 1, 0, 0},
	IncrementAddition:      []int{1, 0, 0, 2, 1},
	BasePrice:              []int{10000, 25000, 50000, 100000, 200000},
}

var TG_65 = TradeGood{
	Code:                    "65",
	TradeGoodType:           goodsType("65"),
	Descriptions:            descriptions("65"),
	AvailableAt:             []tc.Classification{tc.In, tc.Ht},
	MaximumSupplyMultiplier: 5,
	PurchaseDM: map[tc.Classification]int{
		tc.In: 2,
		// tc.Ht: 0,
	},
	SaleDM: map[tc.Classification]int{
		tc.Po:    6,
		tc.Amber: 8,
		tc.Red:   10,
	},
	MaximumRiskAssesmentDM: 5,
	DangerousGoodsDM:       6,
	IncrementBase:          []int{1, 1, 1, 1, 1},
	IncrementMultiplier:    []int{6, 5, 5, 2, 1},
	IncrementAddition:      []int{0, 0, 0, 0, 0},
	BasePrice:              []int{50000, 100000, 150000, 300000, 450000},
}

var TG_66 = TradeGood{
	Code:                    "66",
	TradeGoodType:           goodsType("66"),
	Descriptions:            descriptions("66"),
	AvailableAt:             []tc.Classification{},
	MaximumSupplyMultiplier: 1,
	PurchaseDM:              map[tc.Classification]int{},
	SaleDM:                  map[tc.Classification]int{},
	MaximumRiskAssesmentDM:  0,
	DangerousGoodsDM:        0,
	IncrementBase:           []int{},
	IncrementMultiplier:     []int{},
	IncrementAddition:       []int{},
	BasePrice:               []int{},
}

func descriptions(code string) []string {
	descrMap := make(map[string][]string)

	descrMap["11"] = []string{"Calculators/Adding Machines", "Video Game and Entertainment Systems", "Personal and Commercial Computers", "Banking Machines and Security Systems", "Microprocessor Assemblies"}
	descrMap["12"] = []string{"Stamped/Poured Cogs and Sprockets", "Piping and Attachment Pieces", "Engine Components", "Pneumatics and Hydraulics", "Starship-Quality Components"}
	descrMap["13"] = []string{"Second Stage Components", "Uniforms/Clothing Products", "Residential Appliances", "Furniture/Storage Systems/Tools", "Vehicle/Survival Accessories"}
	descrMap["14"] = []string{"Foundation Stones and Base Elements", "Workable Metals", "Workable Alloys", "Fabricated Plastics", "Chemical Solutions or Compounds"}
	descrMap["15"] = []string{"Feed-grade Vegetation", "Food-grade Vegetation", "Pre-packaged Food and Drink", "Survival Rations and Storage-packed Liquids", "Junk Food/Soda/Beer"}
	descrMap["16"] = []string{"Bornite or Galena or Sedimentary Stone", "Chalcocite or Talc", "Bauxite, Coltan and Wolframite", "Acanthite, Cobaltite or Magnetite", "Chromite or Cinnabar"}

	descrMap["21"] = []string{"Circuitry Bundles", "Fibre-optic Components", "VR Computer and Sensor Packages", "Weapon Components", "Starship Bridge Components"}
	descrMap["22"] = []string{"Alloy and Plastic Tool Kits", "Starship Deckplate/Atmospheric Filters", "Fusion Conduits/Power Plant Shells", "Weapon Cores/Starship Hull", "Gravitic Gyros, Navigation Magnetics"}
	descrMap["23"] = []string{"High-Pressure or Temperature-Resistant Components", "Protective or Specialised Clothing", "Survival Equipment/Colonisation Kits", "Computerised Job-related Gear", "Starship Add-Ons/Powered Armour Components"}
	descrMap["24"] = []string{"(TL7 or less) Slug Weapons", "(TL10 or less) Slug Weapons", "(TL12 or less) Slug or Energy Weapons/Heavy Slug Weapons", "(TL15 or less) Slug or Energy Weapons/Explosives", "Artillery, Heavy Energy Weapons"}
	descrMap["25"] = []string{"Engine Components or Packages", "Seafaring or Mole Vehicle Components or Packages", "Air/Raft Components or Packages", "Grav-Vehicle Components", "Spacecraft Components"}
	descrMap["26"] = []string{"Organic Glues, Acids or Bases/Vegetable Oil", "Ethanol/Fructose Syrup", "Biodiesel/Cooking Compounds", "Oxygenated Cleaner/Biodegradable Concentrates", "Gelid Oxygen-Substitutes/Bio-fusion Cell Fuel"}

	descrMap["31"] = []string{"Rock Salt/Compressed Coal", "Graphite/Quartz", "Silica/Focuser-Quality Gems", "Photonics/Synthetic Gems", "Industrial Diamond/Jewellery-Quality Gems"}
	descrMap["32"] = []string{"Cybernetic Lubricants", "Cybernetic Components/Physical Augments", "Cyber-Prosthetics", "Cosmetic Prosthetics", "Real-Life Replacements and Augments"}
	descrMap["33"] = []string{"Beasts of Burden", "Untrained Riding Animals", "Trained Riding Animals/Common Pets", "Untrained Guard Animals", "Trained Guard Animals/Exotic Pets"}
	descrMap["34"] = []string{"Common Desserts/Rare Food Additives", "Common Desserts/Common Wine", "Rare Foods/Common Liquor", "Exotic Foods/Rare Desserts/Rare Liquor", "Exotic Desserts/Exotic Liquor"}
	descrMap["35"] = []string{"Rare Literature/Art", "Jewellery/Alien Textiles", "Rare Clothing/Home Decorations", "VR Electronic Entertainment Devices", "Exotic Furnishings/Exquisite Jewellery"}
	descrMap["36"] = []string{"Medical Uniforms/Disposable Tools", "Cosmetic Chemicals/Practitioner's Tools", "General Medical Equipment or Supplies", "Specialist Equipment or Supplies", "Micro-surgical Equipment or Supplies"}

	descrMap["41"] = []string{"Crude Oil/Diesel", "Refined Kerosene/Purified Oil", "Gasoline/Machine Lubricants", "Jet Fuel/Gelid Adhesives", "Rocket Fuel/Power Plant Starter Charges"}
	descrMap["42"] = []string{"OTC Drugs/Antibiotics", "Antivenin/Prescription Medications", "Prescription Medications/Surgical", "Anagathics", "Psi-Related Drugs/Viral Therapy Doses"}
	descrMap["43"] = []string{"Rubber/Vinyl Spooling", "Insulation/Polyurethane Foam", "Poured Plastics/Synthetic Fibre Spools", "Kevlar/Teflon", "Advanced Ballistic Weave"}
	descrMap["44"] = []string{"Bismuth/Indium", "Beryllium/Silver", "Ruthenium/Rhenium", "Gold/Osmium/Iridium", "Platinum/Rhodium"}
	descrMap["45"] = []string{"Nuclear Waste/Deactivated Materials", "Industrial Isotopes", "Medical Isotopes/Reactor-Grade Uranium", "Weapons-Grade Plutonium/ Fusion Cell Rods", "Superweapon-grade Isotopes"}
	descrMap["46"] = []string{"Automated Robotics/Cargo Drones", "Industrial or Personal Drones", "Combat or Guardian Drones", "Scout and Sensor Drones", "Advanced Robotics"}

	descrMap["51"] = []string{"Table Salt/Black Pepper", "Adobo/Basil/Sage", "Aniseed/Curry/Fennel/White Pepper", "Cinnamon/Marjoram/Wasabi", "Black Salt/Saffron/Alien Flavourings"}
	descrMap["52"] = []string{"Yarn/Wool/Canvas", "Animal-based Fabrics", "Cotton or Flax-based Fabrics", "Synthetic Silks/Finished Common Clothing", "Organic Silk/Satin/Finished Fine Clothing"}
	descrMap["53"] = []string{"Lead/Zinc", "Copper/Tin", "Nickel/Sodium/Tungsten", "Gold/Silver/Ilmenite", "Platinum/Uranium"}
	descrMap["54"] = []string{"Aluminium/Brass/Calcium", "Carbonate/Magnesium/Meteoric Iron", "Marble/Potassium/Titanium", "Stellite/Tombac", "Depleted Uranium/Ceramic-Alloy"}
	descrMap["55"] = []string{"Low-grade Rough Cuts/Construction Scrap", "High-Grade Rough-Cut", "Construction-grade Timber", "Furniture-grade Timber/Rare Grades", "Exotics (Pernambuco, White Mahogany, etc.)"}
	descrMap["56"] = []string{"Wheeled Repair Components", "Tracked Repair Components", "Wheeled Components or Packages", "Wheeled Vehicles/Tracked Components or Packages", "Tracked Vehicles"}

	descrMap["61"] = []string{"Herbal Stimulants/Ultra-Caffeine", "Raw Growth Hormones", "Chemical Solvents/Protein Duplexer Steroids", "Bio-Acid/Pheromone Extracts", "Genetic Mutagens/Organic Toxins"}
	descrMap["62"] = []string{"Unlicensed Augment Tools and Parts", "Physical Enhancement Tissues", "Unlicensed Augmentatives/Combat Implant Additives", "Combat Prosthetics/Surgical Duplications", "Mimicry Augmetics"}
	descrMap["63"] = []string{"Herbal Stimulants/Biological Hallucinogens", "Chemical Depressants/Natural Narcotics", "Chemical Stimulants and Hallucinogens", "Designer Narcotics", "Alien Synthetics/Psi-Drugs"}
	descrMap["64"] = []string{"Anti-Governmental Propaganda/Endangered Animal Products", "Black-data Recordings/Slaving Gear", "Extinct Animal Products", "BTL Devices/Cloning Equipment", "Forbidden Pleasures"}
	descrMap["65"] = []string{"Chain-drive Weaponry/Armour-Piercing Ammunition", "Protected Technologies/Explosive or Incendiary Ammunition", "Synthetic Poisons/Personal-scale Mass Trauma Explosives", "Arclight Weaponry/Biological or Chemical Weaponry/Naval Starship Weaponry", "Disintegrators/Psi-Weaponry/Weapons of Mass Destruction"}
	descrMap["66"] = []string{"Exotics"}

	return descrMap[code]
}

func goodsType(code string) string {
	typesMap := make(map[string]string)

	typesMap["11"] = "Common Electronics"
	typesMap["12"] = "Common Machine Parts"
	typesMap["13"] = "Common Manufactored Goods"
	typesMap["14"] = "Common Raw Materials"
	typesMap["15"] = "Common Consumables"
	typesMap["16"] = "Common Ore"
	typesMap["21"] = "Advanced Electronics"
	typesMap["22"] = "Advanced Machine Parts"
	typesMap["23"] = "Advanced Manufactored Goods"
	typesMap["24"] = "Advanced Weapons"
	typesMap["25"] = "Advanced Vechicles"
	typesMap["26"] = "Biochemicals"
	typesMap["31"] = "Crystals & Gems"
	typesMap["32"] = "Cybernetics"
	typesMap["33"] = "Live Animals"
	typesMap["34"] = "Luxury Consumables"
	typesMap["35"] = "Luxury Goods"
	typesMap["36"] = "Medical Supplies"
	typesMap["41"] = "Petrochemicals"
	typesMap["42"] = "Pharmaceuticals"
	typesMap["43"] = "Polymers"
	typesMap["44"] = "Precious Metals"
	typesMap["45"] = "Radioactives"
	typesMap["46"] = "Robots"
	typesMap["51"] = "Spices"
	typesMap["52"] = "Textiles"
	typesMap["53"] = "Uncommon Ore"
	typesMap["54"] = "Uncommon Raw Material"
	typesMap["55"] = "Wood"
	typesMap["56"] = "Vechicles"
	typesMap["61"] = "Biochemicals (illegal)"
	typesMap["62"] = "Cybernetics (illegal)"
	typesMap["63"] = "Drugs (illegal)"
	typesMap["64"] = "Luxuries (illegal)"
	typesMap["65"] = "Weapons (illegal)"
	typesMap["66"] = "Exotics"

	return typesMap[code]
}
