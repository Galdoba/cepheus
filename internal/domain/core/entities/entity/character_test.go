package entity

import (
	"fmt"
	"testing"
)

func TestNewSophont(t *testing.T) {
	s, err := NewSophont()
	fmt.Println(err)
	fmt.Println(s)
}
