package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

const (
	step_preparation = "preparation"
	step_01          = "determine the presence of an object in hex"
	//TODO: add real functions
	done = "generation complete"
)

type BuilderConfiguration struct {
	Seed       string
	RegionType string
}

type Builder struct {
	dice      *dice.Manager
	cfg       BuilderConfiguration
	nextStep  string
	buildFunc map[string]buildFunc
}

type buildFunc func(*StarSystem) error

func NewBuilder(cfg BuilderConfiguration) *Builder {
	b := Builder{}
	b.cfg = cfg
	b.buildFunc = make(map[string]buildFunc)
	b.buildFunc[step_preparation] = b.Step_00
	b.buildFunc[step_01] = b.Step_01
	return &b
}

func (b *Builder) Build() (*StarSystem, error) {
	starSystem := &StarSystem{}
	for b.nextStep != done {
		fn := b.buildFunc[b.nextStep]
		if err := fn(starSystem); err != nil {
			return nil, fmt.Errorf("build step %q failed: %w", b.nextStep, err)
		}

	}

	return &StarSystem{}, fmt.Errorf("not implemented")
}

func (b *Builder) Step_00(ss *StarSystem) error {
	dice, err := dice.New(b.cfg.Seed)
	if err != nil {
		return fmt.Errorf("failed to create dice manager: %w", err)
	}
	b.dice = dice
	b.nextStep = step_01
	return nil
}

func (b *Builder) Step_01(ss *StarSystem) error {
	fmt.Println(b.dice.MustRoll("2d6"))
	b.nextStep = done
	return nil
}

func defineNextStep(ss *StarSystem) string {
	return ""
}
