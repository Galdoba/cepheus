package orbit

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/float"
)

func TestOrbitN_To_AU(t *testing.T) {
	for i := -0.02; i < 15.1; i = i + 0.01 {
		fmt.Printf("%v\t=>\t%v\n", float.RoundN(i, 2), OrbitN_To_AU(i))
		fmt.Printf("%v\t=>\t%v\n", float.RoundN(i, 2), OrbitN_To_Mkm(i))
	}
}
