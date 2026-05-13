package star

import (
	"fmt"
	"testing"
)

func TestDefault(t *testing.T) {
	for k, v := range defaultRegistry.Data {
		// fmt.Println(k, v)
		if len(v) != 8 {
			fmt.Printf("%s: have %v elements\n", k, len(v))
		}
		if v[6] < v[5] && v[6] != 0 && v[5] != 0 {
			fmt.Printf("%s: have minmax habitable\n", k)
		}
		if v[6] > v[7] && v[6] != 0 {
			fmt.Printf("%s: have minmax outer\n", k)
		}
		if v[4] > v[5] && v[5] != 0 {
			fmt.Printf("%s: have minmax inner\n", k)
		}
	}

}
