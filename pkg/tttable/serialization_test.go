package tttable

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLiveAction(t *testing.T) {
	table, err := NewTable("test table",
		WithDiceExpression("1d6"),
		WithIndexMods(Flat, map[string]int{"no skill": -3}),
		WithIndexEntries(
			NewTableEntry("1-", "No, and..."),
			NewTableEntry("2", "No."),
			NewTableEntry("3", "No, but..."),
			NewTableEntry("4", "Yes, but..."),
			NewTableEntry("5", "Yes."),
			NewTableEntry("6+", "Yes, and..."),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	hm, _ := os.UserHomeDir()
	path := filepath.Join(hm, "workspace", "test_tab.json")
	os.MkdirAll(filepath.Dir(path), 0755)
	if err := SaveAs(table, path); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = Load(path)

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(loaded.GetName(), "loaded")
}
