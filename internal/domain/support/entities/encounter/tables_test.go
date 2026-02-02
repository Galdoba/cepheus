package encounter

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/tttable"
)

func TestPerson(t *testing.T) {
	dp := dice.New("")
	i1, r1, err1 := EncounterDistance().Roll(dp)
	i2, r2, err2 := Person().Roll(dp)
	i3, r3, err3 := CharacterQuircks().Roll(dp)

	cl, err := tttable.NewTableCollection(
		tttable.WithTables(
			EncounterDistance(), Person(), CharacterQuircks(),
		),
	)
	if err != nil {
		fmt.Println(err)
	}
	cl.SetRoller(dp)
	fmt.Println(cl.Roll("Encounter Distance"))

	fmt.Printf("at %v distance you see %v who seems like %v\n", r1, r2, r3)
	fmt.Println(i1, i2, i3, err1, err2, err3)
}
