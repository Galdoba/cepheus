package services

import (
	"github.com/Galdoba/cepheus/internal/domain/core/aggregates/typeset"
	"github.com/Galdoba/cepheus/internal/domain/core/values/characteristic"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
)

type EntityGenerator struct {
}

func (eg *EntityGenerator) GetCharacteristics() []*typeset.Collection[characteristic.Characteristic] {
	return nil
}
func (eg *EntityGenerator) GetSkills() []*typeset.Collection[skill.Skill] {
	return nil
}
func (eg *EntityGenerator) GetTraits() []characteristic.Characteristic {
	return nil
}

type Asset interface {
	DM() int
}

func (e *EntityGenerator) Roll(code string, mods ...Asset) {

}
