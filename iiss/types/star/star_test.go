package star

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestStarTypeDetermination(t *testing.T) {
	dp := dice.NewDicepool()
	for i := 0; i < 200000; i++ {
		st, err := Generate(dp)
		fmt.Printf("%v\t%v\n", i, st.String())
		if err != nil {
			panic(err)
		}
	}
}
