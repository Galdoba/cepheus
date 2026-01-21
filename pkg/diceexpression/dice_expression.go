package diceexpression

import (
	"fmt"
	"regexp"
	"strings"
)

type DiceExpression string

//Scheme Sum: NdM[+/-n][xn][ro:n][dlN][dhN][rrN][maxS][minS]

const (
	Concat         RollType = "concat"
	Sum            RollType = "sum"
	Values         RollType = "values"
	Additive       ModType  = "additive"
	Multiplicative ModType  = "multiplicative"
	Deletive       ModType  = "deletive"
)

type RollType string

type ModType string

type SumDirectives struct {
	Num     int
	Faces   int
	SumMods map[ModType]int
	Replace map[int]int
}

func newSumDirectives() *SumDirectives {
	sd := SumDirectives{}
	sd.SumMods = make(map[ModType]int)
	sd.SumMods[Multiplicative] = 1
	sd.Replace = make(map[int]int)
	return &sd
}

type RollDirectives struct {
	Origin   string
	RollType RollType
	SumD     *SumDirectives
}

func (de DiceExpression) Parse() (RollDirectives, error) {
	rd := RollDirectives{}
	s := string(de)
	s = strings.TrimSpace(strings.ToLower(s))
	rd.RollType = parseType(s)
	switch rd.RollType {
	case Sum:
		rd.SumD = newSumDirectives()
		n, f, err := parseBaseSum(s)
		if err != nil {
			return rd, fmt.Errorf("failed to parse base Sum")
		}
		rd.SumD.Num = n
		rd.SumD.Faces = f
		n, err = parseAdditiveMod(s)
		if err != nil {
			return rd, fmt.Errorf("failed to parse additive mods (%v): %v", s, err)
		}
		rd.SumD.SumMods[Additive] = n
		n, err = parseMultiplicativeMod(s)
		if err != nil {
			return rd, fmt.Errorf("failed to parse additive mods (%v): %v", s, err)
		}
		rd.SumD.SumMods[Multiplicative] = n
	case Concat:
	default:
		return rd, fmt.Errorf("%v", string(de))
	}

	return rd, fmt.Errorf("incomplete")
}

func parseType(de string) RollType {
	re := regexp.MustCompile(`\A(\d+)d*`)
	found := re.FindString(string(de))
	if found != "" {
		return Sum
	}
	re = regexp.MustCompile(`\Ad*`)
	found = re.FindString(string(de))
	if found != "" {
		return Concat
	}
	return RollType(fmt.Sprintf("unknown type: %s", de))

}
