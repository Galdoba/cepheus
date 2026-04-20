package star

import (
	"fmt"
	"testing"
)

func Test_defineStarKeys(t *testing.T) {
	// "D" заменён на "D7" (валидный), добавлен "BD"
	for _, sk := range Parse("G2 V O3 III BD; D7; L4  D2") {
		fmt.Println(sk)
		s, n, l := parseKey(sk)
		fmt.Printf("stellar=%q; num=%q luma=%q\n", s, n, l)
	}
}
