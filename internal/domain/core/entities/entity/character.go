package entity

import (
	"fmt"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/core/aggregates/typeset"
	"github.com/Galdoba/cepheus/internal/domain/core/entities/value"
	"github.com/Galdoba/cepheus/internal/domain/core/values/characteristic"
	"github.com/Galdoba/cepheus/internal/domain/core/values/race"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
	"github.com/Galdoba/cepheus/pkg/dice"
)

type Sophont struct {
	name               string
	race               race.Race
	seed               int64
	generationComplete bool
	characteristics    *typeset.Collection[characteristic.Characteristic]
	skills             *typeset.Collection[skill.Skill]
}

type sophontGenerator struct {
	dp *dice.Dicepool
}

func NewSophont(opts ...SophontOption) (*Sophont, error) {
	soph := Sophont{}
	soph.seed = time.Now().UnixNano()
	soph.name = fmt.Sprintf("Sophont %v", soph.seed)
	soph.race = race.Human
	for _, modify := range opts {
		modify(&soph)
	}
	// generate characteristics
	soph.characteristics = typeset.New[characteristic.Characteristic]()
	for _, chr := range soph.race.Characteristics() {
		val := value.New(value.ValueFor(value.ValueTypeCharacteristic))
		soph.characteristics.Set(chr, *val)
	}
	// generate skills
	soph.skills = typeset.New[skill.Skill]()
	for _, skl := range soph.race.Skills() {
		val := value.New(value.ValueFor(value.ValueTypeSkill))
		soph.skills.Set(skl, *val)
	}
	// generate career path
	return &soph, nil

}

type SophontOption func(*Sophont)

func WithSeed(val int64) SophontOption {
	return func(s *Sophont) {
		s.seed = val
	}
}

func WithRace(val race.Race) SophontOption {
	return func(s *Sophont) {
		s.race = val
	}
}

func WithName(val string) SophontOption {
	return func(s *Sophont) {
		s.name = val
	}
}
