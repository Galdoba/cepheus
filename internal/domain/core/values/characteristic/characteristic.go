package characteristic

import "strings"

const (
	TypePhysical      CharacteristicType = "Physical"
	TypeMental        CharacteristicType = "Mental"
	TypeObscured      CharacteristicType = "Obscured"
	TypeCore          CharacteristicType = "Core"
	TypeSecondary     CharacteristicType = "Obscured"
	NO_CHARACTERISTIC Characteristic     = "[NO CHARACTERISTIC]"
	Strength          Characteristic     = "Strength"
	Dexterity         Characteristic     = "Dexterity"
	Endurance         Characteristic     = "Endurance"
	Inteligence       Characteristic     = "Inteligence"
	Education         Characteristic     = "Education"
	SocialStanding    Characteristic     = "Social Standing"
	PsionicStrength   Characteristic     = "Psionic Strength"
	Wealth            Characteristic     = "Wealth"
	Sanity            Characteristic     = "Sanity"
	Morale            Characteristic     = "Morale"
	Luck              Characteristic     = "Luck"
	Territory         Characteristic     = "Territory"
	Charisma          Characteristic     = "Charisma"
)

type CharacteristicType string

type Characteristic string

func Call(key string, list ...Characteristic) Characteristic {
	for _, l := range list {
		name := string(l)
		if strings.EqualFold(key, name) {
			return l
		}
		if strings.EqualFold(key, Abbreviation[l]) {
			return l
		}
	}
	return NO_CHARACTERISTIC
}

var Abbreviation = map[Characteristic]string{
	Strength:        "STR",
	Dexterity:       "DEX",
	Endurance:       "END",
	Inteligence:     "INT",
	Education:       "EDU",
	SocialStanding:  "SOC",
	PsionicStrength: "PSI",
	Wealth:          "WLT",
	Sanity:          "SAN",
	Morale:          "MOR",
	Luck:            "LCK",
	Territory:       "TER",
	Charisma:        "CHA",
}

var NONE Characteristic = ""
