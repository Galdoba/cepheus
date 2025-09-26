package interaction

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/huh"
)

func TestSelectSingle(t *testing.T) {
	val := ""
	sel := huh.NewSelect[string]().Value(&val).Options(
		huh.NewOption("one", "one"),
		huh.NewOption("two", "two"),
		huh.NewOption("three", "three"),
	)
	sel.Run()
	fmt.Println(val)
	selected, err := SelectSingle("mego", WithItems(NewItem("one", 1), NewItem("two", "bar")), Auto(true))
	fmt.Println(err)
	fmt.Println(selected)

}
