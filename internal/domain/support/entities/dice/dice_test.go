package dice_test

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/core/entities/dice"
)

func TestDicepool_Roll(t *testing.T) {
	dp := dice.New("galdoba")
	dp.Roll("2D6")

	fmt.Println(dice.Roll("3D"))
	fmt.Println(dice.D66(-9, 3))
}
