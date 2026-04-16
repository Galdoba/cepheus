package uwp

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/engine/ehex"
)

const (
	Port = iota
	Size
	Atmosphere
	Hydrospere
	Population
	Government
	LawLevel
	separator
	TechLevel
)

type UWP struct {
	data           []ehex.Ehex
	representation string
}

func New() *UWP {
	var u UWP
	for _, p := range positions() {
		switch p {
		case separator:
			u.data = append(u.data, ehex.Ignore)
		default:
			u.data = append(u.data, ehex.Placeholder)
		}
	}
	return &u
}

func (u *UWP) Set(field int, value ehex.Ehex) error {
	switch field {
	case Port, Size, Atmosphere, Hydrospere, Population, Government, LawLevel, TechLevel:
		u.data[field] = value
	default:
		return fmt.Errorf("invalid field index: %v", field)
	}
	return nil
}

func positions() []int {
	return []int{
		Port,
		Size,
		Atmosphere,
		Hydrospere,
		Population,
		Government,
		LawLevel,
		separator,
		TechLevel,
	}
}


