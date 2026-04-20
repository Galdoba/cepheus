package uwp

import (
	"fmt"
	"strings"

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
	data []ehex.Ehex
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

func FromString(s string) (*UWP, error) {
	if len(s) != 9 {
		return nil, fmt.Errorf("string must be exactly 9 charaters: %s", s)
	}
	u := New()
	for i, code := range strings.Split(s, "") {
		u.Set(i, ehex.FromCode(code))
	}
	return u, nil
}

func (u *UWP) Port() int {
	return u.data[Port].Value()
}

func (u *UWP) Size() int {
	return u.data[Size].Value()
}

func (u *UWP) Atmosphere() int {
	return u.data[Atmosphere].Value()
}

func (u *UWP) Hydrospere() int {
	return u.data[Hydrospere].Value()
}

func (u *UWP) Population() int {
	return u.data[Population].Value()
}

func (u *UWP) Government() int {
	return u.data[Government].Value()
}

func (u *UWP) LawLevel() int {
	return u.data[LawLevel].Value()
}

func (u *UWP) TL() int {
	return u.data[TechLevel].Value()
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

func (u *UWP) String() string {
	s := ""
	for i := range u.data {
		s += u.data[i].Code()
	}
	return s
}

func (u *UWP) Raw() []ehex.Ehex {
	return u.data
}
