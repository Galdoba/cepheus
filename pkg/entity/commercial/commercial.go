package commercial

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
)

const (
	SkillAdvocasy                              = "Advocasy"
	SkillAgency                                = "Agency"
	SkillBrokerage                             = "Brokerage"
	SkillFabrication                           = "Fabrication"
	SkillInvestiment                           = "Investiment"
	SkillMischief                              = "Mischief"
	SkillNobility                              = "Nobility"
	SkillPropaganda                            = "Propaganda"
	SkillResearch                              = "Research"
	SkillShipping                              = "Shipping"
	CharacteristicControl                      = "Control"
	CharacteristicDependability                = "Dependability"
	CharacteristicGuile                        = "Guide"
	CharacteristicManagement                   = "Management"
	TraitRanking                               = "Ranking"
	TraitLoyalty                               = "Loyalty"
	TraitReputation                            = "Reputation"
	TraitWealth                                = "Wealth"
	MaxValueSkill                              = 6
	MaxValueChar                               = 15
	MaxValueTrait                              = 10
	CreditsWealthRatio                         = 25000
	sellAssetsRatioBaseProfit                  = 3000
	shareProfitBase                            = 10000
	investimentsBase                           = 100000
	Low                         LivingStandard = 0
	Average                     LivingStandard = 1
	Good                        LivingStandard = 2
	High                        LivingStandard = 3
	VeryHigh                    LivingStandard = 4
	Rich                        LivingStandard = 5
)

type LivingStandard int

func (ls LivingStandard) String() string {
	switch ls {
	case Low:
		return "Low"
	case Average:
		return "Average"
	case Good:
		return "Good"
	case High:
		return "High"
	case VeryHigh:
		return "Very High"
	case Rich:
		return "Rich"
	}
	return "error"
}

type Corporation struct {
	//CommercialEntityInfo
	CompanyName          string
	MissionStatement     string
	MissionStatementType string
	//Leaders
	Leaders []*Leader
	//Characteristics
	Control       int
	Dependability int
	Guile         int
	Management    int
	//Skills
	Advocasy    int
	Agency      int
	Brokerage   int
	Fabrication int
	Investiment int
	Mischief    int
	Nobility    int
	Propaganda  int
	Research    int
	Shipping    int
	//Traits
	Loyalty    int
	Reputation int
	//Stats
	Wealth                 int
	wealthReceived         int
	Employees              int
	EmployeeLivingStandard LivingStandard
	RankCurrent            int
	RankMaximumReached     int
	IndustryLines          map[ProductionType]*IndustryLine
	profitPool             int
	//BoonsAndBanes
	NextAttackDM  map[string]int
	NextDefenceDM map[string]int
	//options
	humanPlayable bool
}

type corpCreationOptions struct {
	humanPlayable bool
	leaders       []*Leader
	mission       MissionStatement
}

type corpCreationOption func(*corpCreationOptions)

func IsHumanPlayable(hp bool) corpCreationOption {
	return func(cco *corpCreationOptions) {
		cco.humanPlayable = hp
	}
}

func NewCorporation(options ...corpCreationOption) *Corporation {
	corp := Corporation{}
	opts := corpCreationOptions{
		leaders:       []*Leader{},
		mission:       noStatement,
		humanPlayable: false,
	}
	for _, modify := range options {
		modify(&opts)
	}
	if opts.mission == noStatement {
		opts.mission = randomStatement()
	}
	if len(opts.leaders) == 0 {
		opts.leaders = randomLeaders()
	}
	corp.Leaders = opts.leaders
	corp.humanPlayable = opts.humanPlayable
	corp.NextAttackDM = make(map[string]int)
	corp.NextDefenceDM = make(map[string]int)
	//step 1
	corp.MissionStatementType = string(opts.mission.Type)
	corp.Control = opts.mission.Control
	corp.Dependability = opts.mission.Dependability
	corp.Guile = opts.mission.Guide
	corp.Management = opts.mission.Management
	//step 2
	for _, founder := range opts.leaders {
		skills := founder.ProvideSkills()
		for _, skill := range skills {
			corp.increaseSkill(skill)
		}
	}
	//step 3
	corp.Loyalty = min(corp.Control+corp.Brokerage-(characteristicDM(corp.Guile)+corp.Mischief), MaxValueTrait)
	corp.Reputation = min(corp.Dependability+corp.Propaganda-(characteristicDM(corp.Guile)+corp.Mischief), MaxValueTrait)
	corp.Wealth++
	corp.applyInitialFounderInvesiment()
	corp.applyOpeningBankStructure()
	//step 4
	corp.applyInitialHiring()
	corp.EmployeeLivingStandard = Average
	//step 5

	corp.IndustryLines = make(map[ProductionType]*IndustryLine)
	return &corp
}

func characteristicDM(val int) int {
	switch val {
	case 0:
		return -3
	case 1, 2:
		return -2
	case 3, 4, 5:
		return -1
	case 6, 7, 8:
		return 0
	case 9, 10, 11:
		return +1
	case 12, 13, 14:
		return 2
	case 15:
		return 3
	}
	panic("invalid characteristic value")
}

func (c *Corporation) increaseSkill(skill string) {
	var exceeded bool

	switch skill {
	case SkillAdvocasy:
		c.Advocasy++
		exceeded = c.Advocasy > MaxValueSkill
		if exceeded {
			c.Advocasy = MaxValueSkill
		}
	case SkillAgency:
		c.Agency++
		exceeded = c.Agency > MaxValueSkill
		if exceeded {
			c.Agency = MaxValueSkill
		}
	case SkillBrokerage:
		c.Brokerage++
		exceeded = c.Brokerage > MaxValueSkill
		if exceeded {
			c.Brokerage = MaxValueSkill
		}
	case SkillFabrication:
		c.Fabrication++
		exceeded = c.Fabrication > MaxValueSkill
		if exceeded {
			c.Fabrication = MaxValueSkill
		}
	case SkillInvestiment:
		c.Investiment++
		exceeded = c.Investiment > MaxValueSkill
		if exceeded {
			c.Investiment = MaxValueSkill
		}
	case SkillMischief:
		c.Mischief++
		exceeded = c.Mischief > MaxValueSkill
		if exceeded {
			c.Mischief = MaxValueSkill
		}
	case SkillNobility:
		c.Nobility++
		exceeded = c.Nobility > MaxValueSkill
		if exceeded {
			c.Nobility = MaxValueSkill
		}
	case SkillPropaganda:
		c.Propaganda++
		exceeded = c.Propaganda > MaxValueSkill
		if exceeded {
			c.Propaganda = MaxValueSkill
		}
	case SkillResearch:
		c.Research++
		exceeded = c.Research > MaxValueSkill
		if exceeded {
			c.Research = MaxValueSkill
		}
	case SkillShipping:
		c.Shipping++
		exceeded = c.Shipping > MaxValueSkill
		if exceeded {
			c.Shipping = MaxValueSkill
		}
	}

	if exceeded {
		c.Wealth++
	}
}

func (corp *Corporation) applyInitialFounderInvesiment() {
	investiments := 0
	shares := 0
	for _, founder := range corp.Leaders {
		fi, fs := founder.Invested()
		investiments += fi
		shares += fs
	}
	corp.Wealth += investiments / CreditsWealthRatio
	corp.Wealth += (shares * 5)
}

func (corp *Corporation) applyOpeningBankStructure() {
	corp.Wealth += corp.Management * corp.Investiment
}

func (corp *Corporation) applyInitialHiring() {
	corp.Employees = (corp.Loyalty + corp.Propaganda) * corp.Dependability
}

func (corp *Corporation) HireAdditionalEmployees(dp *dice.Dicepool, spending int) error {
	if spending >= corp.Wealth {
		return fmt.Errorf("can't spend %v Wealth for hiring", spending)
	}
	corp.Wealth -= spending
	corp.Employees += dp.Sum("1d6") * spending
	return nil
}

func (corp *Corporation) RecordSheet() []string {
	s := []string{fmt.Sprintf("Company Name: %v", corp.CompanyName)}
	s = append(s, fmt.Sprintf("Company Mission Statement: %v", corp.MissionStatement))
	s = append(s, fmt.Sprintf(""))
	s = append(s, fmt.Sprintf("Mission Statement Type: %v", corp.MissionStatementType))
	s = append(s, fmt.Sprintf(""))
	s = append(s, fmt.Sprintf("Company Leader(s):"))
	for _, leader := range corp.Leaders {
		s = append(s, "  "+leader.String())
	}
	s = append(s, fmt.Sprintf(""))
	s = append(s, fmt.Sprintf("Entity Characteristics"))
	s = append(s, fmt.Sprintf("Control (Con)      :  %v (%v)", corp.Control, characteristicDM(corp.Control)))
	s = append(s, fmt.Sprintf("Dependability (Dep):  %v (%v)", corp.Dependability, characteristicDM(corp.Dependability)))
	s = append(s, fmt.Sprintf("Guile (Gle)        :  %v (%v)", corp.Guile, characteristicDM(corp.Guile)))
	s = append(s, fmt.Sprintf("Management (Mng)   :  %v (%v)", corp.Management, characteristicDM(corp.Management)))

	s = append(s, fmt.Sprintf(""))
	s = append(s, fmt.Sprintf("Entity Skills"))
	s = append(s, fmt.Sprintf("Advocasy   :   %v", corp.Advocasy))
	s = append(s, fmt.Sprintf("Agency     :   %v", corp.Agency))
	s = append(s, fmt.Sprintf("Brokerage  :   %v", corp.Brokerage))
	s = append(s, fmt.Sprintf("Fabrication:   %v", corp.Fabrication))
	s = append(s, fmt.Sprintf("Investiment:   %v", corp.Investiment))
	s = append(s, fmt.Sprintf("Mischief   :   %v", corp.Mischief))
	s = append(s, fmt.Sprintf("Nobility   :   %v", corp.Nobility))
	s = append(s, fmt.Sprintf("Propaganda :   %v", corp.Propaganda))
	s = append(s, fmt.Sprintf("Research   :   %v", corp.Research))
	s = append(s, fmt.Sprintf("Shipping   :   %v", corp.Shipping))

	s = append(s, fmt.Sprintf(""))
	s = append(s, fmt.Sprintf("Entity Traits"))
	s = append(s, fmt.Sprintf("Company Ranking:   %v", corp.RankCurrent))
	s = append(s, fmt.Sprintf("Loyalty        :   %v", corp.Loyalty))
	s = append(s, fmt.Sprintf("Reputation     :   %v", corp.Reputation))
	s = append(s, fmt.Sprintf("Wealth         :   %v/%v", corp.uninvestedWealth(), corp.Wealth))
	s = append(s, fmt.Sprintf("Employee Pool  :   %v/%v", corp.freeEmployees(), corp.Employees))

	s = append(s, fmt.Sprintf(""))
	s = append(s, fmt.Sprintf("Industry Line Invested"))
	for _, productionType := range listLines() {
		if line, ok := corp.IndustryLines[productionType]; ok {
			s = append(s, line.String())
		}
	}

	return s
}

func (c *Corporation) CalculateProfit() {
	wealthCalculated := 0
	positiveMult := []int{1, 5, 12, 25, 50, 250, 500}
	negativeMult := []int{2, 12, 25, 50, 250, 500, 1000}
	switch c.profitPool > 0 {
	case true:
		wealthCalculated = c.profitPool * positiveMult[c.RankCurrent]
	case false:
		wealthCalculated = c.profitPool * negativeMult[c.RankCurrent]
	}
	c.Wealth += wealthCalculated
	c.wealthReceived += wealthCalculated
	fmt.Println("wealth gained", c.wealthReceived)
}

func (c *Corporation) uninvestedWealth() int {
	wealth := c.Wealth
	for _, industry := range c.IndustryLines {
		wealth -= industry.WealthInvested
	}
	return wealth
}

func (c *Corporation) freeEmployees() int {
	workers := c.Employees
	for _, industry := range c.IndustryLines {
		workers -= industry.WorkersProvided
	}
	return workers
}

func (c *Corporation) PayStaff() {
	wealthToEmployeeRatio := [][]int{}
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{50, 35, 25, 20, 15, 5})
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{100, 75, 60, 40, 25, 10})
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{250, 125, 100, 65, 50, 25})
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{500, 350, 250, 125, 75, 50})
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{1000, 500, 400, 300, 250, 100})
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{5000, 1250, 1000, 750, 500, 250})
	wealthToEmployeeRatio = append(wealthToEmployeeRatio, []int{10000, 7500, 5000, 2500, 1000, 500})
	ratio := wealthToEmployeeRatio[c.RankCurrent][int(c.EmployeeLivingStandard)]
	paymentRequirement := (c.Employees / ratio) + 1
	c.Wealth -= paymentRequirement
	fmt.Println("wages spend", paymentRequirement)
}

func (c *Corporation) LiquefyAssets(sellWealth int) (int, error) {
	if sellWealth >= c.Wealth {
		return 0, fmt.Errorf("entity has no %v wealth", sellWealth)
	}
	assetsGained := dice.FastRandom("1d6") * sellAssetsRatioBaseProfit * sellWealth
	c.Wealth -= sellWealth
	fmt.Println("wealth sold", sellWealth)
	return assetsGained, nil
}

func (c *Corporation) ProfitShare(sellWealth int) (int, error) {
	profitTotal := shareProfitBase * c.wealthReceived
	proffitLiquified, err := c.LiquefyAssets(sellWealth)
	if err != nil {
		return 0, err
	}
	profitTotal += proffitLiquified
	profitMult := 1.0
	for _, leader := range c.Leaders {
		if leader.Founder {
			continue
		}
		profitMult -= float64(sumCareers(leader.Career)*2) / 100.0
	}
	fmt.Println("mult", profitMult)
	profit := max(0, int(float64(profitTotal)*profitMult))

	share := profit / len(c.Leaders)
	return share, nil
}

func (c *Corporation) HireLeader() {
	newLeader := NewRandomLeader(false)
	c.Leaders = append(c.Leaders, newLeader)
}

func (c *Corporation) AddInvestor() {
	newLeader := NewRandomLeader(false)
	newLeader.Career = make(map[string]int)
	investiments := (dice.FastRandom("2d6") / 2) * investimentsBase
	newLeader.WithInvestiments(investiments)
	c.Wealth += investiments / CreditsWealthRatio
	c.Leaders = append(c.Leaders, newLeader)
}
