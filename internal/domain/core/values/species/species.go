package species

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/core/values/characteristic"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
)

const (
	Human Specie = "Human"
	Aslan Specie = "Aslan"
	Vargr Specie = "Vargr"
)

type Specie string

func (r Specie) Characteristics() []characteristic.CharacteristicName {
	switch r {
	default:
		panic(fmt.Sprintf("race %v characteristic list is not implemented", r))
	case Human:
		return []characteristic.CharacteristicName{
			characteristic.Strength,
			characteristic.Dexterity,
			characteristic.Endurance,
			characteristic.Inteligence,
			characteristic.Education,
			characteristic.SocialStanding,
		}
	case Aslan:
		return append(Human.Characteristics(),
			characteristic.Territory,
		)
	case Vargr:
		return append(Human.Characteristics(),
			characteristic.Charisma,
		)
	}
}

func (r Specie) Skills() []skill.Skill {
	switch r {
	default:
		panic(fmt.Sprintf("race %v skill list is not implemented", r))
	case Human, Vargr:
		return []skill.Skill{}
	case Aslan:
		return []skill.Skill{
			skill.Independence,
			skill.Tolerance,
		}
	}
}
