package interaction

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/charmbracelet/huh"
)

type SelectionItem struct {
	Key      string
	Selected bool
	Data     any
}

func NewItem(key string, data any) SelectionItem {
	return SelectionItem{
		Key:      key,
		Selected: false,
		Data:     data,
	}
}

func (i SelectionItem) GetKey() string {
	return i.Key
}
func (i SelectionItem) IsSelected() bool {
	return i.Selected
}
func (i SelectionItem) PayData() any {
	return i.Data
}

type SelectorBuilder struct {
	autoselect bool
	autoMin    int
	autoMax    int
	itemPool   []SelectionItem
	single     *huh.Select[*SelectionItem]
	multi      *huh.MultiSelect[*SelectionItem]
}

func SelectSingle(title string, options ...SelectorOption) (*SelectionItem, error) {
	result := new(SelectionItem)
	prompt := SelectorBuilder{
		single: huh.NewSelect[*SelectionItem]().Title(title).Value(&result),
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

type SelectorOption func(*SelectorBuilder)

func WithItems(items ...SelectionItem) SelectorOption {
	return func(sb *SelectorBuilder) {
		opts := huh.NewOptions[*SelectionItem]()
		for _, item := range items {
			sb.itemPool = append(sb.itemPool, item)
			opts = append(opts, huh.NewOption[*SelectionItem](item.Key, &item))
		}
		sb.single.Options(opts...)
	}
}

func Auto(auto bool, minmax ...int) SelectorOption {
	return func(sb *SelectorBuilder) {
		sb.autoselect = auto
		for i, mm := range minmax {
			switch i {
			case 0:
				sb.autoMin = mm
				sb.autoMax = mm
			case 1:
				sb.autoMax = mm
			}

		}
		if sb.autoMin > sb.autoMax {
			sb.autoMin, sb.autoMax = sb.autoMax, sb.autoMin
		}
		if sb.autoMin+sb.autoMax == 0 {
			r := dice.FastRandom(fmt.Sprintf("1d%v", len(sb.itemPool)))
			sb.autoMin = r
			sb.autoMax = r
		}
	}
}
