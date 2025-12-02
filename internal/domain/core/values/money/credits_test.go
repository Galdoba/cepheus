package money_test

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/core/values/money"
)

func TestMegaCredit_String(t *testing.T) {
	val := money.MegaCredit(1.120909090909)
	fmt.Println(val)
	val2 := money.RU(12)
	fmt.Println(val2.MegaCredit())

}
