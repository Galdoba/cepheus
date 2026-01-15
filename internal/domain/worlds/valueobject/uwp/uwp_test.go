package uwp

import (
	"fmt"
	"testing"
)

func TestUWP_Description(t *testing.T) {
	u, err := New("A123AFF-7")
	fmt.Println(err)
	fmt.Println(u)
	fmt.Println(u.Description())
}
