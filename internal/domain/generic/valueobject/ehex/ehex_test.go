package ehex

import (
	"fmt"
	"testing"
)

func TestFromString(t *testing.T) {
	eh1 := FromInt(8)
	fmt.Println(eh1)
	eh2 := FromString("4")
	fmt.Println(eh2)
	eh3 := Add(eh1, eh2)
	fmt.Println(eh3)
}

func TestPrintRunes(t *testing.T) {
	s := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz?^*&"
	PrintRunes(s)
}
