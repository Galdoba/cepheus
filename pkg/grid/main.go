package main

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/grid/coordinates"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GenerateBeltCenters(n int) [][3]int {
	if n == 0 {
		return [][3]int{{0, 0, 0}}
	}

	total := 2 * n
	var centers [][3]int

	for i := -total; i <= total; i++ {
		rem := total - abs(i)
		for j := -rem; j <= rem; j++ {
			k := -i - j
			if abs(i)+abs(j)+abs(k) == total {
				centers = append(centers, [3]int{24 * i, 24 * j, 24 * k})
			}
		}
	}

	return centers
}

func main() {

	n := 0
	for i := 0; i < 2; i++ {
		for j, crd := range GenerateBeltCenters(i) {
			fmt.Println(n, i, j, crd)
			coords := coordinates.NewSpaceCoordinates(crd[0], crd[1], crd[2])
			fmt.Println(coords.GlobalValues())
			fmt.Println(coords.LocalValues())
			fmt.Println(coords.StringSectorNameHex())
			n++
		}
	}
}
