package statblock

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/core/entities/dice"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
	"github.com/charmbracelet/lipgloss"
)

type CharacterSheet struct {
	PersonalDataFile     PersonalData             `json:"personal_data_file"`
	CoreCharacteristics  CoreCharacteristics      `json:"core_characteristics"`
	OtherCharacteristics SecondaryCharacteristics `json:"other_characteristics"`
	Careers              []CareerTerm             `json:"careers,omitempty"`
	Skills               SkillSummary             `json:"skills"`
	Finances             Finances                 `json:"finances"`
	Armour               []Armor                  `json:"armour,omitempty"`
	Weapons              []Weapon                 `json:"weapons,omitempty"`
	Augments             []Augment                `json:"augments,omitempty"`
	Equipment            EquipmentSummary         `json:"equipment"`
	BackgrondNotes       BackgrondNotes           `json:"backgrond_notes,omitempty"`
	Allies               []Connection             `json:"allies,omitempty"`
	Contacts             []Connection             `json:"contacts,omitempty"`
	Rivals               []Connection             `json:"rivals,omitempty"`
	Enemies              []Connection             `json:"enemies,omitempty"`
	Wounds               []Wound                  `json:"wounds,omitempty"`
	Biography            string                   `json:"biography,omitempty"`
}

type PersonalData struct {
	Name      string   `json:"name,omitempty"`
	Age       int      `json:"age,omitempty"`
	Species   string   `json:"species,omitempty"`
	Traits    []string `json:"traits,omitempty"`
	Homeworld string   `json:"homeworld,omitempty"`
	Rads      int      `json:"rads,omitempty"`
}

type CoreCharacteristics struct {
	Strenght       int `json:"strenght"`
	Dexterity      int `json:"dexterity"`
	Endurance      int `json:"endurance"`
	Inteligence    int `json:"inteligence"`
	Education      int `json:"education"`
	SocialStanding int `json:"social_standing"`
}

type SecondaryCharacteristics struct {
	Psionic   int `json:"psionic,omitempty"`
	Morale    int `json:"morale,omitempty"`
	Luck      int `json:"luck,omitempty"`
	Sanity    int `json:"sanity,omitempty"`
	Charm     int `json:"charm,omitempty"`
	Territory int `json:"territory,omitempty"`
}

type CareerTerm struct {
	Career string `json:"career"`
	Terms  int    `json:"terms"`
	Rank   int    `json:"rank"`
}

type SkillSummary struct {
	SkillInTraining          skill.Skill  `json:"skill_in_training"`
	TrainingPeriodsCompleted int          `json:"training_periods_completed,omitempty"`
	TrainingPeriodsRequired  int          `json:"training_periods_required,omitempty"`
	Skills                   []SkillEntry `json:"skills,omitempty"`
}

type SkillEntry struct {
	Skill  skill.Skill `json:"skill"`
	Rating int         `json:"rating"`
}

type Finances struct {
	Pension    int `json:"pension,omitempty"`
	Debt       int `json:"debt,omitempty"`
	CashOnHand int `json:"cash_on_hand,omitempty"`
	LivingCost int `json:"living_cost,omitempty"`
}

type Armor struct {
	Type             string   `json:"type"`
	Rad              int      `json:"rad,omitempty"`
	Protection       int      `json:"protection"`
	ProtectionEnergy int      `json:"protection_energy,omitempty"`
	Mass             float64  `json:"mass"`
	Options          []string `json:"options,omitempty"`
	Equiped          bool     `json:"equiped"`
}

type Weapon struct {
	Weapon   string  `json:"weapon"`
	TL       int     `json:"tl"`
	Range    string  `json:"range"`
	Damage   string  `json:"damage"`
	Mass     float64 `json:"mass"`
	Magazine int     `json:"magazine,omitempty"`
	Equiped  bool    `json:"equiped"`
}

type Augment struct {
	Type        string `json:"type"`
	TL          int    `json:"tl"`
	Improvement string `json:"improvement,omitempty"`
}

type EquipmentSummary struct {
	EquipmentList   []Equipment `json:"equipment_list,omitempty"`
	TotalMassCaried float64     `json:"total_mass_caried,omitempty"`
}

type Equipment struct {
	Type   string  `json:"type"`
	OnSelf bool    `json:"on_self"`
	Mass   float64 `json:"mass"`
}

type BackgrondNotes []string

type Connection struct {
	Name      string `json:"name"`
	Occupancy string `json:"occupancy"`
	Relations int    `json:"relations"`
	Power     int    `json:"power"`
	Influence int    `json:"influence"`
	Note      string `json:"note"`
}

type Wound struct {
	Type            string `json:"type"`
	Location        string `json:"location"`
	RecoveryPreriod string `json:"recovery_preriod"`
	Notes           string `json:"notes"`
}

/*

STR:  8 (+0) DEX: 10 (+1) END:  8 (+0) INT: 10 (+1) EDU:  8 (+0) SOC: 10 (+1)


STR: __8_ DEX: _10_ END: __5_ INT: ____ EDU: ____ SOC: ____ |
     (+1)      (+1)      (-1)      ____      ____      ____ |



STR:  8 (+0)
DEX: 10 (+1)
END:  8 (+0)
INT: 10 (+1)
EDU:  8 (+0)
SOC: 10 (+1)


STR:  8 (+0) INT: 10 (+1)
DEX: 10 (+1) EDU:  8 (+0)
END:  8 (+0) SOC:  2 (-2)


STR:  8 INT: 10
DEX: 10 EDU:  8
END:  8 SOC:  2
*/

var Leeroy = CharacterSheet{
	PersonalDataFile: PersonalData{
		Name:      "Leeroy Jenkins",
		Age:       42,
		Species:   "Human",
		Traits:    []string{},
		Homeworld: "Drinax",
		Rads:      15,
	},
	CoreCharacteristics: CoreCharacteristics{
		Strenght:       8,
		Dexterity:      9,
		Endurance:      10,
		Inteligence:    11,
		Education:      12,
		SocialStanding: 6,
	},
	OtherCharacteristics: SecondaryCharacteristics{
		Psionic:   0,
		Morale:    11,
		Luck:      7,
		Sanity:    6,
		Charm:     8,
		Territory: 0,
	},
	Careers: []CareerTerm{
		{Career: "Rogue", Terms: 1, Rank: 1},
		{Career: "Military", Terms: 3, Rank: 2},
		{Career: "Drifter", Terms: 1, Rank: 0},
	},
	Skills: SkillSummary{
		SkillInTraining:          skill.Admin,
		TrainingPeriodsCompleted: 1,
		TrainingPeriodsRequired:  5,
		Skills: []SkillEntry{
			{Skill: skill.ScienceLife, Rating: 0},
			{Skill: skill.Pilot_SmallCraft, Rating: 1},
			{Skill: skill.Tactics_Military, Rating: 1},
			{Skill: skill.Medic, Rating: 4},
		},
	},
	Finances: Finances{
		Pension:    10000,
		Debt:       5000,
		CashOnHand: 200,
		LivingCost: 1600,
	},
	Armour: []Armor{
		{
			Type:             "Cloth",
			Rad:              0,
			Protection:       8,
			ProtectionEnergy: 8,
			Mass:             14.5,
			Options:          []string{},
			Equiped:          false,
		},
	},
	Weapons: []Weapon{
		{
			Weapon:   "Laser Pistol",
			TL:       11,
			Range:    "Long",
			Damage:   "3D+3",
			Mass:     5,
			Magazine: 100,
			Equiped:  true,
		},
	},
	Augments: []Augment{
		{
			Type:        "Wafer jack",
			TL:          12,
			Improvement: "Rating 2; Bandwith /8",
		},
	},
	Equipment:      EquipmentSummary{},
	BackgrondNotes: BackgrondNotes{},
	Allies:         []Connection{},
	Contacts:       []Connection{},
	Rivals:         []Connection{},
	Enemies: []Connection{
		{Name: "Elsa Brimor", Occupancy: "Agent", Relations: -2, Power: 1, Influence: 4, Note: "Loves Leeroy"},
	},
	Wounds: []Wound{
		{
			Type:            "Gun Shot",
			Location:        "Head",
			RecoveryPreriod: "None",
			Notes:           "",
		},
	},
	Biography: "This is Leeroys BIO",
}

func (cs CharacterSheet) View() string {
	s := renderCoreCharacteristics(cs.CoreCharacteristics)
	s2 := renderCoreCharacteristics(cs.CoreCharacteristics)
	ss := lipgloss.Place(20, 15, 10, 5, lipgloss.JoinVertical(lipgloss.Center, s, s2))
	return ss

}

func renderCoreCharacteristics(core CoreCharacteristics) string {
	s := ""
	s += fmt.Sprintf(" STR: %v \n", statValue(core.Strenght))
	s += fmt.Sprintf(" DEX: %v \n", statValue(core.Dexterity))
	s += fmt.Sprintf(" END: %v \n", statValue(core.Endurance))
	s += fmt.Sprintf(" INT: %v \n", statValue(core.Inteligence))
	s += fmt.Sprintf(" EDU: %v \n", statValue(core.Education))
	s += fmt.Sprintf(" SOC: %v ", statValue(core.SocialStanding))
	r := lipgloss.DefaultRenderer()
	st := lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:          "-",
			Bottom:       "-",
			Left:         "|",
			Right:        "|",
			TopLeft:      "1",
			TopRight:     "+",
			BottomLeft:   "+",
			BottomRight:  "+",
			MiddleLeft:   "!",
			MiddleRight:  "!",
			Middle:       "!",
			MiddleTop:    "=",
			MiddleBottom: "=",
		}, true).
		Renderer(r)
	return st.Render(s)
}

func statValue(i int) string {
	s := " "
	s += fmt.Sprintf("%d", i)
	for len(s) < 3 {
		s = " " + s
	}
	s += " "
	dm := dice.CharacteristicDM(i)
	s += "("
	if dm >= 0 {
		s += "+"
	}
	s += fmt.Sprintf("%v)", dm)
	return s
}
