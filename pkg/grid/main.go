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
	for i := 0; i < 23; i++ {
		for _, crd := range GenerateBeltCenters(i) {
			// fmt.Println(n, i, j, crd)
			coords := coordinates.NewSpaceCoordinates(crd[0], crd[1], crd[2])
			x, y := coords.GlobalValues()
			q, r, s := coords.CubeValues()
			fmt.Println(n, x, y, "|", q, r, s)
			n++
		}
	}
}
