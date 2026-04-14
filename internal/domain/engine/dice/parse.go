package dice

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	diceTypeNormal      = "d"
	diceTypeConcat      = "D"
	diceTypeDestructive = "DD"
	addEachSuffix       = "e"
	addIndividualSuffix = ">>"
	dropLowPrefix       = "dl"
	dropHighPrefix      = "dh"
	dividePrefix        = "/"
	multiplyPrefix1     = "x"
	multiplyPrefix2     = "*"

	complexModsSeparator = ":"
)

func parseExpression(expr string) (Dicepool, []Mod, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return Dicepool{}, nil, fmt.Errorf("empty expression")
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
		return Dicepool{}, nil, err
	}

	var simpleAdd *AddConst
	if strings.TrimSpace(leftover) != "" {
		val, err := parseSimpleAdditive(leftover)
		if err != nil {
			return Dicepool{}, nil, fmt.Errorf("invalid simple additive modifier '%s': %w", leftover, err)
		}
		simpleAdd = &AddConst{value: val}
	}

	dice := make([]Die, count)
	for i := range count {
		switch diceType {
		case diceTypeNormal:
			dice[i] = NewDice(sides)
		case diceTypeConcat, diceTypeDestructive:
			//TODO:special types (not implemented yet)
			dice[i] = NewDice(0).WithMeta(map[string]string{
				"special": diceType,
				"faces":   strconv.Itoa(sides),
			})
		default:
			return Dicepool{}, nil, fmt.Errorf("unknown dice type %s", diceType)
		}
	}

	complexMods, err := parseComplexModifiers(modsPart)
	if err != nil {
		return Dicepool{}, nil, err
	}

	allMods := []Mod{}
	allMods = append(allMods, complexMods...)
	if simpleAdd != nil {
		allMods = append(allMods, *simpleAdd)
	}
	// Добавляем Sum, если ещё нет агрегатора (пока всегда)
	allMods = append(allMods, Sum{})

	allMods = sortModifiers(allMods)

	dp := NewDicepool(dice...).WithMods(allMods...)
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

	// в) Грани (целое положительное число)
	if len(part) == 0 || part[0] < '0' || part[0] > '9' {
		return 0, "", 0, "", fmt.Errorf("missing sides after %s", diceType)
	}
	i := 0
	for i < len(part) && part[i] >= '0' && part[i] <= '9' {
		i++
	}
	sides, err = strconv.Atoi(part[:i])
	if err != nil || sides < 2 {
		return 0, "", 0, "", fmt.Errorf("invalid sides: %s", part[:i])
	}
	part = part[i:]

	leftover = strings.TrimSpace(part)
	// Проверяем, что leftover состоит только из +, -, цифр, пробелов
	for _, ch := range leftover {
		if ch != '+' && ch != '-' && ch != ' ' && (ch < '0' || ch > '9') {
			return 0, "", 0, "", fmt.Errorf("unexpected characters in additive part: %q", leftover)
		}
	}
	return count, diceType, sides, leftover, nil
}

// parseSimpleAdditive parse additives like "+3", "-2", "  +5  "
func parseSimpleAdditive(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// parseComplexModifiers разбирает последовательность модификаторов, разделённых ':'
func parseComplexModifiers(modsStr string) ([]Mod, error) {
	if modsStr == "" {
		return nil, nil
	}
	parts := strings.Split(modsStr, complexModsSeparator)
	var mods []Mod
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

func parseOneComplexModifier(modStr string) (Mod, error) {
	// +Ne / -Ne
	if strings.HasSuffix(modStr, addEachSuffix) {
		numStr := strings.TrimSuffix(modStr, addEachSuffix)
		val, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AddToEach: %s", modStr)
		}
		return AddToEach{value: val}, nil
	}
	// +NtoM / -NtoM
	if strings.Contains(modStr, addIndividualSuffix) {
		// формат: +3to2 или -5to1
		var sign int = 1
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
		return AddIndividual{position: pos, value: sign * val}, nil
	}
	// dlN
	if strings.HasPrefix(modStr, dropLowPrefix) {
		numStr := modStr[2:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid dl count: %s", numStr)
		}
		return DropLowest{quantity: n}, nil
	}
	// dhN
	if strings.HasPrefix(modStr, dropHighPrefix) {
		numStr := modStr[2:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid dh count: %s", numStr)
		}
		return DropHighest{quantity: n}, nil
	}
	// /N
	if strings.HasPrefix(modStr, dividePrefix) {
		numStr := modStr[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n == 0 {
			return nil, fmt.Errorf("invalid divisor: %s", numStr)
		}
		return Divide{value: n}, nil
	}
	// xN или *N
	if strings.HasPrefix(modStr, multiplyPrefix1) || strings.HasPrefix(modStr, multiplyPrefix2) {
		numStr := modStr[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid multiplier: %s", numStr)
		}
		return Multiply{value: n}, nil
	}
	return nil, fmt.Errorf("unknown complex modifier: %s", modStr)
}

func sortModifiers(mods []Mod) []Mod {
	sort.SliceStable(mods, func(i, j int) bool {
		return mods[i].Priority() < mods[j].Priority()
	})
	return mods
}

func ValidateExpression(expr string) error {
	_, _, err := parseExpression(expr)
	return err
}
