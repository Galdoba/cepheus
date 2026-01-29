package starsystem

import (
	"fmt"
	"testing"
)

func TestBuilder(t *testing.T) {
	for i := 160000; i < 1000000; i++ {
		b, err := NewBuilder(fmt.Sprintf("%v", i))
		if err != nil {
			fmt.Println(err)
			return
		}
		ss, err := b.Build()
		if err != nil {
			fmt.Println(err)
			return
		}
		// if ss.PrimaryStar.Class == "IV" {
		fmt.Println(i, ss.PrimaryStar)
		fmt.Println(ss)
		// }
		if ss.PrimaryStar.Class == "IV" && ss.PrimaryStar.Type == "O" {
			panic(1)

		}
	}
}
