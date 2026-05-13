package starsystem

import (
	"fmt"
	"testing"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder(BuilderConfiguration{
		Seed:       "",
		RegionType: Dense,
		SystemType: Realistic,
	})
	ss, err := b.Build()
	fmt.Println(err)
	fmt.Println(ss)
}
