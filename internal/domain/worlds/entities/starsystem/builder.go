package starsystem

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
	"github.com/Galdoba/cepheus/pkg/dice"
)

type Builder struct {
	rng      *dice.Roller
	options  map[string]bool
	imported t5ss.WorldData
	step1    *primaryStarDeterminator
	step2    *secondaryStarsDeterminator
}

type primaryStarDeterminator struct {
	tablesStarType rtg.RandomTableGenerator
	activeMods     map[string]bool
	completed      bool
}

type secondaryStarsDeterminator struct {
	tables     rtg.RandomTableGenerator
	starSchema []stellar.StarDesignation
	completed  bool
}

type BuildOption func(*Builder)

func NewBuilder(seed string, options ...BuildOption) (*Builder, error) {
	b := Builder{}
	b.rng = dice.New(seed)
	b.options = make(map[string]bool)

	psd := primaryStarDeterminator{}
	rtg1, err := rtg.NewStarTypeDeterminationGenerator(b.rng)
	if err != nil {
		return nil, fmt.Errorf("failed to create RNG: Star Type Determination: %v", err)
	}
	psd.tablesStarType = rtg1

	ssd := secondaryStarsDeterminator{}
	ssd.tables = rtg1

	b.step1 = &psd
	return &b, nil
}

func newStarSystem() *StarSystem {
	ss := StarSystem{}
	ss.Stars = make(map[orbit.Orbit]*Star)
	return &ss
}

func (b *Builder) Build(directives ...string) (*StarSystem, error) {

	ss := &StarSystem{}
	if err := b.runStep1(ss); err != nil {
		return nil, fmt.Errorf("step 1 failed: %v", err)
	}

	return ss, nil
}
