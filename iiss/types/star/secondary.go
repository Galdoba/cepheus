package star

import (
	"github.com/Galdoba/cepheus/pkg/dice"
)

func (s Star) Lesser(dp *dice.Dicepool) Star {
	second := Star{}
	primary := s.Index()
	primClass, primType, primSub := FromIndex(primary)
	stellar := ""
	for stellar == "" {
		second, _ = Generate(dp, KnownClass(primClass), KnownType(lesserType(primType)))
		sub, _ := StarSubTypeDetermination(dp, second)

		if sub != nil {
			if primSub > *sub {
				continue
			}
			second.SubType = sub
		}
		if second.Type == "M" && s.Type == "M" && *second.SubType > *s.SubType {
			second, _ = Generate(dp, KnownStellar("BD"))
		}
		stellar = second.ToStellar()

	}
	return second
}

func (s Star) Random(dp *dice.Dicepool) Star {
	second, _ := Generate(dp)
	if s.Mass < second.Mass {
		return s.Lesser(dp)
	}
	return second
}

func (s Star) Sibling(dp *dice.Dicepool) Star {
	second, _ := Generate(dp, KnownStellar(s.ToStellar()))
	sub := *second.SubType
	sub += dp.Sum1D()
	tp := s.Type
	newStellar := ""
	if sub >= 10 {
		sub -= 10
		second.Type = lesserType(tp)
		second.SubType = &sub
		newStellar = second.ToStellar()
		if s.Type == "M" {
			newStellar = "BD"
		}
	} else {
		second.SubType = &sub
		newStellar = second.ToStellar()
	}
	second, _ = Generate(dp, KnownStellar(newStellar))
	return second
}

func (s Star) Twin(dp *dice.Dicepool) Star {
	second, _ := Generate(dp, KnownStellar(s.ToStellar()))
	return second
}

func lesserType(tp string) string {
	switch tp {
	case "BH":
		return "NS"
	case "NS", "PSR":
		return "D"
	case "D", "BD", "L", "T", "Y":
		return "BD"
	case "O":
		return "B"
	case "B":
		return "A"
	case "A":
		return "F"
	case "F":
		return "G"
	case "G":
		return "K"
	case "K", "M":
		return "M"
	}
	return "?"
}
