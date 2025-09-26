package commercial

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
)

const (
	Processing                         ProductionType = "Processing"
	AutomatedMining                    ProductionType = "Automated Mining"
	Manufacturing                      ProductionType = "Manufacturing"
	Machining                          ProductionType = "Machining"
	ForestryAndFarming                 ProductionType = "Forestry and Farming"
	Motorworks                         ProductionType = "Motorworks"
	CivilEngineering                   ProductionType = "Civil Engineering"
	ArmouryTrade                       ProductionType = "Armoury Trade"
	ChemicalEngineering                ProductionType = "Chemical Engineering"
	MedicalEngineering                 ProductionType = "Medical Engineering"
	MechanicalEngineering              ProductionType = "Mechanical Engineering"
	Refinement                         ProductionType = "Refinement"
	MilitaryEngineering                ProductionType = "Military Engineering"
	AdvancedRobotics                   ProductionType = "Advanced Robotics"
	TechnicalMining                    ProductionType = "Technical Mining"
	AeronauticsAndGravimetrics         ProductionType = "Aeronautics and Gravimetrics"
	WorkforceManagement                ProductionType = "Workforce Management"
	CargoLineTransit                   ProductionType = "Cargo Line Transit"
	JunkyardManagement                 ProductionType = "Junkyard Management"
	ProfessionalContracting            ProductionType = "Professional Contracting"
	CommercialTransit                  ProductionType = "Commercial Transit"
	QualityControl                     ProductionType = "Quality Control"
	BrainTrust                         ProductionType = "Brain Trust"
	Banking                            ProductionType = "Banking"
	ArtisticCreation                   ProductionType = "Artistic Creation"
	ConlgomerationControl              ProductionType = "Conlgomeration Control"
	notFound                           ProductionType = "not found"
	staffRequirementMultiplier                        = 10
	marketExpertsRequirementMultiplier                = 10
)

func listLines() []ProductionType {
	return []ProductionType{
		Processing,
		AutomatedMining,
		Manufacturing,
		Machining,
		ForestryAndFarming,
		Motorworks,
		CivilEngineering,
		ArmouryTrade,
		ChemicalEngineering,
		MedicalEngineering,
		MechanicalEngineering,
		Refinement,
		MilitaryEngineering,
		AdvancedRobotics,
		TechnicalMining,
		AeronauticsAndGravimetrics,
		WorkforceManagement,
		CargoLineTransit,
		JunkyardManagement,
		ProfessionalContracting,
		CommercialTransit,
		QualityControl,
		BrainTrust,
		Banking,
		ArtisticCreation,
		ConlgomerationControl,
	}
}

type ProductionType string

type IndustryLine struct {
	Type                     ProductionType
	PossibleGoodsProduced    string
	ServicesOffered          string
	SkillUsed                string
	InvestimentRate          int
	StaffRequirementFactor   float64
	ProfitModifier           int
	WealthInvested           int
	WorkersProvided          int
	LinesOperational         int
	QuaterlyProductionResult int
}

func NewProductionLine(productionType ProductionType) *IndustryLine {
	il := IndustryLine{}
	switch productionType {
	case Processing:
		il.Type = productionType
		il.PossibleGoodsProduced = "Basic Raw Materials, Polymers, Textiles"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 10
		il.ProfitModifier = -1
	case AutomatedMining:
		il.Type = productionType
		il.PossibleGoodsProduced = "Basic Ore"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 10
		il.ProfitModifier = -1
	case Manufacturing:
		il.Type = productionType
		il.PossibleGoodsProduced = "Basic Consumables, Basic Electronics, Basic Manufactured Goods"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 15
		il.ProfitModifier = 0
	case Machining:
		il.Type = productionType
		il.PossibleGoodsProduced = "Basic Machine Parts"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 15
		il.ProfitModifier = 0
	case ForestryAndFarming:
		il.Type = productionType
		il.PossibleGoodsProduced = "Live Animals, Wood"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 15
		il.ProfitModifier = 0
	case Motorworks:
		il.Type = productionType
		il.PossibleGoodsProduced = "Vechicles"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 25
		il.ProfitModifier = 0
	case CivilEngineering:
		il.Type = productionType
		il.PossibleGoodsProduced = "Advanced Electronics, Advanced Manufactured Goods, Luxury Consumables, Luxury Goods, Spices"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 30
		il.ProfitModifier = 0
	case ArmouryTrade:
		il.Type = productionType
		il.PossibleGoodsProduced = "Armor, Weapons"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 50
		il.ProfitModifier = 1
	case ChemicalEngineering:
		il.Type = productionType
		il.PossibleGoodsProduced = "Biochemicals, Petrochemicals"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 50
		il.ProfitModifier = 1
	case MedicalEngineering:
		il.Type = productionType
		il.PossibleGoodsProduced = "Medical Supplies, Pharmaceuticals"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 50
		il.ProfitModifier = 1
	case MechanicalEngineering:
		il.Type = productionType
		il.PossibleGoodsProduced = "Advanced Machine Parts"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 50
		il.ProfitModifier = 1
	case Refinement:
		il.Type = productionType
		il.PossibleGoodsProduced = "Precious Metals, Uncommon Raw Materials"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 75
		il.ProfitModifier = 1
	case MilitaryEngineering:
		il.Type = productionType
		il.PossibleGoodsProduced = "Advanced Armor, Advanced Weapons"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 100
		il.ProfitModifier = 2
	case AdvancedRobotics:
		il.Type = productionType
		il.PossibleGoodsProduced = "Cybernetics, Robots"
		il.SkillUsed = SkillResearch
		il.InvestimentRate = 100
		il.ProfitModifier = 2
	case TechnicalMining:
		il.Type = productionType
		il.PossibleGoodsProduced = "Cristals & Gems, Radioactives, Uncommon Ores"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 100
		il.ProfitModifier = 2
	case AeronauticsAndGravimetrics:
		il.Type = productionType
		il.PossibleGoodsProduced = "Advances Vechicles"
		il.SkillUsed = SkillFabrication
		il.InvestimentRate = 200
		il.ProfitModifier = 3
	case WorkforceManagement:
		il.Type = productionType
		il.ServicesOffered = "Unskilled Labour"
		il.SkillUsed = SkillBrokerage
		il.InvestimentRate = 10
		il.ProfitModifier = -1
	case CargoLineTransit:
		il.Type = productionType
		il.ServicesOffered = "Freight and mail Transport"
		il.SkillUsed = SkillShipping
		il.InvestimentRate = 25
		il.ProfitModifier = 0
	case JunkyardManagement:
		il.Type = productionType
		il.ServicesOffered = "Junk Dealing, Salvage"
		il.SkillUsed = SkillShipping
		il.InvestimentRate = 25
		il.ProfitModifier = 0
	case ProfessionalContracting:
		il.Type = productionType
		il.ServicesOffered = "Agent, Mercenary or Merchant Brokerage"
		il.SkillUsed = SkillAgency
		il.InvestimentRate = 25
		il.ProfitModifier = 0
	case CommercialTransit:
		il.Type = productionType
		il.ServicesOffered = "Passenger Transport"
		il.SkillUsed = SkillShipping
		il.InvestimentRate = 50
		il.ProfitModifier = 1
	case QualityControl:
		il.Type = productionType
		il.ServicesOffered = "Prototype Testing"
		il.SkillUsed = SkillAgency
		il.InvestimentRate = 50
		il.ProfitModifier = 1
	case BrainTrust:
		il.Type = productionType
		il.ServicesOffered = "Research and Development"
		il.SkillUsed = SkillResearch
		il.InvestimentRate = 100
		il.ProfitModifier = 2
	case Banking:
		il.Type = productionType
		il.ServicesOffered = "Investiments and Stock Brokering"
		il.SkillUsed = SkillInvestiment
		il.InvestimentRate = 100
		il.ProfitModifier = 2
	case ArtisticCreation:
		il.Type = productionType
		il.ServicesOffered = "Exotic Refinement"
		il.SkillUsed = SkillNobility
		il.InvestimentRate = 200
		il.ProfitModifier = 3
	case ConlgomerationControl:
		il.Type = productionType
		il.ServicesOffered = "Management of Commercial Entities"
		il.SkillUsed = SkillBrokerage
		il.InvestimentRate = 500
		il.ProfitModifier = 5
	default:
		panic("unknown industry type")
	}
	il.StaffRequirementFactor = 1.0
	return &il
}

func (il *IndustryLine) String() string {
	s := string(il.Type)
	for len(s) < len(string(AeronauticsAndGravimetrics)) {
		s += " "
	}
	s += "\t" + fmt.Sprintf("(x%v)\t", il.LinesOperational)
	sup := il.SkillUsed
	for len(sup) < len(SkillInvestiment) {
		sup += " "
	}
	s += sup + "\t"
	pms := ""
	switch il.ProfitModifier {
	case 0:
		pms = "--"
	case 1, 2, 3, 4, 5:
		pms = fmt.Sprintf("+%v", il.ProfitModifier)
	case -1:
		pms = fmt.Sprintf("%v", il.ProfitModifier)
	}
	s += pms + "\t"
	s += fmt.Sprintf("%v\t", il.WealthInvested*il.LinesOperational)
	s += fmt.Sprintf("%v\t", il.QuaterlyProductionResult)
	return s
}

func (c *Corporation) StartIndustryLines(production ProductionType, lines int) error {
	if c.IndustryLines == nil {
		return fmt.Errorf("industry lines are not initiated")
	}
	industryLine := &IndustryLine{Type: notFound}
	for _, line := range c.IndustryLines {
		if line.Type == production {
			industryLine = line
			break
		}
	}
	if industryLine.Type == notFound {
		industryLine = NewProductionLine(production)
	}

	staffRequired := max(c.RankCurrent, 1) * staffRequirementMultiplier * lines
	if staffRequired >= c.freeEmployees() {
		return fmt.Errorf("not enough free employees to start %v Industry Line: %v of %v available", string(production), c.freeEmployees(), staffRequired)
	}

	wealthRequired := industryLine.InvestimentRate * max(c.RankCurrent, 1) * lines
	if wealthRequired >= c.uninvestedWealth() {
		return fmt.Errorf("not enough uninvested Wealth to start %v Industry Line: %v of %v available", string(production), c.uninvestedWealth(), wealthRequired)
	}

	industryLine.WorkersProvided += staffRequired
	industryLine.WealthInvested += wealthRequired
	industryLine.LinesOperational += lines
	c.IndustryLines[industryLine.Type] = industryLine
	fmt.Println(c.IndustryLines[industryLine.Type].LinesOperational)
	return nil
}

func (c *Corporation) WorkOnIndustries() {
	profitPool := 0
	dp := dice.NewDicepool()
	for _, industry := range c.IndustryLines {
		charDM := 0
		switch industry.ServicesOffered == "" {
		case true:
			charDM = characteristicDM(c.Control)
		case false:
			charDM = characteristicDM(c.Management)
		}
		roll := dp.SkillCheck(dice.CheckAverage, charDM, c.SkillValue(industry.SkillUsed))
		industry.QuaterlyProductionResult = roll
		profitPool += roll * industry.LinesOperational

	}
	c.profitPool = profitPool
}

func (c *Corporation) MegaTrade(wealthRisked int) error {
	if wealthRisked >= c.uninvestedWealth() {
		return fmt.Errorf("wealth risked exeeds or equal uninvested wealth: %v >= %v", wealthRisked, c.uninvestedWealth())
	}
	investimentSkillValue := c.SkillValue(SkillInvestiment)
	if investimentSkillValue < 1 {
		return fmt.Errorf("invesiment expetise is absent")
	}
	workersRequired := investimentSkillValue * marketExpertsRequirementMultiplier
	if workersRequired >= c.freeEmployees() {
		return fmt.Errorf("not enough workers to send to market: %v of %v available", c.freeEmployees(), workersRequired)
	}
	roll := minmax(dice.NewDicepool().SkillCheck(dice.CheckAverage, characteristicDM(c.Guile), investimentSkillValue), -6, 6)
	mult := 1.0
	switch roll {
	case -6:
		mult = 0.0
	case -5:
		mult = 0.2
	case -4:
		mult = 0.4
	case -3:
		mult = 0.6
	case -2:
		mult = 0.8
	case -1:
		mult = 0.9
	case 0:
		mult = 0.95
	case 1:
		mult = 1.0
	case 2:
		mult = 1.1
	case 3:
		mult = 1.25
	case 4:
		mult = 1.5
	case 5:
		mult = 1.75
	case 6:
		mult = 2.0
	}
	wealthRewarded := int(float64(wealthRisked)*mult) + 1
	c.Wealth = c.Wealth - wealthRisked + wealthRewarded
	fmt.Println("risked", wealthRisked, "rewarded", wealthRewarded)
	c.wealthReceived += (wealthRewarded - wealthRisked)
	return nil
}

func minmax(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
