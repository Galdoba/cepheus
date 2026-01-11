package characteristic

import "strings"

const (
	TypePhysical      CharacteristicType = "Physical"
	TypeMental        CharacteristicType = "Mental"
	TypeObscured      CharacteristicType = "Obscured"
	TypeCore          CharacteristicType = "Core"
	TypeSecondary     CharacteristicType = "Obscured"
	NO_CHARACTERISTIC CharacteristicName = "[NO CHARACTERISTIC]"
	Strength          CharacteristicName = "Strength"
	Dexterity         CharacteristicName = "Dexterity"
	Endurance         CharacteristicName = "Endurance"
	Inteligence       CharacteristicName = "Inteligence"
	Education         CharacteristicName = "Education"
	SocialStanding    CharacteristicName = "Social Standing"
	PsionicStrength   CharacteristicName = "Psionic Strength"
	Wealth            CharacteristicName = "Wealth"
	Sanity            CharacteristicName = "Sanity"
	Morale            CharacteristicName = "Morale"
	Luck              CharacteristicName = "Luck"
	Territory         CharacteristicName = "Territory"
	Charisma          CharacteristicName = "Charisma"
)

type CharacteristicType string

type CharacteristicName string

func Call(key string, list ...CharacteristicName) CharacteristicName {
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

var Abbreviation = map[CharacteristicName]string{
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

type CharacteristicValue int

func (c CharacteristicValue) BaseValue() int {
	return int(c)
}
