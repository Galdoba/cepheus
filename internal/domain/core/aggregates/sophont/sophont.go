package sophont

import (
	"fmt"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/core/entities/value"
	"github.com/Galdoba/cepheus/internal/domain/core/values/characteristic"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
	"github.com/Galdoba/cepheus/internal/domain/core/values/species"
)

type Sophont struct {
	uuid      string
	name      string
	specie    species.Specie
	char      *characteristics
	skill     *skills
	modifiers map[string]bool
}

func NewSophont(opts ...SophontOption) *Sophont {
	s := Sophont{
		uuid:   fmt.Sprintf("Sophont %v", time.Now().Unix()),
		name:   "Some Name",
		specie: species.Human,
	}
	for _, modify := range opts {
		modify(&s)
	}
	s.char = newCharacteristics(s.specie)
	s.skill = newSkills(s.specie)
	s.modifiers = make(map[string]bool)
	return &s
}

type SophontOption func(*Sophont)

type characteristics struct {
	byName map[characteristic.CharacteristicName]*value.CharacteristicValue
}

func newCharacteristics(specie species.Specie) *characteristics {
	charSet := make(map[characteristic.CharacteristicName]*value.CharacteristicValue)
	for _, chr := range specie.Characteristics() {
		charSet[chr] = value.NewCharacteristicValue()
	}
	return &characteristics{charSet}
}

type skills struct {
	byName map[skill.Skill]*value.SkillValue
}

func newSkills(specie species.Specie) *skills {
	skillSet := make(map[skill.Skill]*value.SkillValue)
	for _, skl := range specie.Skills() {
		skillSet[skl] = value.NewSkillValue()
	}
	return &skills{skillSet}
}

func (s *Sophont) CreationSheet() string {

}
