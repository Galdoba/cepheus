package money

import (
	"fmt"
	"testing"
)

func TestMegaCredit_String(t *testing.T) {
	val := MegaCredit(1.120909090909)
	fmt.Println(val)
	val2 := RU(12)
	fmt.Println(val2.MegaCredit())

}
