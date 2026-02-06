package starsystem

import (
	"fmt"
	"testing"
)

func TestBuilder(t *testing.T) {
	lastPr := ""
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
		if ss.Profile() != lastPr {
			fmt.Println(ss.Profile())
			lastPr = ss.Profile()

		}
	}
}
