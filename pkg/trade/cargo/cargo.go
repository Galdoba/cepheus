package cargo

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
)

const (
	TGT_Agricultural = iota
	TGT_Criminal
	TGT_Desert
	TGT_Industrial
	TGT_Information
	TGT_Processed
	TGT_Resources
	TGT_Water
	TGT_Vacuum
)

type Lot struct {
	Tonns      int
	Content    string
	CargoCode  string
	Hazards    []string
	Intrigue   string
	Event      string
	needReroll bool
	Stolen     bool
	Criminal   bool
}

func NewCargo(dp *dice.Dicepool, tradecodes []string) (*Lot, error) {
	l := Lot{}
	tables := selectTables(dp, tradecodes)
	if err := l.rollCargo(dp, tables); err != nil {
		return nil, err
	}
	if l.Content == "" {
		fmt.Println(l, tables, tradecodes)
		panic(0)
	}
	return &l, nil
}

func selectTables(dp *dice.Dicepool, codes []string) []int {
	choose := []int{}
	chosen := []int{}
	for _, code := range codes {
		switch code {
		case "Ag", "Ga":
			choose = []int{
				TGT_Information,
				TGT_Information,
				TGT_Agricultural,
				TGT_Agricultural,
				TGT_Agricultural,
				TGT_Agricultural,
				TGT_Processed,
				TGT_Resources,
				TGT_Resources,
				TGT_Processed,
				TGT_Criminal}
		case "In", "Ht":
			choose = []int{
				TGT_Industrial,
				TGT_Industrial,
				TGT_Industrial,
				TGT_Industrial,
				TGT_Industrial,
				TGT_Resources,
				TGT_Processed,
				TGT_Processed,
				TGT_Information,
				TGT_Information,
				TGT_Criminal}
		case "Fl", "Ri", "Wa":
			choose = []int{
				TGT_Water,
				TGT_Water,
				TGT_Water,
				TGT_Water,
				TGT_Water,
				TGT_Resources,
				TGT_Processed,
				TGT_Processed,
				TGT_Information,
				TGT_Information,
				TGT_Criminal}
		case "Va", "As":
			choose = []int{
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Industrial,
				TGT_Industrial,
				TGT_Processed,
				TGT_Processed,
				TGT_Information,
				TGT_Criminal}
		case "Ba", "De", "Na":
			choose = []int{
				TGT_Industrial,
				TGT_Industrial,
				TGT_Desert,
				TGT_Desert,
				TGT_Desert,
				TGT_Desert,
				TGT_Desert,
				TGT_Desert,
				TGT_Processed,
				TGT_Information,
				TGT_Criminal}
		case "Ic":
			choose = []int{
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Vacuum,
				TGT_Industrial,
				TGT_Agricultural,
				TGT_Resources,
				TGT_Resources,
				TGT_Processed,
				TGT_Information,
				TGT_Criminal}
		default:
			choose = []int{
				TGT_Industrial,
				TGT_Industrial,
				TGT_Water,
				TGT_Agricultural,
				TGT_Resources,
				TGT_Resources,
				TGT_Processed,
				TGT_Processed,
				TGT_Information,
				TGT_Vacuum,
				TGT_Criminal,
			}
		}
		chosen = append(chosen, choose[dp.Sum("2d6")-2])
	}
	if len(chosen) == 0 {
		choose = []int{
			TGT_Industrial,
			TGT_Industrial,
			TGT_Water,
			TGT_Agricultural,
			TGT_Resources,
			TGT_Resources,
			TGT_Processed,
			TGT_Processed,
			TGT_Information,
			TGT_Vacuum,
			TGT_Criminal,
		}
		chosen = append(choose, choose[dp.Sum("2d6")-2])
	}
	code := fmt.Sprintf("1d%v", len(chosen))
	fmt.Printf("code: %v of %v\n", code, len(chosen))
	return chosen
}

func (l *Lot) rollCargo(dp *dice.Dicepool, tables []int) error {
	for {
		l.needReroll = false
		table := tables[dp.Sum(fmt.Sprintf("1d%v", len(tables)))-1]
		fmt.Println("roll table", table)
		switch table {
		default:
			return fmt.Errorf("unknown table: %v", table)
		case 0:
			l.rollArgicultural(dp)
		case 1:
			l.rollCriminal(dp)
			if l.needReroll {
				continue
			}
		case 2:
			l.rollDesert(dp)
		case 3:
			l.rollIndustrial(dp)
		case 4:
			l.rollInformation(dp)
		case 5:
			l.rollProcessed(dp)
		case 6:
			l.rollResources(dp)
		case 7:
			l.rollWater(dp)
		case 8:
			l.rollVacuum(dp)
		}
		break
	}
	return nil
}

type lotDetails struct {
	desc string
	cor  int
	fla  int
	fra  int
	rad  int
	per  int
}

func (l *Lot) setDetails(details lotDetails, dp *dice.Dicepool) {
	l.Content = details.desc
	if dp.Sum("2d6") <= details.cor {
		l.Hazards = append(l.Hazards, "Cor")
	}
	if dp.Sum("2d6") <= details.fla {
		l.Hazards = append(l.Hazards, "Fla")
	}
	if dp.Sum("2d6") <= details.fra {
		l.Hazards = append(l.Hazards, "Fra")
	}
	if dp.Sum("2d6") <= details.rad {
		l.Hazards = append(l.Hazards, "Rad")
	}
	if dp.Sum("2d6") <= details.per {
		l.Hazards = append(l.Hazards, "Per")
	}
}

func (l *Lot) rollArgicultural(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["4"] = lotDetails{"Wild plants (live)", 3, 2, 3, 2, 2}
	values["5"] = lotDetails{"Food plants (live)", 2, 2, 3, 0, 2}
	values["6"] = lotDetails{"Livestock", 2, 2, 4, 0, 0}
	values["7"] = lotDetails{"Livestock", 2, 2, 4, 0, 0}
	values["8"] = lotDetails{"Fibres", 0, 3, 2, 0, 3}
	values["9"] = lotDetails{"Meat", 2, 2, 0, 0, 9}
	values["10"] = lotDetails{"Meat", 2, 2, 0, 0, 9}
	values["11"] = lotDetails{"Vegetables", 0, 2, 3, 0, 7}
	values["12"] = lotDetails{"Grain", 0, 4, 2, 0, 5}
	values["13"] = lotDetails{"Grain", 0, 4, 2, 0, 5}
	values["14"] = lotDetails{"Grain", 0, 4, 2, 0, 5}
	values["15"] = lotDetails{"Processed foods", 0, 2, 2, 0, 3}
	values["16"] = lotDetails{"Processed foods", 0, 2, 2, 0, 3}
	values["17"] = lotDetails{"Forest products (wood)", 0, 4, 3, 2, 3}
	values["18"] = lotDetails{"Fruit", 2, 2, 4, 0, 7}
	values["19"] = lotDetails{"Textiles", 0, 5, 2, 0, 3}
	values["20"] = lotDetails{"Liquor/Wine", 2, 5, 4, 0, 4}
	values["21"] = lotDetails{"Herbs/spices", 0, 4, 3, 0, 5}
	values["22"] = lotDetails{"Pharmaceuticals", 3, 4, 4, 2, 5}
	values["23"] = lotDetails{"Rare plants (live)", 3, 2, 4, 2, 2}
	values["24"] = lotDetails{"Rare animals (live)", 3, 2, 5, 2, 0}
	key := fmt.Sprintf("%v", dp.Sum("4d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollCriminal(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["3"] = lotDetails{"Plants (poisonous)", 2, 2, 3, 0, 4}
	values["4"] = lotDetails{"Plants (large carnivorous)", 2, 2, 3, 0, 4}
	values["5"] = lotDetails{fmt.Sprintf("Fugitive(s) (%v persons)", dp.Sum("1d6")), 0, 0, 0, 0, 0}
	values["6"] = lotDetails{"Drugs (alien)", 3, 4, 3, 2, 3}
	values["7"] = lotDetails{"Drugs (hard)", 2, 4, 3, 0, 3}
	values["8"] = lotDetails{"Drugs (soft)", 2, 4, 3, 0, 3}
	values["11"] = lotDetails{"Information (illegal)", 0, 0, 5, 0, 0}
	values["12"] = lotDetails{"Weapons", 0, 4, 4, 2, 0}
	values["13"] = lotDetails{"Erotica/sex aids", 0, 2, 3, 0, 2}
	values["14"] = lotDetails{"Warbots", 0, 3, 3, 2, 0}
	values["15"] = lotDetails{"Chemical weapons", 3, 3, 5, 2, 3}
	values["16"] = lotDetails{"Atomic weapons", 2, 3, 4, 9, 0}
	values["17"] = lotDetails{"Bacteriological weapons", 0, 2, 6, 0, 4}
	values["18"] = lotDetails{"Genetic weapons", 0, 2, 6, 0, 4}
	key := fmt.Sprintf("%v", dp.Sum("3d6"))
	if key == "9" || key == "10" {
		l.needReroll = true
		l.Stolen = true
		return
	}
	l.Criminal = true
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollDesert(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["1"] = lotDetails{"Water condensers", 0, 0, 3, 0, 0}
	values["2"] = lotDetails{"Hydroponics equipment", 0, 0, 4, 3, 0}
	values["3"] = lotDetails{"Stilsuits", 0, 0, 2, 0, 0}
	values["4"] = lotDetails{"Stellar power systems", 0, 2, 3, 0, 0}
	values["5"] = lotDetails{"Food synthesisers", 0, 3, 4, 0, 0}
	values["6"] = lotDetails{"ATVs (desert", 0, 0, 0, 2, 0}
	key := fmt.Sprintf("%v", dp.Sum("1d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollIndustrial(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["4"] = lotDetails{"Weapons/ammunition", 2, 8, 4, 2, 0}
	values["5"] = lotDetails{"Semi-finished metal products", 0, 0, 2, 0, 0}
	values["6"] = lotDetails{"Fusion power plants", 0, 3, 3, 5, 0}
	values["7"] = lotDetails{"Vehicle (grav)", 0, 2, 3, 2, 0}
	values["8"] = lotDetails{"Vehicle (air)", 0, 2, 3, 0, 0}
	values["9"] = lotDetails{"Vehicle (water)", 0, 2, 3, 0, 0}
	values["10"] = lotDetails{"Chemicals", 5, 5, 4, 2, 3}
	values["11"] = lotDetails{"Vehicle (ground)", 0, 2, 3, 0, 0}
	values["12"] = lotDetails{"Plastics", 0, 4, 2, 0, 2}
	values["13"] = lotDetails{"Computer/electronics", 0, 2, 5, 0, 0}
	values["14"] = lotDetails{"Mining/farm/building equipment", 0, 3, 3, 0, 0}
	values["15"] = lotDetails{"Consumer goods", 0, 4, 4, 0, 2}
	values["16"] = lotDetails{"Machinery/tools", 2, 3, 3, 0, 0}
	values["17"] = lotDetails{"Clothing", 0, 5, 2, 0, 4}
	values["18"] = lotDetails{"Polymers", 2, 6, 3, 0, 3}
	values["19"] = lotDetails{"Petrochemicals", 3, 10, 3, 0, 0}
	values["20"] = lotDetails{"Medical supplies", 3, 5, 5, 0, 5}
	values["21"] = lotDetails{"Special alloys", 2, 2, 2, 2, 2}
	values["22"] = lotDetails{"Grav components", 0, 2, 4, 0, 0}
	values["23"] = lotDetails{"Cybernetics", 0, 3, 5, 0, 0}
	values["24"] = lotDetails{"Prosthetics", 0, 3, 5, 0, 0}
	key := fmt.Sprintf("%v", dp.Sum("4d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollInformation(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["4"] = lotDetails{"Exotic art forms", 2, 6, 7, 0, 3}
	values["5"] = lotDetails{"Sculpture", 0, 4, 8, 0, 2}
	values["6"] = lotDetails{"Paintings", 0, 7, 7, 0, 3}
	values["7"] = lotDetails{"Writings (paper)", 0, 8, 3, 0, 2}
	values["8"] = lotDetails{"Writings (data)", 0, 3, 2, 0, 0}
	values["9"] = lotDetails{"Still pictures", 0, 2, 3, 0, 0}
	values["10"] = lotDetails{"Software (robot)", 0, 2, 3, 0, 0}
	values["11"] = lotDetails{"Software (starship)", 0, 2, 3, 0, 0}
	values["12"] = lotDetails{"Software (computer)", 0, 2, 3, 0, 0}
	values["13"] = lotDetails{"Still holo pictures", 0, 2, 2, 0, 0}
	values["14"] = lotDetails{"Audio recordings", 0, 4, 3, 0, 2}
	values["15"] = lotDetails{"Video recordings", 0, 5, 3, 0, 2}
	values["16"] = lotDetails{"Holo recordings", 0, 3, 2, 0, 0}
	values["17"] = lotDetails{"Holo recordings", 0, 3, 2, 0, 0}
	values["18"] = lotDetails{"Records (data)", 0, 3, 2, 0, 0}
	values["19"] = lotDetails{"Records (paper)", 0, 7, 3, 0, 2}
	values["20"] = lotDetails{"Raw data (data)", 0, 3, 2, 0, 0}
	values["21"] = lotDetails{"Raw data (paper)", 0, 7, 3, 0, 2}
	values["22"] = lotDetails{"Credit (data)", 0, 2, 2, 0, 0}
	values["23"] = lotDetails{"Currency", 0, 3, 3, 0, 2}
	values["24"] = lotDetails{"Erotica", 0, 4, 3, 0, 0}
	key := fmt.Sprintf("%v", dp.Sum("4d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollProcessed(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["11"] = lotDetails{"Iron", 0, 0, 0, 0, 0}
	values["12"] = lotDetails{"Steel", 0, 0, 0, 0, 0}
	values["13"] = lotDetails{"Aluminium", 0, 0, 0, 3, 0}
	values["14"] = lotDetails{"Copper", 0, 0, 0, 3, 0}
	values["15"] = lotDetails{"Tin", 0, 0, 0, 3, 0}
	values["16"] = lotDetails{"Zinc", 0, 0, 0, 3, 0}
	values["21"] = lotDetails{"Special alloys", 2, 2, 3, 3, 2}
	values["22"] = lotDetails{"Precious metals", 2, 2, 4, 0, 2}
	values["23"] = lotDetails{"Processed radioactives", 4, 3, 3, 12, 2}
	values["24"] = lotDetails{"Plant compounds", 3, 4, 2, 2, 5}
	values["25"] = lotDetails{"Animal compounds", 3, 3, 2, 2, 5}
	values["26"] = lotDetails{"Petrochemicals", 3, 9, 0, 2, 2}
	values["31"] = lotDetails{"Textiles", 2, 6, 3, 0, 2}
	values["32"] = lotDetails{"Explosives", 3, 11, 5, 2, 3}
	values["33"] = lotDetails{"Polymers", 3, 7, 0, 0, 3}
	values["34"] = lotDetails{"Paper", 0, 9, 2, 0, 4}
	values["35"] = lotDetails{"Pharmaceuticals", 4, 5, 5, 2, 5}
	values["36"] = lotDetails{"Preserved foods", 0, 3, 3, 0, 4}
	values["41"] = lotDetails{"Spices", 2, 4, 4, 0, 5}
	values["42"] = lotDetails{"Gourmet foods", 0, 3, 3, 0, 6}
	values["43"] = lotDetails{"Alcoholic beverages", 3, 6, 5, 0, 4}
	values["44"] = lotDetails{"Milks", 2, 2, 4, 0, 5}
	values["45"] = lotDetails{"Nectars", 2, 2, 4, 0, 4}
	values["46"] = lotDetails{"Syrups", 2, 2, 4, 0, 3}
	values["51"] = lotDetails{"Teas", 2, 2, 4, 0, 3}
	values["52"] = lotDetails{"Exotic fluids", 4, 3, 4, 2, 4}
	values["53"] = lotDetails{"Aromatics", 3, 4, 4, 0, 5}
	values["54"] = lotDetails{"Disposables", 0, 4, 2, 3, 0}
	values["55"] = lotDetails{"Protective gear", 0, 3, 2, 0, 0}
	values["56"] = lotDetails{"Metal parts", 0, 0, 2, 2, 0}
	values["61"] = lotDetails{"Electronic parts", 0, 3, 3, 2, 2}
	values["62"] = lotDetails{"High Tech", 0, 3, 4, 2, 3}
	values["63"] = lotDetails{"Tools", 0, 2, 2, 2, 0}
	values["64"] = lotDetails{"Entertainment equipment", 0, 4, 4, 2, 0}
	values["65"] = lotDetails{"Appliances", 0, 5, 3, 2, 0}
	values["66"] = lotDetails{"Furniture", 0, 5, 5, 0, 2}
	key := fmt.Sprintf("%v", dp.Sum("6d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollResources(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["3"] = lotDetails{"Processed radioactives", 4, 2, 3, 12, 0}
	values["4"] = lotDetails{"Ore (radioactives)", 3, 0, 2, 12, 0}
	values["5"] = lotDetails{"Crystals", 3, 2, 5, 2, 2}
	values["6"] = lotDetails{"Refined hydrocarbons", 2, 9, 0, 0, 0}
	values["7"] = lotDetails{"Refined ferrous metals", 0, 0, 2, 0, 0}
	values["8"] = lotDetails{"Refined ferrous metals", 0, 0, 2, 0, 0}
	values["9"] = lotDetails{"Refined non-ferrous metals", 2, 0, 2, 2, 0}
	values["10"] = lotDetails{"Refined non-metallics", 3, 3, 3, 2, 0}
	values["11"] = lotDetails{"Ore (ferrous metal)", 0, 0, 0, 0, 0}
	values["12"] = lotDetails{"Ore (ferrous metal)", 0, 0, 0, 0, 0}
	values["13"] = lotDetails{"Ore (non-ferrous metal)", 0, 0, 0, 3, 0}
	values["14"] = lotDetails{"Ore (non-metallic)", 3, 2, 0, 2, 0}
	values["15"] = lotDetails{"Nitrates (fertiliser)", 0, 7, 3, 0, 4}
	values["16"] = lotDetails{"Nitrates (explosive)", 2, 11, 5, 0, 3}
	values["17"] = lotDetails{"Refined precious metals", 0, 0, 3, 2, 0}
	values["18"] = lotDetails{"Refined rare earths", 2, 2, 3, 3, 0}

	key := fmt.Sprintf("%v", dp.Sum("3d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollWater(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["2"] = lotDetails{"Submarine", 0, 2, 2, 2, 0}
	values["3"] = lotDetails{"Domestic marine", 2, 0, 4, 0, 3}
	values["4"] = lotDetails{"Live seafood", 3, 0, 4, 0, 3}
	values["5"] = lotDetails{"Artificial gills", 0, 2, 3, 0, 0}
	values["6"] = lotDetails{"Refined light", 2, 2, 0, 2, 0}
	values["7"] = lotDetails{"Organic chemicals", 5, 6, 3, 0, 3}
	values["8"] = lotDetails{"Frozen seafood", 0, 0, 2, 0, 4}
	values["9"] = lotDetails{"Protein concentrate", 0, 3, 2, 0, 4}
	values["10"] = lotDetails{"Precious metals", 0, 0, 0, 2, 0}
	values["11"] = lotDetails{"Pharmaceuticals", 3, 4, 5, 2, 5}
	values["12"] = lotDetails{"Hovercraft", 0, 2, 2, 2, 0}
	key := fmt.Sprintf("%v", dp.Sum("2d6"))
	details := values[key]
	l.setDetails(details, dp)
}

func (l *Lot) rollVacuum(dp *dice.Dicepool) {
	values := map[string]lotDetails{}
	values["2"] = lotDetails{"Small spacecraft", 0, 3, 3, 2, 0}
	values["3"] = lotDetails{"Explosives", 4, 11, 7, 0, 3}
	values["4"] = lotDetails{"Frozen gasses", 3, 4, 2, 2, 0}
	values["5"] = lotDetails{"Radioactive ores", 3, 3, 0, 12, 2}
	values["6"] = lotDetails{"Non-metallic ores", 2, 2, 0, 2, 2}
	values["7"] = lotDetails{"Non-metallic ores", 2, 2, 0, 2, 2}
	values["8"] = lotDetails{"Ferrous ores", 0, 0, 0, 2, 0}
	values["9"] = lotDetails{"Ferrous ores", 0, 0, 0, 2, 0}
	values["10"] = lotDetails{"Ferrous ores", 0, 0, 0, 2, 0}
	values["11"] = lotDetails{"Non-ferrous metal ores", 2, 0, 0, 2, 0}
	values["12"] = lotDetails{"Non-ferrous metal ores", 2, 0, 0, 2, 0}
	values["13"] = lotDetails{"Vehicle (grav)", 0, 3, 2, 0, 0}
	values["14"] = lotDetails{"Vac suits", 0, 3, 3, 0, 0}
	values["15"] = lotDetails{"Pressure tents", 0, 3, 2, 0, 0}
	values["16"] = lotDetails{"Mining lasers", 0, 2, 3, 0, 0}
	values["17"] = lotDetails{"Vacuum processed parts", 0, 2, 3, 0, 0}
	values["18"] = lotDetails{"Vacuum processed chemicals", 4, 4, 4, 3, 2}

	key := fmt.Sprintf("%v", dp.Sum("3d6"))
	details := values[key]
	l.setDetails(details, dp)
}
