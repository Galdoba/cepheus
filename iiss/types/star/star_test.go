package star

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestStarTypeDetermination(t *testing.T) {
	dp := dice.NewDicepool()
	for i := 0; i < 200000; i++ {
		st, err := StarTypeDetermination(dp)
		if err != nil {
			panic(i)
		}
		pr := strings.Split(st, "|")
		fmt.Printf("%v\t%v%v %v\n", i, pr[0], pr[1], pr[2])
		if pr[2] == "VI" {
			panic(1)
		}
	}
}
