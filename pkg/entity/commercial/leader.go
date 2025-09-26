package commercial

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/names"
)

const (
	CareerAgent       = "Agent"
	CareerArmy        = "Army"
	CareerCitizen     = "Citizen"
	CareerDrifter     = "Drifter"
	CareerEntertainer = "Entertainer"
	CareerMarines     = "Marines"
	CareerMerchants   = "Merchants"
	CareerNavy        = "Navy"
	CareerNobility    = "Nobility"
	CareerRogue       = "Rogue"
	CareerScholar     = "Scholar"
	CareerScout       = "Scout"
)

type Leader struct {
	Name         string
	Founder      bool
	Career       map[string]int
	ShipShares   int
	Investiments int
}

func randomLeaders() []*Leader {
	lNum := dice.FastRandom("1d5")
	leaders := []*Leader{}
	for i := 0; i < lNum; i++ {
		leaders = append(leaders, NewRandomLeader(true))
	}
	return leaders
}

func NewRandomLeader(isFounder bool) *Leader {
	male := dice.FastRandom("1d6") > 3
	name := ""
	switch male {
	case true:
		name = names.RandomMaleName()
	case false:
		name = names.RandomFemaleName()
	}
	name += " " + names.RandomLastName()
	l := Leader{}
	l.Name = name
	l.Career = randomCareers()
	l.Founder = isFounder
	l.Investiments = max(0, (dice.FastRandom("2d6")-7)*25000)
	l.ShipShares = max(0, (dice.FastRandom("2d6") - 7))
	return &l
}

func randomCareers() map[string]int {
	careersMap := make(map[string]int)
	for i := 1; i <= 12; i++ {
		careersMap[careerIndexToString(i)] = 0
	}
	for {

		careerSum := sumCareers(careersMap)
		if careerSum > 6 {
			break
		}
		switch careerSum {
		case 0:
			careersMap[newCareer(careersMap)]++
		default:
			r := dice.FastRandom("2d6")
			switch r {
			case 2:
				break
			default:
				if r/2 < careerSum {
					return careersMap
				}
				if r >= 11 {
					careersMap[newCareer(careersMap)]++
					continue
				}
				careersMap[existantCareer(careersMap)]++
			}
		}
	}
	return careersMap
}

func sumCareers(cm map[string]int) int {
	s := 0
	for _, v := range cm {
		s += v
	}
	return s
}

func existantCareer(cm map[string]int) string {
	have := []string{}
	for k, v := range cm {
		if v > 0 {
			have = append(have, k)
		}
	}
	i := dice.FromSliceRandom(have)
	return have[i]
}

func newCareer(cm map[string]int) string {
	haveNot := []string{}
	for k, v := range cm {
		if v == 0 {
			haveNot = append(haveNot, k)
		}
	}
	i := dice.FromSliceRandom(haveNot)
	return haveNot[i]
}

func careerIndexToString(i int) string {
	switch i {
	case 1:
		return CareerAgent
	case 2:
		return CareerArmy
	case 3:
		return CareerCitizen
	case 4:
		return CareerDrifter
	case 5:
		return CareerEntertainer
	case 6:
		return CareerMarines
	case 7:
		return CareerMerchants
	case 8:
		return CareerNavy
	case 9:
		return CareerNobility
	case 10:
		return CareerRogue
	case 11:
		return CareerScholar
	case 12:
		return CareerScout
	default:
		panic("no index")
	}
}

func NewFounder(terms map[string]int) *Leader {
	l := Leader{
		Founder: true,
		Career:  terms,
	}
	return &l
}

func (l *Leader) WithShipShares(s int) *Leader {
	l.ShipShares = s
	return l
}

func (l *Leader) WithInvestiments(i int) *Leader {
	l.Investiments = i
	return l
}

func (l *Leader) IsFounder() bool {
	return l.Founder
}

func (l *Leader) ProvideSkills() []string {
	skills := []string{}
	for career, terms := range l.Career {
		for i := 0; i < terms; i++ {
			skillPool := aplicableSkills(career)
			index := dice.RandomIndex(skillPool)
			skills = append(skills, skillPool[index])
		}
	}

	return skills
}

func (l *Leader) Invested() (int, int) {
	return l.Investiments, l.ShipShares
}

type Founder interface {
	IsFounder() bool
	Invested() (int, int)
	ProvideSkills(*dice.Dicepool) []string
}

func aplicableSkills(career string) []string {
	switch career {
	case CareerAgent:
		return []string{SkillAdvocasy, SkillAgency, SkillMischief, SkillPropaganda}
	case CareerArmy:
		return []string{SkillAgency, SkillMischief, SkillResearch}
	case CareerCitizen:
		return []string{SkillAdvocasy, SkillAgency, SkillBrokerage, SkillFabrication, SkillInvestiment, SkillMischief, SkillPropaganda, SkillResearch, SkillShipping}
	case CareerDrifter:
		return []string{SkillAgency, SkillFabrication, SkillMischief, SkillResearch, SkillShipping}
	case CareerEntertainer:
		return []string{SkillNobility, SkillPropaganda}
	case CareerMarines:
		return []string{SkillAgency, SkillMischief, SkillResearch, SkillShipping}
	case CareerMerchants:
		return []string{SkillAdvocasy, SkillBrokerage, SkillInvestiment, SkillNobility, SkillPropaganda, SkillResearch, SkillShipping}
	case CareerNavy:
		return []string{SkillAgency, SkillMischief, SkillResearch, SkillShipping}
	case CareerNobility:
		return []string{SkillAdvocasy, SkillInvestiment, SkillNobility, SkillPropaganda}
	case CareerRogue:
		return []string{SkillAgency, SkillBrokerage, SkillInvestiment, SkillMischief, SkillPropaganda}
	case CareerScholar:
		return []string{SkillBrokerage, SkillInvestiment, SkillResearch}
	case CareerScout:
		return []string{SkillAgency, SkillBrokerage, SkillMischief, SkillPropaganda, SkillResearch}
	default:
		return nil
	}
}

func (l *Leader) String() string {
	s := l.Name
	careerStr := []string{}
	for i := 1; i <= 12; i++ {
		careerName := careerIndexToString(i)
		if val, ok := l.Career[careerName]; ok {
			if val == 0 {
				continue
			}
			careerStr = append(careerStr, fmt.Sprintf("%v (%v)", careerName, val))
		}
	}
	for len(s) < 26 {
		s += " "
	}
	s += "\t" + strings.Join(careerStr, ", ")
	if !l.Founder {
		s += fmt.Sprintf("[profit modifier=-%v", l.profitCut()) + `%]`
	}
	if sumCareers(l.Career) == 0 {
		s += "[shareholder]"
	}
	return s
}

func (l *Leader) profitCut() int {
	if l.Founder {
		return 0
	}
	sum := sumCareers(l.Career)
	switch sum {
	case 0:
		return 1
	default:
		return sum * 2
	}
}
