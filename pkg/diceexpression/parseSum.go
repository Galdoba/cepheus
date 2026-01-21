package diceexpression

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseBaseSum(s string) (int, int, error) {
	n, f := -1, -1
	re := regexp.MustCompile(`^\A(\d+)d+(\d+)`)
	normal := strings.TrimSpace(strings.ToLower(s))
	found := re.FindStringSubmatch(normal)
	if len(found) != 3 {
		return n, f, fmt.Errorf("expect exactly 3 submatched from '%s'", s)
	}
	v, err := strconv.Atoi(found[1])
	if err != nil {
		return n, f, fmt.Errorf("failed to parse number of dice (%v): %v", s, err)
	}
	n = v
	v, err = strconv.Atoi(found[2])
	if err != nil {
		return n, f, fmt.Errorf("failed to parse dice faces (%v): %v", s, err)
	}
	f = v
	return n, f, nil
}

func parseAdditiveMod(s string) (int, error) {
	sum := 0
	re := regexp.MustCompile(`[\+-]\d+`)
	normal := strings.TrimSpace(strings.ToLower(s))
	found := re.FindAllString(normal, -1)
	for i := range found {
		v, err := strconv.Atoi(found[i])
		if err != nil {
			return 0, fmt.Errorf("failed to parse additive mods (%v): %v", s, err)
		}
		sum += v

	}
	return sum, nil
}

func parseMultiplicativeMod(s string) (int, error) {
	mult := 1
	re := regexp.MustCompile(`x\d+`)
	normal := strings.TrimSpace(strings.ToLower(s))
	found := re.FindAllString(normal, -1)
	fmt.Println(found)
	for i := range found {
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "x"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse multiplicative mods (%v): %v", s, err)
		}
		mult = mult * v

	}
	return mult, nil
}

func parseDeletiveMod(s string) (int, error) {
	mult := 1
	re := regexp.MustCompile(`\/\d+`)
	normal := strings.TrimSpace(strings.ToLower(s))
	found := re.FindAllString(normal, -1)
	for i := range found {
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], `/`))
		if err != nil {
			return 0, fmt.Errorf("failed to parse deletive mods (%v): %v", s, err)
		}
		mult = mult * v

	}
	return mult, nil
}
func parseReplacements(s string) (map[int]int, error) {
	repl := make(map[int]int)
	re := regexp.MustCompile(`r\d+:\d+`)
	normal := strings.TrimSpace(strings.ToLower(s))
	found := re.FindAllString(normal, -1)
	for i := range found {
		parts := strings.Split(found[i], ":")
		parts[0] = strings.TrimPrefix(parts[0], "r")
		o, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse replacement mods (%v): %v", s, err)
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse replacement mods (%v): %v", s, err)
		}
		if _, ok := repl[o]; ok {
			return nil, fmt.Errorf("double asignment to replace value %v", o)
		}
		repl[o] = n

	}
	return repl, nil
}

func parseDropLow(s string) (int, error) {
	dl := 0
	set := false
	re := regexp.MustCompile(`dl\d+`)
	normal := strings.TrimSpace(strings.ToLower(s))
	found := re.FindAllString(normal, -1)
	fmt.Println(found)
	for i := range found {
		if set {
			return 0, fmt.Errorf("multiple directives to drop low dices: %v", s)
		}
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "dl"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse drop low value (%v): %v", found[i], err)
		}
		dl = v
		set = true
	}
	return dl, nil
}

func parseDropHigh(s string) (int, error) {
	dh := 0
	set := false
	re := regexp.MustCompile(`dh\d+`)
	found := re.FindAllString(s, -1)
	fmt.Println(found)
	for i := range found {
		if set {
			return 0, fmt.Errorf("multiple directives to drop high dices: %v", s)
		}
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "dh"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse drop high value (%v): %v", found[i], err)
		}
		dh = v
		set = true
	}
	return dh, nil
}
func parseIndividualMod(s string) (int, error) {
	individualMod := 0
	re := regexp.MustCompile(`i[\+-]\d+`)
	found := re.FindAllString(s, -1)
	fmt.Println(found)
	for i := range found {
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "i"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse drop high value (%v): %v", found[i], err)
		}
		individualMod += v
	}
	return individualMod, nil
}
func parseBottom(s string) (int, error) {
	m := 0
	set := false
	re := regexp.MustCompile(`b[\+-]\d+`)
	found := re.FindAllString(s, -1)
	fmt.Println(found)
	for i := range found {
		if set {
			return 0, fmt.Errorf("multiple derectives: bottom '%v'", s)
		}
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "min"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse drop high value (%v): %v", found[i], err)
		}
		m = v
		set = true
	}
	return m, nil
}
func parseTop(s string) (int, error) {
	m := 0
	set := false
	re := regexp.MustCompile(`t[\+-]\d+`)
	found := re.FindAllString(s, -1)
	fmt.Println(found)
	for i := range found {
		if set {
			return 0, fmt.Errorf("multiple derectives: top '%v'", s)
		}
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "min"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse drop high value (%v): %v", found[i], err)
		}
		m = v
		set = true
	}
	return m, nil
}
func parseTop(s string) (int, error) {
	m := 0
	set := false
	re := regexp.MustCompile(`t[\+-]\d+`)
	found := re.FindAllString(s, -1)
	fmt.Println(found)
	for i := range found {
		if set {
			return 0, fmt.Errorf("multiple derectives: top '%v'", s)
		}
		v, err := strconv.Atoi(strings.TrimPrefix(found[i], "min"))
		if err != nil {
			return 0, fmt.Errorf("failed to parse drop high value (%v): %v", found[i], err)
		}
		m = v
		set = true
	}
	return m, nil
}
