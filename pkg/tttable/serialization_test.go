package tttable

import (
	"fmt"
	"os"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestLiveAction(t *testing.T) {
	table, err := NewTable("test table",
		WithDiceExpression("1d6"),
		WithMods(map[string]int{"no skill": -3}),
		WithRoller(dice.New("42")),
		WithRows(
			NewRow("1-", "No, and..."),
			NewRow("2", "No."),
			NewRow("3", "No, but..."),
			NewRow("4", "Yes, but..."),
			NewRow("5", "Yes."),
			NewRow("6+", "Yes, and..."),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := table.SaveAs(`/home/galdoba/workspace/test_tab.json`); err != nil {
		fmt.Println(err)
	}
	result, err := table.Roll(" ")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("result:", result)
	table, err = Load(`/home/galdoba/workspace/test_tab.json`)
	table.SetRoller(dice.New("42"))
	table.SetWriter(os.Stderr)
	result, err = table.Roll("no skill")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("result:", result)
}
