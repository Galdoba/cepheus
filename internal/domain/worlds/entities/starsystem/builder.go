package starsystem

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Galdoba/cepheus/internal/domain/support/services/float"
	"github.com/Galdoba/cepheus/internal/domain/worlds/services/interpolate"
	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
	"github.com/Galdoba/cepheus/pkg/dice"
)

type Builder struct {
	rng     *dice.Roller
	options map[string]bool
	step1   *primaryStarDeterminator
}

type primaryStarDeterminator struct {
	tablesStarType    rtg.RandomTableGenerator
	activeMods        map[string]bool
	tablesStarSubType rtg.RandomTableGenerator
	completed         bool
}

type BuildOption func(*Builder)

func NewBuilder(seed string, options ...BuildOption) (*Builder, error) {
	b := Builder{}
	b.rng = dice.New(seed)
	b.options = make(map[string]bool)

	pssd := primaryStarDeterminator{}
	rtg1, err := rtg.NewStarTypeDeterminationGenerator(b.rng)
	if err != nil {
		return nil, fmt.Errorf("failed to create RNG: Star Type Determination: %v", err)
	}

	pssd.tablesStarType = rtg1

	b.step1 = &pssd
	return &b, nil
}

const (
	undecided         = "Special/Unusual"
	classIV_MTypeBump = 5
)

func (b *Builder) Build(directives ...string) (*StarSystem, error) {

	ss := &StarSystem{}
	ss.PrimaryStar = &Star{Designation: Primary}
	//step 1
	if err := b.determinePrimaryStarTypeAndClass(ss); err != nil {
		return nil, fmt.Errorf("failed to determine primary star type and class: %v", err)
	}

	//step 1a
	if err := b.determinePrimaryStarSubtype(ss); err != nil {
		return nil, fmt.Errorf("failed to determine primary star subtype: %v", err)
	}
	ss.Dead = ss.PrimaryStar.Dead
	ss.Primordial = ss.PrimaryStar.Protostar

	//step 1b
	if err := determineMassDiameterAgeTemperature(b.rng, ss.PrimaryStar); err != nil {
		return nil, err
	}

	//step 1c
	ss.PrimaryStar.Luminocity = float.Round(luminocity(ss.PrimaryStar.Diameter, ss.PrimaryStar.Temperature))

	//step 1d
	ss.Age = ss.PrimaryStar.Age

	return ss, nil
}

func (b *Builder) determinePrimaryStarTypeAndClass(ss *StarSystem) error {
	activeMods1 := []string{}
primary_star_class_generation:
	for {
		i, res, err := b.step1.tablesStarType.Roll("Type")
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		switch res {
		case "O", "B", "A", "F", "G", "K", "M":
			switch ss.PrimaryStar.Class {
			case "":
				ss.PrimaryStar.Class = "V"
			case "IV":
				if res == "M" {
					res, err = b.step1.tablesStarType.FindByIndex("Type", i+classIV_MTypeBump)
				}
				if res == "O" {
					res = "B"
				}
			case "VI":
				if res == "F" {
					res = "G"
				}
				if res == "A" {
					res = "B"
				}
			}
			if ss.Primordial {
				switch res {
				case "O", "B":
					continue
				}
			}
			ss.PrimaryStar.Type = res
		case "Ia", "Ib", "II", "III", "IV", "VI":
			if res != "IV" && res != "VI" {
				activeMods1 = append(activeMods1, rtg.MOD_NonMainSequenceClass)
			}
			ss.PrimaryStar.Class = res
		case "BD":
			ss.PrimaryStar.Class = res
			ss.Empty = true
			break primary_star_class_generation
		case "D", "NS", "BH":
			ss.PrimaryStar.Type = res
			ss.PrimaryStar.Dead = true
			ss.Empty = true
			break primary_star_class_generation
		case "Nb":
			if ss.NebulaType == 0 {
				ss.NebulaType = rollNebula(b.rng)
			}
		case "Star Cluster":
			ss.Clustered = true
		case "Protostar":
			ss.PrimaryStar.Protostar = true
			activeMods1 = append(activeMods1, rtg.MOD_ProtostarSystem)
		case "PSR":
			ss.PrimaryStar.Type = res
			ss.Dead = true
			break primary_star_class_generation
		case "Anomaly":
			ss.PrimaryStar.Type = res
			break primary_star_class_generation
		default:
			panic(fmt.Sprintf("dev error: invalid value rolled: %v", res))
		}
		if ss.PrimaryStar.Class != "" && ss.PrimaryStar.Type != "" {
			break primary_star_class_generation
		}
	}
	return nil
}

func (b *Builder) determinePrimaryStarSubtype(ss *StarSystem) error {
	switch ss.PrimaryStar.Type {
	case "M":
		_, res, err := b.step1.tablesStarType.Roll("M Type Primary")
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		ss.PrimaryStar.SubType = res
	case "O", "B", "A", "F", "G", "K":
		_, res, err := b.step1.tablesStarType.Roll("M Type Primary")
		if err != nil {
			return fmt.Errorf("failed to roll on RTG1: %v", err)
		}
		n, err := strconv.Atoi(res)
		if err != nil {
			return fmt.Errorf("expect number for subtype: '%v'", res)
		}
		if ss.PrimaryStar.Class == "IV" && ss.PrimaryStar.Type == "K" && n > 4 {
			n = n - 5 //For a K-type Class IV star, subtract 5 (make lower) any subtype result above 4 (p. 16)
		}
		ss.PrimaryStar.SubType = fmt.Sprintf("%v", n)
	}
	return nil
}

func determineMassDiameterAgeTemperature(r *dice.Roller, s *Star) error {
	i := interpolate.Index(s.Type, s.SubType, s.Class)
	switch s.Type {
	case "D":
		s.Mass = whiteDwarfMass(r)
		s.Diameter = float.Round(whiteDwarfDiameter(s.Mass))
		s.Age = starAge(r, s)
	case "BH":
		s.Mass = blackHoleMass(r)
		s.Diameter = float.Round(2.95 * s.Mass) //km
		s.Age = starAge(r, s)
	case "NS", "PSR":
		s.Mass = neitronStarMass(r)
		s.Diameter = 19 + float64(r.Roll("1d6")) //km
		s.Age = starAge(r, s)
	case "Anomaly":
	default:
		s.Mass = float.Round(interpolate.MassByIndex(i))
		if s.Mass == 0 {
			fmt.Println(i)
			return fmt.Errorf("failed to determine by interpolation: star mass (%v%v %v)", s.Type, s.SubType, s.Class)
		}
		s.Diameter = float.Round(interpolate.DiamByIndex(i))
		if s.Diameter == 0 {
			return fmt.Errorf("failed to determine by interpolation: star diameter (%v%v %v)", s.Type, s.SubType, s.Class)
		}
		s.Age = starAge(r, s)
	}
	s.Temperature = float.Round(interpolate.TempByIndex(i))
	return nil
}

func whiteDwarfMass(r *dice.Roller) float64 {
	r1 := float64(r.Roll("2d6"))
	r2 := float64(r.Roll("1d10"))
	m := float.Round(((r1 - 1) / 10) + (r2 / 100))
	if m > 1.44 {
		m = 1.34 + float64(r.Roll("1d100"))/1000
	}
	return m
}

func whiteDwarfDiameter(m float64) float64 {
	return (1.0 / m) * 0.01
}

func blackHoleMass(r *dice.Roller) float64 {
	r6 := r.Roll("1d6")
	r10 := r.Roll("1d10")
	m := 2.1 + float64(r6) - 1 + (float64(r10) / 10)
	for r6 == 6 {
		r6 = r.Roll("1d6")
		m += float64(r6)
	}
	return float.Round(m)
}

func neitronStarMass(r *dice.Roller) float64 {
	r1 := r.Roll("1d6")
	m := 1 + (float64(r1) / 10)
	for r1 == 6 {
		r1 = r.Roll("1d6")
		m += (float64(r1) - 1.0) / 10.0
	}
	return m
}

func variance(r *dice.Roller) float64 {
	r1 := r.Roll("1d1001-1")
	return float64(r1) / 1000.0
}

func mainSequanceLifespan(m float64) float64 {
	return float.Round(10 / math.Pow(m, 2.5))
}

func smallStarAge(r *dice.Roller) float64 {
	return float64(r.Roll("1d6x2")-r.Roll("1d3-2")) + float64(r.Roll("1d10")/10.0)
}

func subgiantLifeSpan(msl, m float64) float64 {
	return msl / (4 + m)
}

func giantLifeSpan(msl, m float64) float64 {
	return msl / 10.0 * math.Pow(m, 3)
}

func deadStarMass(r *dice.Roller, m float64) float64 {
	return float.Round(float64(r.Roll("1d3+2")) * m)
}

func starAge(r *dice.Roller, s *Star) float64 {
	age := 0.0
	mass := s.Mass
	if s.Dead {
		mass = deadStarMass(r, mass)
	}
	msl := mainSequanceLifespan(mass)
	if mass < 0.9 {
		msl = smallStarAge(r)
	}
	age = msl * variance(r)
	switch s.Class {
	case "BD":
		age = smallStarAge(r)
	case "D", "NS", "BH":
		mass := deadStarMass(r, s.Mass)
		msl := mainSequanceLifespan(deadStarMass(r, mass))
		age = smallStarAge(r) + msl + subgiantLifeSpan(msl, mass) + giantLifeSpan(msl, mass)
	case "PSR":
		mass := deadStarMass(r, s.Mass)
		msl := mainSequanceLifespan(deadStarMass(r, mass))
		age = (0.1 * float64(r.Roll("2d10"))) + msl + subgiantLifeSpan(msl, mass) + giantLifeSpan(msl, mass)
	case "IV":
		age = msl + (subgiantLifeSpan(msl, mass) * variance(r))
	case "III":
		age = msl + subgiantLifeSpan(msl, mass) + (giantLifeSpan(msl, mass) * variance(r))

	}
	if s.Protostar {
		age = float64(r.Roll("2d10")) * 0.01
	}
	age = min(13.8, float.Round(age))
	for age <= 0 {
		age = variance(r) * 13.8
	}
	return age
}

func luminocity(diameter, temperature float64) float64 {
	return math.Pow(diameter, 2) * math.Pow(temperature/float64(5772), 4)
}
