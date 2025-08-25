package uwp

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/pkg/ehex"
	"github.com/Galdoba/cepheus/pkg/rules"
)

func listDataTypes() []string {
	return []string{
		Port,
		Size,
		Atmo,
		Hydr,
		Pops,
		Govr,
		Laws,
		TL,
	}
}

// UWP - Represents Universal World Profile - short descriptive string with world stats.
type UWP struct {
	Data           map[string]ProfileValue `json:"profile"`
	RuleSystem     rules.RuleSystem        `json:"rules,omitempty"`
	GenerationTags []string                `json:"tags,omitempty"`
	Error          string                  `json:"error,omitempty"`
}

func New(opts ...UWP_Option) UWP {
	u := UWP{}
	u.RuleSystem = rules.T5
	u.Data = make(map[string]ProfileValue)
	for _, data := range listDataTypes() {
		u.Data[data] = ProfileValue{
			Category:  data,
			Code:      "?",
			Numerical: 0,
		}
	}
	for _, modify := range opts {
		modify(&u)
	}
	return u
}

type UWP_Option func(*UWP)

func FromString(s string) UWP_Option {
	return func(u *UWP) {
		if !StringValid(s) {
			fmt.Println("invalid uwp")
			return
		}
		data := strings.Split(s, "")
		for i, dataType := range listDataTypes() {
			if i == 7 {
				i++
			}
			u.Data[dataType] = ProfileValue{Category: dataType, Code: data[i], Numerical: ehex.FromString(data[i]).Value()}
		}
	}
}

func (u *UWP) ValueOf(key string) int {
	switch key {
	case Port, Size, Atmo, Hydr, Pops, Govr, Laws, TL:
		return u.Data[key].Numerical
	default:
		return -1
	}
}

func (u *UWP) CodeOf(key string) string {
	switch key {
	case Port, Size, Atmo, Hydr, Pops, Govr, Laws, TL:
		return u.Data[key].Code
	default:
		return ""
	}
}

// func New(dp dice.Dicepool, options ...UWP_Option) *UWP {
// 	uwp := UWP{}
// 	uwp.RuleSystem = rules.MgT2
// 	uwp.Data = make(map[string]ProfileValue)
// 	for _, modify := range options {
// 		modify(&uwp)
// 	}
// 	for _, datatype := range listDataTypes() {
// 		if uwp.Data[datatype] == nil {
// 			uwp.Data[datatype] = ehex.FromString("?")
// 		}
// 	}
// 	return &uwp
// }

// type UWP_Option func(*UWP)

// // KnownProfile - Adds known profile data
// func KnownProfile(profile string) UWP_Option {
// 	return func(u *UWP) {
// 		maps.Copy(u.Data, profileToData(profile))
// 	}
// }

// func profileToData(profile string) map[string]ehex.Ehex {
// 	data := make(map[string]ehex.Ehex)
// 	values := strings.Split(profile, "")
// 	if len(values) < 9 {
// 		return data
// 	}
// 	for i, val := range values {
// 		switch i {
// 		case 0:
// 			data[Port] = ehex.FromString(val)
// 		case 1:
// 			data[Size] = ehex.FromString(val)
// 		case 2:
// 			data[Atmosphere] = ehex.FromString(val)
// 		case 3:
// 			data[Hydrosphere] = ehex.FromString(val)
// 		case 4:
// 			data[Population] = ehex.FromString(val)
// 		case 5:
// 			data[Government] = ehex.FromString(val)
// 		case 6:
// 			data[Laws] = ehex.FromString(val)
// 		case 8:
// 			data[TL] = ehex.FromString(val)
// 		}
// 	}
// 	return data
// }
