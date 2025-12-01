package statblock

type CharacterSheet struct {
	PersonalDataFile     PersonalData             `json:"personal_data_file"`
	CoreCharacteristics  CoreCharacteristics      `json:"core_characteristics"`
	OtherCharacteristics SecondaryCharacteristics `json:"other_characteristics"`
	Careers              CareerSummary            `json:"careers,omitempty"`
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
	Name      string        `json:"name,omitempty"`
	Age       int           `json:"age,omitempty"`
	Species   string        `json:"species,omitempty"`
	Traits    []string      `json:"traits,omitempty"`
	Homeworld string        `json:"homeworld,omitempty"`
	Rads      int           `json:"rads,omitempty"`
	Careers   CareerSummary `json:"careers,omitempty"`
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

type CareerSummary []CareerTerm

type CareerTerm struct {
	Career string `json:"career,omitempty"`
	Terms  string `json:"terms,omitempty"`
	Rank   int    `json:"rank,omitempty"`
}

type SkillSummary struct {
	SkillInTraining          string       `json:"skill_in_training"`
	TrainingPeriodsCompleted int          `json:"training_periods_completed,omitempty"`
	TrainingPeriodsRequired  int          `json:"training_periods_required,omitempty"`
	Skills                   []SkillEntry `json:"skills,omitempty"`
}

type SkillEntry struct {
	Skill  string `json:"skill"`
	Rating int    `json:"rating"`
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
