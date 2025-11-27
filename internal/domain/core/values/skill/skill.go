package skill

import (
	"fmt"
	"strings"
)

const (
	Admin           SkillGroup = "Admin"
	Advocate        SkillGroup = "Advocate"
	Animals         SkillGroup = "Animals"
	Art             SkillGroup = "Art"
	Astrogation     SkillGroup = "Astrogation"
	Athletics       SkillGroup = "Athletics"
	Broker          SkillGroup = "Broker"
	Carouse         SkillGroup = "Carouse"
	Deception       SkillGroup = "Deception"
	Diplomat        SkillGroup = "Diplomat"
	Drive           SkillGroup = "Drive"
	Electronics     SkillGroup = "Electronics"
	Engineer        SkillGroup = "Engineer"
	Explosives      SkillGroup = "Explosives"
	Flyer           SkillGroup = "Flyer"
	Gambler         SkillGroup = "Gambler"
	Gunner          SkillGroup = "Gunner"
	GunCombat       SkillGroup = "Gun Combat"
	HeavyWeapons    SkillGroup = "Heavy Weapons"
	Independence    SkillGroup = "Independence"
	Investigate     SkillGroup = "Investigate"
	JOAT            SkillGroup = "Jack-of-All-Trades"
	Language        SkillGroup = "Language"
	Leadership      SkillGroup = "Leadership"
	Mechanic        SkillGroup = "Mechanic"
	Medic           SkillGroup = "Medic"
	Melee           SkillGroup = "Melee"
	Navigation      SkillGroup = "Navigation"
	Persuade        SkillGroup = "Persuade"
	Pilot           SkillGroup = "Pilot"
	Profession      SkillGroup = "Profession"
	Recon           SkillGroup = "Recon"
	ScienceLife     SkillGroup = "Life Science"
	SciencePhysical SkillGroup = "Physical Science"
	ScienceRobotics SkillGroup = "Robotic Science"
	ScienceSocial   SkillGroup = "Social Science"
	ScienceSpaces   SkillGroup = "Space Science"
	Seafarer        SkillGroup = "Seafarer"
	Stealth         SkillGroup = "Stealth"
	Steward         SkillGroup = "Steward"
	Streetwise      SkillGroup = "Streetwise"
	Survival        SkillGroup = "Survival"
	Tactics         SkillGroup = "Tactics"
	Tolerance       SkillGroup = "Tolerance"
	VaccSuit        SkillGroup = "Vacc Suit"
)

type SkillGroup string

type Skill struct {
	Group      SkillGroup `json:"group"`
	Speciality string     `json:"speciality"`
}

var NoSkill = Skill{}

func (s Skill) String() string {
	str := fmt.Sprintf("%s", s.Group)
	if s.Speciality != "" {
		str += fmt.Sprintf(" (%v)", s.Speciality)
	}
	return str
}

// func groupString(group SkillGroup) string {
// 	switch group {
// 	case ScienceLife, SciencePhysical, ScienceRobotics, ScienceSocial, ScienceSpaces:
// 		return "Science"
// 	default:
// 		return fmt.Sprintf("%s", group)
// 	}
// }

func Call(key string, skills ...Skill) Skill {
	for _, s := range skills {
		if strings.Contains(key, s.Speciality) {
			return s
		}
		if key == string(s.Group) {
			return s
		}
	}
	return NoSkill
}

func New(gr SkillGroup) Skill {
	return Skill{
		Group: gr,
	}
}

func (s Skill) WithSpec(spec string) Skill {
	for _, sp := range Speciality[s.Group] {
		if sp == spec {
			s.Speciality = sp
			break
		}
	}
	return s
}

var Speciality = map[SkillGroup][]string{
	Animals:         []string{"Handling", "Veterinary", "Training"},
	Art:             []string{"Performer", "Holograpy", "Instrument", "Visual Media", "Write"},
	Athletics:       []string{"Dexterity", "Endurance", "Strength"},
	Drive:           []string{"Hovercraft", "Mole", "Track", "Walker", "Wheel"},
	Electronics:     []string{"Comms", "Computers", "Remote Ops", "Sensors"},
	Engineer:        []string{"M-drive", "J-drive", "Life Support", "Power"},
	Flyer:           []string{"Airship", "Grav", "Ornithopter", "Rotor", "Wing"},
	Gunner:          []string{"Turret", "Ortilery", "Screen", "Capital"},
	GunCombat:       []string{"Archaic", "Energy", "Slug"},
	HeavyWeapons:    []string{"Artilery", "Portable", "Vechicle"},
	Language:        []string{"Galinglic", "[ByRace]"},
	Melee:           []string{"Unarmed", "Blade", "Bludgeon", "Natural"},
	Pilot:           []string{"Small Craft", "Spaceship", "Capital Ships"},
	Profession:      []string{"[ByProfession]"},
	ScienceLife:     []string{"Biology", "Genetics", "Psionicology", "Xenology"},
	SciencePhysical: []string{"Chemistry", "Physics", "Jumpspace Physics"},
	ScienceRobotics: []string{"Cybernetics", "Robotics"},
	ScienceSocial:   []string{"Archaeology", "Economics", "History", "Linguistics", "Philosophy", "Psychology", "Sophontology"},
	ScienceSpaces:   []string{"Astronomy", "Cosmology", "Planetology"},
	Seafarer:        []string{"Ocean Ships", "Personal", "Sail", "Submarines"},
	Tactics:         []string{"Military", "Naval"},
}

func BackgroundSkillList() []Skill {
	return []Skill{
		Skill{Group: Admin},
		Skill{Group: Animals},
		Skill{Group: Art},
		Skill{Group: Athletics},
		Skill{Group: Carouse},
		Skill{Group: Drive},
		Skill{Group: Electronics},
		Skill{Group: Flyer},
		Skill{Group: Language},
		Skill{Group: Mechanic},
		Skill{Group: Medic},
		Skill{Group: Profession},
		Skill{Group: ScienceLife},
		Skill{Group: SciencePhysical},
		Skill{Group: ScienceRobotics},
		Skill{Group: ScienceSocial},
		Skill{Group: ScienceSpaces},
		Skill{Group: Streetwise},
		Skill{Group: Survival},
		Skill{Group: VaccSuit},
	}
}
