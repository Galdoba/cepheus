package star

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestStarTypeDetermination(t *testing.T) {
	dp := dice.NewDicepool()
	for i := 0; i < 50; i++ {
		st, err := Generate(dp)
		if err != nil {
			panic(err)
		}
		data, err := json.Marshal(&st)
		fmt.Println(i, string(data))
		if st.Class == "BD" {
			panic(2)
		}
	}
	st2, err2 := Generate(dp, KnownClass("BD"))
	fmt.Println("===")
	fmt.Println(err2)
	data, err3 := json.Marshal(&st2)
	fmt.Println(string(data), err3)

}
