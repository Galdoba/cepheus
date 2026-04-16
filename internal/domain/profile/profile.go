package profile

import "github.com/Galdoba/cepheus/internal/domain/engine/ehex"

type Profile struct {
	blocks []*profileBlock
}

type profileBlock struct {
	eh        ehex.Ehex
	float     float64
	prefix    string
	text      string
	suffix    string
	separator string
}

var uwpMap = map[int]int{}

func (p *Profile) get(index int) *profileBlock {

}

// String return profile's string representation
func (p *Profile) String() string {
	return ""
}



