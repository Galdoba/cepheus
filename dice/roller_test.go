package dice

import (
	"fmt"
	"testing"
)

func Test_parseDiceCode(t *testing.T) {
	parseDiceCode("dd    + 2")
	dp := NewDicepool(WithSeed(111))
	mdMap := make(map[string]int)
	mdMap["pistol"] = 2
	mdMap["auto"] = 6
	mdMap["laser"] = 1
	r := dp.Sum("3d1", DM_conditional(mdMap, "slug", "laser"), ForEveryDice(1))
	fmt.Println(r, mdMap)
}
