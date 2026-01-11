package worldsize

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/support/valueobject/uwp"
)

func TestWorldSize(t *testing.T) {
	uwp := uwp.UWP("A623456-8")
	// dp := dice.NewDicepool(dice.WithSeedString("test"))
	ws1 := New(uwp.Size())
	fmt.Println(ws1)
	fmt.Println(ws1.Profile())
}
