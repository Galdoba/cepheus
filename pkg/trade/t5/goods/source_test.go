package goods

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	g := New("E888787-2")
	fmt.Println(g.ID())
	fmt.Println(g)
}
