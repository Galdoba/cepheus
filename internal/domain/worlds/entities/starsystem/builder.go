package starsystem

import (
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/tttable"
)

type Builder struct {
	rng     *dice.Roller
	options map[string]bool
	step1   *primaryStarDeterminator
}

type primaryStarDeterminator struct {
	tablesStarType    *tttable.TableCollection
	starTypeMods      map[string]bool
	tablesStarSubType *tttable.TableCollection
	completed         bool
}

type BuildOption func(*Builder)

func NewBuilder(options ...BuildOption) (*Builder, error) {
	b := Builder{}
	b.options = make(map[string]bool)
	return &b, nil
}

func (b *Builder) Build(directives ...string) (*StarSystem, error) {
	return &StarSystem{}, nil
}
