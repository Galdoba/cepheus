package star

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestStarTypeDetermination(t *testing.T) {
	dp := dice.NewDicepool()
	for i := 0; i < 20; i++ {
		st, err := Generate(dp)
		if err != nil {
			panic(err)
		}
		// data, err := json.Marshal(&st)
		// fmt.Println(i, string(data))
		fmt.Println(st.String())
		if st.Class == "BD" {
			panic(2)
		}
		st2 := st.Sibling(dp)
		fmt.Println(st2.String())
		fmt.Println("===")

	}
	st2, err2 := Generate(dp, KnownStellar("M5 IV"))
	fmt.Println("===")
	fmt.Println(err2)
	data, err3 := json.Marshal(&st2)
	fmt.Println(string(data), err3)

}
