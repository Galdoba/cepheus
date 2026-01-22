package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// compileRegexps компилирует все регулярные выражения один раз
var (
	baseRe           = regexp.MustCompile(`^(\d+)d(\d+)`)
	additiveRe       = regexp.MustCompile(`[\+-]\d+`)
	multiplicativeRe = regexp.MustCompile(`x[\+-]?\d+`)
	deletiveRe       = regexp.MustCompile(`/[\+-]?\d+`)
	replaceRe        = regexp.MustCompile(`r(\d+(?:;\d+)*):(\d+)`)
	dropLowRe        = regexp.MustCompile(`dl(\d+)`)
	dropHighRe       = regexp.MustCompile(`dh(\d+)`)
	individualRe     = regexp.MustCompile(`i([\+-]?\d+)`)
	sumMinRe         = regexp.MustCompile(`min([\+-]?\d+)`)
	sumMaxRe         = regexp.MustCompile(`max([\+-]?\d+)`)
	rerollRe         = regexp.MustCompile(`rr(\d+(?:;\d+)*)`)
)

// parseSumString однопроходный парсер для dice expression
func parseSumString(s string) (SumDirectives, error) {
	sd := newSumDirectives()
	remaining := strings.ToLower(strings.TrimSpace(s))
	// 1. Парсим базовую часть (обязательная)
	baseMatch := baseRe.FindStringSubmatch(remaining)
	if baseMatch == nil {
		return sd, fmt.Errorf("failed to parse base Sum")
	}

	// Извлекаем и удаляем базовую часть
	num, err := strconv.Atoi(baseMatch[1])
	if err != nil {
		return sd, fmt.Errorf("failed to parse number of dice: %v", err)
	}
	sd.Num = num

	faces, err := strconv.Atoi(baseMatch[2])
	if err != nil {
		return sd, fmt.Errorf("failed to parse dice faces: %v", err)
	}
	sd.Faces = faces

	remaining = strings.TrimPrefix(remaining, baseMatch[0])

	// Функция для удаления найденной подстроки
	removeFound := func(match string) string {
		// Находим позицию совпадения в оставшейся строке
		if idx := strings.Index(remaining, match); idx != -1 {
			// Удаляем совпадение
			remaining = remaining[:idx] + remaining[idx+len(match):]
		}
		return remaining
	}

	// 2. Парсим reroll значения (rr) - ДО замен (r), чтобы избежать конфликта
	for {
		match := rerollRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}
		// Разбираем значения, разделенные точкой с запятой
		valuesStr := match[1]
		values := strings.Split(valuesStr, ";")

		for _, valStr := range values {
			value, err := strconv.Atoi(valStr)
			if err != nil {
				return sd, fmt.Errorf("failed to parse reroll value '%s': %v", valStr, err)
			}
			if _, exists := sd.ReRoll[value]; exists {
				return sd, fmt.Errorf("duplicated reroll value %v", value)
			}
			sd.ReRoll[value] = true
		}

		remaining = removeFound(match[0])
	}

	// 3. Парсим замены (могут быть несколько)
	for {
		match := replaceRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}
		// Разбираем исходные значения, разделенные точкой с запятой
		sourcesStr := match[1]
		sources := strings.Split(sourcesStr, ";")

		replacement, err := strconv.Atoi(match[2])
		if err != nil {
			return sd, fmt.Errorf("failed to parse replacement value: %v", err)
		}

		for _, srcStr := range sources {
			original, err := strconv.Atoi(srcStr)
			if err != nil {
				return sd, fmt.Errorf("failed to parse replacement original: %v", err)
			}
			if _, exists := sd.Replace[original]; exists {
				return sd, fmt.Errorf("double assignment to replace value %v", original)
			}
			sd.Replace[original] = replacement
		}

		remaining = removeFound(match[0])
	}

	// 4. Парсим drop low (только один)
	if match := dropLowRe.FindStringSubmatch(remaining); match != nil {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return sd, fmt.Errorf("failed to parse drop low value: %v", err)
		}
		sd.SumMods[DropLow] = value
		remaining = removeFound(match[0])
	}

	// 5. Парсим drop high (только один)
	if match := dropHighRe.FindStringSubmatch(remaining); match != nil {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return sd, fmt.Errorf("failed to parse drop high value: %v", err)
		}
		sd.SumMods[DropHigh] = value
		remaining = removeFound(match[0])
	}

	// 6. Парсим индивидуальные модификаторы (могут быть несколько)
	individualSum := 0
	for {
		match := individualRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}
		// Парсим значение (может быть со знаком или без)
		valueStr := match[1]
		// Если нет знака, считаем положительным
		if valueStr[0] != '+' && valueStr[0] != '-' {
			valueStr = "+" + valueStr
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return sd, fmt.Errorf("failed to parse individual mod: %v", err)
		}
		individualSum += value
		remaining = removeFound(match[0])
	}
	sd.SumMods[Individual] = individualSum

	// 7. Парсим минимальную сумму (только один)
	if match := sumMinRe.FindStringSubmatch(remaining); match != nil {
		valueStr := match[1]
		// Если нет знака, считаем положительным
		if valueStr[0] != '+' && valueStr[0] != '-' {
			valueStr = "+" + valueStr
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return sd, fmt.Errorf("failed to parse minimum sum: %v", err)
		}
		sd.SumMods[SumMininum] = value
		remaining = removeFound(match[0])
	}

	// 8. Парсим максимальную сумму (только один)
	if match := sumMaxRe.FindStringSubmatch(remaining); match != nil {
		valueStr := match[1]
		// Если нет знака, считаем положительным
		if valueStr[0] != '+' && valueStr[0] != '-' {
			valueStr = "+" + valueStr
		}
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return sd, fmt.Errorf("failed to parse maximum sum: %v", err)
		}
		sd.SumMods[SumMaximum] = value
		remaining = removeFound(match[0])
	}

	// 9. Собираем мультипликативные модификаторы (могут быть несколько)
	multiplicativeProd := 1
	for {
		match := multiplicativeRe.FindString(remaining)
		if match == "" {
			break
		}
		value, err := strconv.Atoi(strings.TrimPrefix(match, "x"))
		if err != nil {
			return sd, fmt.Errorf("failed to parse multiplicative mod: %v", err)
		}
		multiplicativeProd *= value
		remaining = removeFound(match)
	}
	sd.SumMods[Multiplicative] = multiplicativeProd

	// 10. Собираем делительные модификаторы (могут быть несколько)
	deletiveProd := 1
	for {
		match := deletiveRe.FindString(remaining)
		if match == "" {
			break
		}
		value, err := strconv.Atoi(strings.TrimPrefix(match, "/"))
		if err != nil {
			return sd, fmt.Errorf("failed to parse deletive mod: %v", err)
		}
		deletiveProd *= value
		remaining = removeFound(match)
	}
	sd.SumMods[Deletive] = deletiveProd

	// 11. Собираем аддитивные модификаторы (в самом конце, чтобы избежать конфликтов)
	additiveSum := 0
	for {
		match := additiveRe.FindString(remaining)
		if match == "" {
			break
		}
		value, err := strconv.Atoi(match)
		if err != nil {
			return sd, fmt.Errorf("failed to parse additive mod: %v", err)
		}
		additiveSum += value
		remaining = removeFound(match)
	}
	sd.SumMods[Additive] = additiveSum

	// Проверяем, что вся строка была обработана
	if strings.TrimSpace(remaining) != "" {
		return sd, fmt.Errorf("unrecognized tokens in expression: %s", remaining)
	}

	// Валидация значений
	if sd.Num <= 0 {
		return sd, fmt.Errorf("number of dice must be positive, got %d", sd.Num)
	}
	if sd.Faces <= 0 {
		return sd, fmt.Errorf("number of faces must be positive, got %d", sd.Faces)
	}
	if dl, ok := sd.SumMods[DropLow]; ok && dl > 0 {
		if dl >= sd.Num {
			return sd, fmt.Errorf("drop low count (%d) must be less than number of dice (%d)", dl, sd.Num)
		}
	}
	if dh, ok := sd.SumMods[DropHigh]; ok && dh > 0 {
		if dh >= sd.Num {
			return sd, fmt.Errorf("drop high count (%d) must be less than number of dice (%d)", dh, sd.Num)
		}
	}
	if sd.SumMods[Multiplicative] == 1 {
		delete(sd.SumMods, Multiplicative)
	}
	if sd.SumMods[Deletive] == 1 {
		delete(sd.SumMods, Deletive)
	}
	if sd.SumMods[Individual] == 0 {
		delete(sd.SumMods, Individual)
	}
	if sd.SumMods[Additive] == 0 {
		delete(sd.SumMods, Additive)
	}

	return sd, nil
}
