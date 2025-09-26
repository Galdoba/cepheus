package interaction

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/charmbracelet/huh"
)

func MultiSelect(title string, options ...SelectorOption) (*SelectionItem, error) {
	result := []*SelectionItem{}
	prompt := SelectorBuilder{
		multi: huh.NewMultiSelect[*SelectionItem]().Title(title).Value(&result),
	}
	for _, modify := range options {
		modify(&prompt)
	}

	if len(prompt.itemPool) == 0 {
		return nil, fmt.Errorf("nothing to select from")
	}
	switch prompt.autoselect {
	case true:
		result = &prompt.itemPool[dice.RandomIndex(prompt.itemPool)]
	case false:
		if err := prompt.single.Run(); err != nil {
			return nil, err
		}
	}

	return result, nil
}
