package dice_test

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

func TestNew(t *testing.T) {
	m, err := dice.New(nil)
	fmt.Println(err)
	fmt.Println(m.Roll(" 2d6+100"))

	fmt.Println(dice.D66(-9, -9))
}
