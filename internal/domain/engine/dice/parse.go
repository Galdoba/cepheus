package dice

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func ParseExpression(expr string) (Dicepool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return Dicepool{}, fmt.Errorf("empty expression")
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
		return Dicepool{}, err
	}

	var simpleAdd *AddConst
	if strings.TrimSpace(leftover) != "" {
		val, err := parseSimpleAdditive(leftover)
		if err != nil {
			return Dicepool{}, fmt.Errorf("invalid simple additive modifier '%s': %w", leftover, err)
		}
		simpleAdd = &AddConst{value: val}
	}

	dice := make([]Die, count)
	for i := range count {
		switch diceType {
		case "d":
			dice[i] = NewDice(sides)
		case "D", "DD":
			// Специальные типы – пока только мета, бросок не реализован
			dice[i] = NewDice(0).WithMeta(map[string]string{
				"special": diceType,
				"faces":   strconv.Itoa(sides),
			})
		default:
			return Dicepool{}, fmt.Errorf("unknown dice type %s", diceType)
		}
	}

	complexMods, err := parseComplexModifiers(modsPart)
	if err != nil {
		return Dicepool{}, err
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
	return dp, nil
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
	case strings.HasPrefix(part, "DD"):
		diceType = "DD"
	case strings.HasPrefix(part, "D"):
		diceType = "D"
	case strings.HasPrefix(part, "d"):
		diceType = "d"
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

// parseSimpleAdditive парсит строку вида "+3", "-2", "  +5  "
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
	parts := strings.Split(modsStr, ":")
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
	if strings.HasSuffix(modStr, "e") {
		numStr := strings.TrimSuffix(modStr, "e")
		val, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid AddToEach: %s", modStr)
		}
		return AddToEach{value: val}, nil
	}
	// +NtoM / -NtoM
	if strings.Contains(modStr, "to") {
		// формат: +3to2 или -5to1
		var sign int = 1
		s := modStr
		if s[0] == '-' {
			sign = -1
			s = s[1:]
		} else if s[0] == '+' {
			s = s[1:]
		}
		parts := strings.Split(s, "to")
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
	if strings.HasPrefix(modStr, "dl") {
		numStr := modStr[2:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid dl count: %s", numStr)
		}
		return DropLowest{quantity: n}, nil
	}
	// dhN
	if strings.HasPrefix(modStr, "dh") {
		numStr := modStr[2:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid dh count: %s", numStr)
		}
		return DropHighest{quantity: n}, nil
	}
	// /N
	if strings.HasPrefix(modStr, "/") {
		numStr := modStr[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil || n == 0 {
			return nil, fmt.Errorf("invalid divisor: %s", numStr)
		}
		return Divide{value: n}, nil
	}
	// xN или *N
	if strings.HasPrefix(modStr, "x") || strings.HasPrefix(modStr, "*") {
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
