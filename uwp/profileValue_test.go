package uwp

import (
	"fmt"
	"testing"
)

func Test_gravity(t *testing.T) {
	for i := 0; i <= 15; i++ {
		fmt.Println(i, ":", gravity(i))
		fmt.Println(Description(Size, fmt.Sprintf("%v", i))["ru"])
	}
}
