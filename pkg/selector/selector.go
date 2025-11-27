package selector

import (
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/consolio/prompt"
)

type Selector struct {
	dp *dice.Dicepool
}

func ManualSkillCli(skills ...skill.Skill) skill.Skill {
	items := []*prompt.Item{}
	for _, s := range skills {
		items = append(items, prompt.NewItem(s.String(), s))
	}
	chosen, err := prompt.SelectSingle(
		prompt.WithTitle("select skill"),
		prompt.FromItems(items),
	)
	if err != nil {
		panic(err)
	}
	selected := chosen.Payload().(skill.Skill)

	return selected
}
