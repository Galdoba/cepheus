package uwp

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/pkg/ehex"
)

type UWP string

//A123456-8

func (u UWP) Port() string {
	sl := strings.Split(string(u), "")
	return sl[0]
}

func (u UWP) Size() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[1])
}

func (u UWP) Atmo() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[2])
}

func (u UWP) Hydr() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[2])
}

func (u UWP) Pops() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[4])
}

func (u UWP) Govr() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[5])
}

func (u UWP) Laws() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[6])
}

func (u UWP) TL() ehex.Ehex {
	sl := strings.Split(string(u), "")
	return ehex.FromString(sl[8])
}

func (u UWP) StatBlock() string {
	sb := ""
	sb += fmt.Sprintf("Spaceport: %v (%v)", u.Port(), describePort(u.Port()))
	return sb
}

func describePort(port string) string {
	s := ""
	switch port {
	case "A":
		s = "Exelent"
	case "B":
		s = "Good"
	case "C":
		s = "Routine"
	case "D":
		s = "Poor"
	case "E":
		s = "Frontier"
	case "F":
		s = "Good"
	case "G":
		s = "Poor"
	case "H":
		s = "Basic"
	case "Y":
		s = "None"
	case "X":
		s = "None"
	default:
		s = "invalid port code"
	}
	return s
}
