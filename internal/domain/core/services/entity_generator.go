package services

import (
	"time"

	"github.com/Galdoba/cepheus/internal/domain/core/aggregates/sophont"
	"github.com/Galdoba/cepheus/pkg/dice"
	tea "github.com/charmbracelet/bubbletea"
)

type CharacterGenerator struct {
	dp              *dice.Dicepool
	manualDecidions bool
	character       *sophont.Sophont
}

type CharacterGenerationOptions struct {
	ManualDecidions bool
	Seed            int64
}

func NewCharacterGenerator(cfg CharacterGenerationOptions) CharacterGenerator {
	cg := CharacterGenerator{}
	cg.manualDecidions = cfg.ManualDecidions
	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().Unix()
	}
	cg.dp = dice.NewDicepool(dice.WithSeed(seed))
	return cg
}

func (cg CharacterGenerator) Init() tea.Cmd {
	return nil
}

func (cg CharacterGenerator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return cg, tea.Quit
		}
	}
	return cg, nil
}

func (cg CharacterGenerator) View() string {
	return "cg view"
}

type CharacterExport struct {
	Name string
}
