package dice

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type roller struct {
	src      rand.Source
	rand     *rand.Rand
	dicecode string
}

func newRoller(seed int64) *roller {
	r := roller{}
	r.src = rand.NewSource(seed)
	r.rand = rand.New(r.src)
	return &r
}

type rollRule struct {
	diceNum   int
	diceEdges int
}

func (r *roller) setCode(code string) {
	r.dicecode = code
}

func parseDiceCode(code string) rollRule {
	rr := rollRule{}
	rr.diceNum = -1
	code = strings.ToUpper(code)
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "FLUX", "2d6-7")
	re := regexp.MustCompile(`(\d*)([D]+)(\d*)`)
	subs := re.FindStringSubmatch(code)
	rr = parse(subs)
	if rr.diceEdges < 1 || rr.diceNum < 1 {
		panic("faild to parse: " + code)
	}
	return rr
}

/*
1D
1D+1
D3
D3+1
1 - ошибка
2 - 1D или D3
3 = 2d6
4 = D3+1
5 = 2d6+2
*/

func parse(subs []string) rollRule {
	rr := rollRule{}
	if !isDiceTag(subs[2]) {
		return rr
	}
	n, err := strconv.Atoi(subs[1])
	if err != nil {
		n = 1
	}
	rr.diceNum = n
	n, err = strconv.Atoi(subs[3])
	if err != nil {
		n = 6
	}
	rr.diceEdges = n
	return rr
}

func isDiceTag(tag string) bool {
	if tag == "D" || tag == "DD" {
		return true
	}
	return false
}
func isExplosive(tag string) bool {
	if tag == "DD" {
		return true
	}
	return false
}

func isModded(tag string) bool {
	if tag == "+" || tag == "-" {
		return true
	}
	return false
}

/*
dice.Sum("2d6",	WithBoon(),	WithDMs(skill.Admin, charcteristic.Edu))
)
*/
