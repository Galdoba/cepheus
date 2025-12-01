package services

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewCharacterGenerator(t *testing.T) {
	cg := NewCharacterGenerator(CharacterGenerationOptions{
		ManualDecidions: false,
		Seed:            42,
	})
	p := tea.NewProgram(cg)
	lastState, err := p.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(lastState)
}
