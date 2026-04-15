package dice

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	diceTypeNormal       = "d"
	diceTypeConcat       = "D"
	diceTypeDestructive  = "DD"
	addEachSuffix        = "e"
	addIndividualSuffix  = ">>"
	dropLowPrefix        = "dl"
	dropHighPrefix       = "dh"
	dividePrefix         = "/"
	multiplyPrefix1      = "x"
	multiplyPrefix2      = "*"
	complexModsSeparator = ":"
)

// parseExpression converts a dice expression string into a dicepool and a list of mods.
func parseExpression(expr string) (dicepool, []mod, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return dicepool{}, nil, fmt.Errorf("empty expression")
	}

	dicePart, modsPart := "", ""
	if before, after, ok := strings.Cut(expr, ":"); ok {
		dicePart = before
		modsPart = after
	} else {
		dicePart = expr
	}

	count, diceType, sides, leftover, err := parseDicePart(dicePart)
	if err != nil {
		return dicepool{}, nil, err
	}

	var simpleAdd *addConst
	if strings.TrimSpace(leftover) != "" {
		val, err := parseSimpleAdditive(leftover)
		if err != nil {
			return dicepool{}, nil, fmt.Errorf("invalid simple additive modifier '%s': %w", leftover, err)
		}
		simpleAdd = &addConst{value: val}
	}

	dice := make([]die, count)
	for i := range count {
		switch diceType {
		case diceTypeNormal:
			dice[i] = newDie(sides)
		case diceTypeConcat, diceTypeDestructive:
			return dicepool{}, nil, fmt.Errorf("special dice types not implemented")
		default:
			return dicepool{}, nil, fmt.Errorf("unknown dice type %s", diceType)
		}
	}

	complexMods, err := parseComplexModifiers(modsPart)
	if err != nil {
		return dicepool{}, nil, err
	}

	allMods := []mod{}
	allMods = append(allMods, complexMods...)
	if simpleAdd != nil {
		allMods = append(allMods, *simpleAdd)
	}
	allMods = append(allMods, summ{})
	allMods = sortModifiers(allMods)

	dp := newDicepool(dice...)
	return dp, allMods, nil
}

func parseDicePart(part string) (count int, diceType string, sides int, leftover string, err error) {
	count = 1
	if len(part) > 0 && part[0] >= '0' && part[0] <= '9' {
		i := 0
		for i < len(part) && part[i] >= '0' && part[i] <= '9' {
			i++
		}
		count, err = strconv.Atoi(part[:i])
		if err != nil {
			return 0, "", 0, "", fmt.Errorf("invalid count: %s", part[:i])
		}
		part = part[i:]
	}
	if len(part) == 0 {
		return 0, "", 0, "", fmt.Errorf("missing dice type")
	}

	switch {
	case strings.HasPrefix(part, diceTypeDestructive):
		diceType = diceTypeDestructive
	case strings.HasPrefix(part, diceTypeConcat):
		diceType = diceTypeConcat
	case strings.HasPrefix(part, diceTypeNormal):
		diceType = diceTypeNormal
	default:
		return 0, "", 0, "", fmt.Errorf("invalid dice type")
	}
	part = strings.TrimPrefix(part, diceType)

	if len(part) == 0 || part[0] < '0' || part[0] > '9' {
		return 0, "", 0, "", fmt.Errorf("missing sides after %s", diceType)
	}
	i := 0
	for i < len(part) && part[i] >= '0' && part[i] <= '9' {
		i++
	}
	sides, err = strconv.Atoi(part[:i])
	if err != nil || sides < 1 {
		return 0, "", 0, "", fmt.Errorf("invalid sides: %s", part[:i])
	}
	part = part[i:]

	leftover = strings.TrimSpace(part)
	for _, ch := range leftover {
		if ch != '+' && ch != '-' && ch != ' ' && (ch < '0' || ch > '9') {
			return 0, "", 0, "", fmt.Errorf("unexpected characters in additive part: %q", leftover)
		}
	}
	return count, diceType, sides, leftover, nil
}

func parseSimpleAdditive(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if s[0] == '+' {
		s = s[1:]
	}
	return strconv.Atoi(s)
}

func parseComplexModifiers(modsStr string) ([]mod, error) {
	if modsStr == "" {
		return nil, nil
	}
	parts := strings.Split(modsStr, complexModsSeparator)
	var mods []mod
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		m, err := parseOneComplexModifier(p)
		if err != nil {
			return nil, err
		}
		mods = append(mods, m)
	}
	return mods, nil
}

func parseOneComplexModifier(modStr string) (mod, error) {
	if before, ok := strings.CutSuffix(modStr, addEachSuffix); ok {
		numStr := before
		val, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AddToEach: %s", modStr)
		}
		return addToEach{value: val}, nil
	}
	if strings.Contains(modStr, addIndividualSuffix) {
		sign := 1
		s := modStr
		switch s[0] {
		case '-':
			sign = -1
			s = s[1:]
		case '+':
			s = s[1:]
		}
		parts := strings.Split(s, addIndividualSuffix)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid AddIndividual format: %s", modStr)
		}
		val, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		pos, err := strconv.Atoi(parts[1])
		if err != nil || pos < 1 {
			return nil, fmt.Errorf("invalid position: %s", parts[1])
		}
		return addIndividual{position: pos, value: sign * val}, nil
	}
	if strings.HasPrefix(modStr, dropLowPrefix) {
		numStr := modStr[2:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid dl count: %s", numStr)
		}
		return dropLowest{quantity: n}, nil
	}
	if strings.HasPrefix(modStr, dropHighPrefix) {
		numStr := modStr[2:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid dh count: %s", numStr)
		}
		return dropHighest{quantity: n}, nil
	}
	if strings.HasPrefix(modStr, dividePrefix) {
		numStr := modStr[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n == 0 {
			return nil, fmt.Errorf("invalid divisor: %s", numStr)
		}
		return divide{value: n}, nil
	}
	if strings.HasPrefix(modStr, multiplyPrefix1) || strings.HasPrefix(modStr, multiplyPrefix2) {
		numStr := modStr[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid multiplier: %s", numStr)
		}
		return multiply{value: n}, nil
	}
	return nil, fmt.Errorf("unknown complex modifier: %s", modStr)
}

func sortModifiers(mods []mod) []mod {
	sort.SliceStable(mods, func(i, j int) bool {
		return mods[i].priority() < mods[j].priority()
	})
	return mods
}

// ValidateExpression checks whether a dice expression is syntactically valid.
func ValidateExpression(expr string) error {
	_, _, err := parseExpression(expr)
	return err
}
