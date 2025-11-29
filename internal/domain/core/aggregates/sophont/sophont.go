package sophont

import (
	"github.com/Galdoba/cepheus/internal/domain/core/entities/value"
	"github.com/Galdoba/cepheus/internal/domain/core/values/characteristic"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
	"github.com/Galdoba/cepheus/internal/domain/core/values/species"
)

type Sophont struct {
	uuid      string
	name      string
	specie    species.Specie
	chars     *characteristics
	modifiers map[string]bool
}

type characteristics struct {
	byName map[characteristic.CharacteristicValue]*value.CharacteristicValue
}

type skills struct {
	byName map[skill.Skill]*value.SkillValue
}
