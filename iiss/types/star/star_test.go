package star

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestStarTypeDetermination(t *testing.T) {
	dp := dice.NewDicepool()
	types := make(map[string]int)
	for i := 0; i < 1000000; i++ {
		st, err := Generate(dp)

		if err != nil {
			panic(err)
		}
		types[st.Type]++
		if st.ProtoStar {
			fmt.Println("")
			fmt.Println(st.String())
		}
		switch st.Class {
		case "proto", "cluster", "anomaly":
			fmt.Println("")
			fmt.Println(st.String())
		}
		fmt.Printf("%v %v  \r", i, types)

		// // data, err := json.Marshal(&st)
		// // fmt.Println(i, string(data))
		// fmt.Println(st.String())
		// if st.Class == "BD" {
		// 	panic(2)
		// }
		// st2 := st.Sibling(dp)
		// fmt.Println(st2.String())
		// fmt.Println("===")

	}
	// st2, err2 := Generate(dp, KnownStellar("M5 IV"))
	// fmt.Println("===")
	// fmt.Println(err2)
	// data, err3 := json.Marshal(&st2)
	// fmt.Println(string(data), err3)

}
