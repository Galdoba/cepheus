package characteristic

import "strings"

const (
	TypePhysical  CharacteristicType = "Physical"
	TypeMental    CharacteristicType = "Mental"
	TypeObscured  CharacteristicType = "Obscured"
	TypeCore      CharacteristicType = "Core"
	TypeSecondary CharacteristicType = "Obscured"
)

type CharacteristicType string

type Characteristic struct {
	Name  string               `json:"name"`
	Abb   string               `json:"abb"`
	Types []CharacteristicType `json:"types"`
}

func Call(key string, list ...Characteristic) Characteristic {
	key = strings.ToLower(key)
	for _, l := range list {
		if key == strings.ToLower(l.Name) || key == strings.ToLower(l.Abb) {
			return l
		}
	}
	return NONE
}

var NONE = Characteristic{
	Name:  "",
	Abb:   "",
	Types: []CharacteristicType{},
}

var Strenght = Characteristic{
	Name:  "Strength",
	Abb:   "STR",
	Types: []CharacteristicType{TypeCore, TypePhysical},
}

var Dexterity = Characteristic{
	Name:  "Dexterity",
	Abb:   "DEX",
	Types: []CharacteristicType{TypeCore, TypePhysical},
}

var Endurance = Characteristic{
	Name:  "Endurance",
	Abb:   "END",
	Types: []CharacteristicType{TypeCore, TypePhysical},
}

var Inteligence = Characteristic{
	Name:  "Intelligence",
	Abb:   "INT",
	Types: []CharacteristicType{TypeCore, TypeMental},
}

var Education = Characteristic{
	Name:  "Education",
	Abb:   "EDU",
	Types: []CharacteristicType{TypeCore, TypeMental},
}

var SocialStanding = Characteristic{
	Name:  "Social Standing",
	Abb:   "SOC",
	Types: []CharacteristicType{TypeCore, TypeMental},
}

var PsionicStrength = Characteristic{
	Name:  "Psionic Strength",
	Abb:   "PSI",
	Types: []CharacteristicType{TypeCore, TypeObscured},
}
var Wealth = Characteristic{
	Name:  "Wealth",
	Abb:   "WLT",
	Types: []CharacteristicType{TypeSecondary},
}
var Sanity = Characteristic{
	Name:  "Sanity",
	Abb:   "SAN",
	Types: []CharacteristicType{TypeCore, TypeObscured},
}
var Morale = Characteristic{
	Name:  "Morale",
	Abb:   "MOR",
	Types: []CharacteristicType{TypeCore, TypeObscured},
}
var Luck = Characteristic{
	Name:  "Luck",
	Abb:   "LCK",
	Types: []CharacteristicType{TypeCore, TypeObscured},
}
