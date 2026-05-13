package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

const (
	step_preparation = "preparing star system builder"
	step_01          = "determining the presence of an object in hex"
	step_02          = "determining the type of object"
	step_03          = "determining the type of brown dwarf"
	step_04          = "determining the type of a star in a system"
	step_05          = "determining the numerical classification"
	//TODO: add real functions
	done = "generation complete"
)

type BuilderConfiguration struct {
	Seed       string
	RegionType SubsectorType
	SystemType SystemType
}

type Builder struct {
	verbose   bool
	dice      *dice.Manager
	cfg       BuilderConfiguration
	nextStep  string
	buildFunc map[string]buildFunc
	precursor *StarSystem
}

type buildFunc func(*StarSystem) error

func NewBuilder(cfg BuilderConfiguration) *Builder {
	b := Builder{}
	b.cfg = cfg

	b.buildFunc = make(map[string]buildFunc)
	b.buildFunc[step_preparation] = b.Step_00
	b.buildFunc[step_01] = b.Step_01
	b.buildFunc[step_02] = b.Step_02
	b.buildFunc[step_03] = b.Step_03
	b.buildFunc[step_04] = b.Step_04
	b.nextStep = step_preparation
	b.precursor = &StarSystem{}
	return &b
}

func (b *Builder) Build() (*StarSystem, error) {
	starSystem := &StarSystem{}
	i := 0
	b.verbose = true
	for b.nextStep != done {
		b.print("%s...\n", b.nextStep)
		i++
		if fn, ok := b.buildFunc[b.nextStep]; ok {
			if err := fn(starSystem); err != nil {
				return nil, fmt.Errorf("build step %q failed: %w", b.nextStep, err)
			}

		} else {
			fmt.Println("no build func for", b.nextStep)
			break
		}
		if i > 50 {
			break
		}
	}

	return b.precursor, fmt.Errorf("not implemented")
}

func (b *Builder) Step_00(ss *StarSystem) error {
	dice, err := dice.New(b.cfg.Seed)
	if err != nil {
		return fmt.Errorf("failed to create dice manager: %w", err)
	}
	b.dice = dice
	b.precursor = &StarSystem{}
	b.nextStep = step_01
	return nil
}

func (b *Builder) print(f string, args ...any) {
	if b.verbose {
		fmt.Printf(f, args...)
	}
}
