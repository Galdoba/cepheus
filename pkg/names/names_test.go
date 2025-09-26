package names_test

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/names"
)

func TestPrintMap(t *testing.T) {
	// names.PrintMap(2)
	fmt.Println(names.RandomMaleName())
	fmt.Println(names.RandomFemaleName())
	fmt.Println(names.RandomLastName())
}
