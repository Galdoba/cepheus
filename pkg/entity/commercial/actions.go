package commercial

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/interaction"
)

const (
	redTapeProtection = "Red Tape Protection"
)

func (c *Corporation) SkillValue(skill string) int {
	switch skill {
	case SkillAdvocasy:
		return c.Advocasy
	case SkillAgency:
		return c.Agency
	case SkillBrokerage:
		return c.Brokerage
	case SkillFabrication:
		return c.Fabrication
	case SkillInvestiment:
		return c.Investiment
	case SkillMischief:
		return c.Mischief
	case SkillNobility:
		return c.Nobility
	case SkillPropaganda:
		return c.Propaganda
	case SkillResearch:
		return c.Research
	case SkillShipping:
		return c.Shipping
	default:
		return -3
	}
}

func (c *Corporation) CharacteristicValue(char string) int {
	switch char {
	case CharacteristicControl:
		return c.Control
	case CharacteristicGuile:
		return c.Guile
	case CharacteristicDependability:
		return c.Dependability
	case CharacteristicManagement:
		return c.Management
	default:
		panic("invalid characteristic")
	}
}

type Action struct {
	Name              string
	DurationFactor    int
	Characteristic    string
	Skill             string
	Difficulty        int
	Description       string
	RankingInfluenced bool
	// PrequisiteWealthPerRank    int
	PrequisiteEmployeesPerRank int
	PrequisiteWealthFlat       int
	// PrequisiteEmployeesFlat    int
	TargetProduction *IndustryLine
	TargetCorp       *Corporation
	KeyOptions       []string
	SelectedKey      string
	EffectFunc       func(Action, *Corporation) error
	effect           int
}

var ActionProtectCompanyThrueeLegalRedTape = Action{
	Name:                       redTapeProtection,
	DurationFactor:             1,
	Characteristic:             CharacteristicManagement,
	Skill:                      SkillAdvocasy,
	Difficulty:                 dice.CheckAverage,
	Description:                "The Effect of this skill check will reduce the Effect of the next Advocacy, Agency, Nobility or Propaganda skill roll targeting this Entity.",
	RankingInfluenced:          true,
	PrequisiteEmployeesPerRank: 10,
	PrequisiteWealthFlat:       0,
	TargetProduction:           &IndustryLine{},
	TargetCorp:                 &Corporation{},
	KeyOptions:                 []string{SkillAdvocasy, SkillAgency, SkillNobility, SkillPropaganda},
	SelectedKey:                "",
	EffectFunc:                 applyEffectRedTapeProtection,
	effect:                     0,
}

func rankDifference(c1, c2 *Corporation) int {
	d := c1.RankCurrent - c2.RankCurrent
	if d > 0 {
		d = d * -1
	}
	return d
}

func attackDM(dmKey string, c *Corporation) int {
	if val, ok := c.NextAttackDM[dmKey]; ok {
		return val
	}
	return 0
}

func (action *Action) AssesPrequisites(c *Corporation) error {
	if action.PrequisiteEmployeesPerRank > 0 && c.freeEmployees() < action.PrequisiteEmployeesPerRank*c.RankCurrent {
		return fmt.Errorf("prequisite not met: workers %v of %v", c.freeEmployees(), action.PrequisiteEmployeesPerRank*c.RankCurrent)
	}
	if action.PrequisiteWealthFlat > 0 && c.uninvestedWealth() < action.PrequisiteWealthFlat {
		return fmt.Errorf("prequisite not met: uninvested wealth %v of %v", c.uninvestedWealth(), action.PrequisiteWealthFlat)
	}

	return nil
}

func (c *Corporation) getNextUseAttackDMs(keys ...string) []int {
	dms := []int{}
	for _, key := range keys {
		if val, ok := c.NextAttackDM[key]; ok {
			dms = append(dms, val)
		} else {
			dms = append(dms, 0)
		}
	}
	return dms
}
func (c *Corporation) getNextUseDefenceDMs(keys ...string) []int {
	dms := []int{}
	for _, key := range keys {
		if val, ok := c.NextDefenceDM[key]; ok {
			dms = append(dms, val)
		} else {
			dms = append(dms, 0)
		}
	}
	return dms
}
func (c *Corporation) clearNextUseAttackDMs(keys ...string) {
	for _, key := range keys {
		delete(c.NextAttackDM, key)
	}
}
func (c *Corporation) clearNextUseDefenceDMs(keys ...string) {
	for _, key := range keys {
		delete(c.NextDefenceDM, key)
	}
}

func collectDMs(action Action, c *Corporation) []int {
	dms := []int{}
	dms = append(dms, characteristicDM(c.CharacteristicValue(action.Characteristic)))
	dms = append(dms, c.SkillValue(action.Skill))
	dms = append(dms, action.Difficulty)
	dms = append(dms, c.getNextUseAttackDMs(action.Characteristic, action.Skill)...)
	if action.TargetCorp != nil {
		dms = append(dms, action.TargetCorp.getNextUseDefenceDMs(action.Characteristic, action.Skill)...)
		if action.RankingInfluenced {
			dms = append(dms, rankDifference(c, action.TargetCorp))
		}
	}
	return dms
}

func resetDMs(action Action, c *Corporation) {
	c.clearNextUseAttackDMs(action.Characteristic, action.Skill)
	if action.TargetCorp != nil {
		action.TargetCorp.clearNextUseDefenceDMs(action.Characteristic, action.Skill)
	}
}

func (c *Corporation) Commence(action Action) error {

	if len(action.KeyOptions) > 0 {
		items := []interaction.SelectionItem{}
		for _, option := range action.KeyOptions {
			items = append(items, interaction.NewItem(option, option))
		}
		selected, err := interaction.SelectSingle("select target item",
			interaction.WithItems(items...),
			interaction.Auto(!c.humanPlayable),
		)
		if err != nil {
			return fmt.Errorf("failed to select target item: %v", err)
		}
		action.SelectedKey = selected.PayData().(string)
	}

	dms := collectDMs(action, c)
	action.effect = dice.NewDicepool().SkillCheck(dms...)

	if err := action.EffectFunc(action, c); err != nil {
		return fmt.Errorf("failed action %v: %v", action.Name, err)
	}
	fmt.Println(action)
	fmt.Println(action.effect)

	resetDMs(action, c)

	return nil
}

func applyEffectRedTapeProtection(action Action, c *Corporation) error {
	if action.effect < 1 {
		return nil
	}
	keySelected := false
	for _, option := range action.KeyOptions {
		if option == action.SelectedKey {
			keySelected = true
		}
	}
	if !keySelected {
		return fmt.Errorf("target skill not selected")
	}
	c.NextDefenceDM[action.SelectedKey] = action.effect
	return nil
}
