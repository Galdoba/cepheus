package star

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/Galdoba/cepheus/pkg/dice"
)

type Star struct {
	Type           string   `json:"type,omitempty"`
	SubType        *int     `json:"subtype,omitempty"`
	Class          string   `json:"class,omitempty"`
	Mass           float64  `json:"mass,omitempty"`
	Temperature    int      `json:"temperature,omitempty"`
	Diameter       float64  `json:"diameter,omitempty"`
	Luminocity     float64  `json:"luminocity,omitempty"`
	Designation    *string  `json:"designation,omitempty"`
	OrbitN         *float64 `json:"orbit#,omitempty"`
	SystemAU       *float64 `json:"au primary,omitempty"`
	Eccentricity   *float64 `json:"eccentricity,omitempty"`
	Age            float64  `json:"age,omitempty"`
	realetdPrimary *Star
}

func Generate(dp *dice.Dicepool, knownData ...KnownStarData) (Star, error) {
	st := Star{}
	for _, add := range knownData {
		add(&st)
	}
	if st.Type+st.Class == "" && st.realetdPrimary == nil {
		tpe, cls, err := StarTypeAndClassDetermination(dp)
		if err != nil {
			return st, fmt.Errorf("failed to determine type and class of primary star: %v", err)
		}
		st.Type = tpe
		st.Class = cls

		st.SubType, err = StarSubTypeDetermination(dp, st)
		if err != nil {
			return st, fmt.Errorf("failed to determine subtype of the star: %v", err)
		}
	}

	st = fixClass(st, dp)
	st = SetParameters(st, dp)
	return st, nil
}

func SetParameters(st Star, dp *dice.Dicepool) Star {
	mass := 0.0
	switch st.Class {
	case "BD":
		st.Type, st.SubType, mass = bdTypeDetails(dp)
	case "D":
		mass = whiteDwarfMass(dp)
	case "NS", "BH", "PSR":
		mass = heavyDwarfMass(dp)
		if mass > 2.16 {
			st.Class = "BH"
		}
		st.Mass = roundFloat(mass)
	default:
		mass = adjust(dp, massByIndex(st.index()))
	}
	st.Mass = roundFloat(mass)
	diam := adjust(dp, diamByIndex(st.index()))
	st.Diameter = roundFloat(diam)
	temp := adjust(dp, tempByIndex(st.index()))
	st.Temperature = int(temp)
	st.Luminocity = roundFloat(calculateLuminosity(st.Diameter, temp))
	st.Age = roundFloat(generateAge(dp, st))
	if st.Class == "D" {
		st.Diameter = roundFloat(whiteDwarfDiameter(st.Mass))
		st.Temperature = int(whiteDwarfTemp(st.Age))
		st.Luminocity = roundFloat(calculateLuminosity(st.Diameter, float64(st.Temperature)))
	}
	return st
}

func heavyDwarfMass(dp *dice.Dicepool) float64 {
	m := 1.0
	r := dp.Sum1D()
	m += (float64(r) / 10.0) + (float64(dp.Sum1D()) / 100.0)
	for r == 6 {
		r = dp.Sum1D()
		m += float64(r-1) / 10.0
	}
	return m
}

func adjust(dp *dice.Dicepool, fl float64) float64 {
	adj := float64(dp.Flux())*0.01 + 1.0 + float64(dp.Flux())*0.001
	return fl * adj
}

func calculateLuminosity(diameter, temperature float64) float64 {
	diameterComponent := math.Pow(diameter, 2)
	temperatureComponent := math.Pow(temperature/5772, 4)
	luminosity := diameterComponent * temperatureComponent
	return luminosity
}

func generateAge(dp *dice.Dicepool, st Star) float64 {
	age := 0.0
	msls := mainSequanceLifespan(st.Mass)
	sgls := msls / (4.0 / st.Mass)
	glls := msls / (10.0 / math.Pow(st.Mass, 3))

	switch st.Class {
	case "Ia", "Ib", "II", "V", "VI":
		switch st.Mass > 0.9 {
		case false:
			age = smallStarAge(dp)
		case true:
			age = msls
		}
	case "IV":
		age = msls + sgls
	case "III":
		age = msls + sgls + glls
	case "BD":
		age = smallStarAge(dp)
	default:
		mass := originalMass(dp, st.Mass)
		msls = mainSequanceLifespan(mass)
		sgls = msls / (4.0 / st.Mass)
		glls = msls / (10.0 / math.Pow(mass, 3))
		finalAge := msls + sgls + glls
		switch st.Class {
		case "D", "NS", "BH", "VII":
			age = finalAge + smallStarAge(dp)
		case "PSR":
			age = 0.1/(float64(dp.Sum("2d10"))) + finalAge
		case "Protostar":
			age = 0.01 / float64(dp.Sum("2d10"))
		}

	}
	age = adjust(dp, age*variance(dp))
	if st.Mass < 4.7 && age < 0.01 {
		age = 0.01
	}
	for age > 13.0 {
		age = age / 10
	}
	return roundFloat(age)
}

func bdTypeDetails(dp *dice.Dicepool) (string, *int, float64) {
	mass := float64(dp.Sum("1d6"))/100.0 + float64(dp.Sum("4d6")-1)/1000.0
	vals := []float64{0.080, 0.076, 0.072, 0.068, 0.064, 0.060, 0.058, 0.056, 0.054, 0.052, 0.050, 0.048, 0.046, 0.044, 0.042, 0.040, 0.037, 0.034, 0.031, 0.028, 0.025, 0.022, 0.019, 0.016, 0.014, 0.013, 0.012, 0.011, 0.010, 0.009}
	index := -1
	for i, v := range vals {
		if v <= mass {
			index = i
			break
		}
	}
	tp := ""
	switch index / 10 {
	case 0:
		tp = "L"
	case 1:
		tp = "T"
	case 2:
		tp = "Y"
	}
	stp := index % 10
	return tp, &stp, mass
}

func whiteDwarfMass(dp *dice.Dicepool) float64 {
	return float64(dp.Sum("2d6")-1)/10.0 - float64(dp.Sum("1d10"))/100.0
}

func whiteDwarfDiameter(mass float64) float64 {
	return (1 / mass) * 0.01
}

func whiteDwarfTemp(age float64) float64 {
	return interpolateWhiteDwarfTemp(age)
}

type KnownStarData func(*Star)

func KnownType(sType string) KnownStarData {
	return func(s *Star) {
		s.Type = sType
	}
}

func KnownClass(class string) KnownStarData {
	return func(s *Star) {
		s.Class = class
	}
}

func (st Star) String() string {
	data, err := json.Marshal(st)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (s Star) ToStellar() string {
	switch s.Class {
	case "D", "BD", "BH", "NS", "PSR", "":
		return s.Class
	}
	return fmt.Sprintf("%v%v %v", s.Type, *s.SubType, s.Class)
}

func fixClass(st Star, dp *dice.Dicepool) Star {
	switch st.Type {
	case "O", "B", "A", "F", "G", "K", "M", "L", "T", "Y", "D", "BD":
		if st.SubType == nil {
			r := dp.Sum("1d10") - 1
			st.SubType = &r
		}
	}
	if st.Class == "IV" {
		if st.Type == "O" || st.Type == "M" {
			st.Class = "V"
		}
		if st.Type == "K" && *st.SubType > 5 {
			st.Class = "V"
		}
	}
	if st.Class == "VI" {
		if st.Type == "A" || st.Type == "F" {
			st.Class = "V"
		}
	}
	return st
}
