package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	concatBaseRe = regexp.MustCompile(`^d(\d+)`)
	concatModRe  = regexp.MustCompile(`cm(\d+):([+-]?\d+)`)
)

func parseConcatString(s string) (ConcatDirectives, error) {
	cd := ConcatDirectives{}
	remaining := strings.ToLower(strings.TrimSpace(s))

	// 1. Парсим базовую часть (обязательная) - dN
	baseMatch := concatBaseRe.FindStringSubmatch(remaining)
	if baseMatch == nil {
		return cd, fmt.Errorf("failed to parse base Concat: expected format dN")
	}

	// Извлекаем и удаляем базовую часть
	diceStr := baseMatch[1]

	// Разбиваем число на цифры - каждая цифра это грань кубика
	for _, char := range diceStr {
		face, err := strconv.Atoi(string(char))
		if err != nil {
			// Это не должно случиться, так как регулярка гарантирует цифры
			return cd, fmt.Errorf("invalid digit in dice specification: %v", char)
		}
		cd.Faces = append(cd.Faces, face)
	}

	// Инициализируем модификаторы нулями
	cd.Mods = make([]int, len(cd.Faces))

	remaining = strings.TrimPrefix(remaining, baseMatch[0])

	// Функция для удаления найденной подстроки
	removeFound := func(match string) string {
		if idx := strings.Index(remaining, match); idx != -1 {
			remaining = remaining[:idx] + remaining[idx+len(match):]
		}
		return remaining
	}

	// 2. Парсим модификаторы (cmA:X)
	modMap := make(map[int]int) // Для проверки дублирования

	for {
		match := concatModRe.FindStringSubmatch(remaining)
		if match == nil {
			break
		}

		// Парсим номер кубика (A)
		diceIndex, err := strconv.Atoi(match[1])
		if err != nil {
			return cd, fmt.Errorf("failed to parse dice index in modifier: %v", match[1])
		}

		// Проверяем что индекс в пределах
		if diceIndex < 1 || diceIndex > len(cd.Faces) {
			return cd, fmt.Errorf("dice index %d out of range (1-%d)", diceIndex, len(cd.Faces))
		}

		// Проверяем на дублирование
		if _, exists := modMap[diceIndex]; exists {
			return cd, fmt.Errorf("duplicate modifier for dice %d", diceIndex)
		}

		// Парсим значение модификатора (X)
		modValue, err := strconv.Atoi(match[2])
		if err != nil {
			return cd, fmt.Errorf("failed to parse modifier value: %v", match[2])
		}

		// Сохраняем модификатор (индекс 0-based)
		cd.Mods[diceIndex-1] = modValue
		modMap[diceIndex] = modValue

		remaining = removeFound(match[0])
	}

	// Проверяем, что вся строка была обработана
	if strings.TrimSpace(remaining) != "" {
		return cd, fmt.Errorf("unrecognized tokens in expression: %s", remaining)
	}

	return cd, nil
}
