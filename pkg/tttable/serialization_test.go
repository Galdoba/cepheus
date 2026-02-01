package tttable

import (
	"fmt"
	"os"
	"testing"
)

func TestLiveAction(t *testing.T) {
	table, err := NewTable("test table",
		WithDiceExpression("1d6"),
		WithMods(map[string]int{"no skill": -3}),
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

	if err != nil {
		fmt.Println(err)
		return
	}
	table, err = Load(`/home/galdoba/workspace/test_tab.json`)
	table.SetWriter(os.Stderr)

	if err != nil {
		fmt.Println(err)
		return
	}
}
