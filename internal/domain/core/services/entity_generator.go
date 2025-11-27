package services

import (
	"github.com/Galdoba/cepheus/internal/domain/core/aggregates/typeset"
	"github.com/Galdoba/cepheus/internal/domain/core/entities/entity"
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

func NewSophont(opts ...entity.SophontOption) *entity.Sophont {
}
