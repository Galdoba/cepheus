package star

import (
	"fmt"
	"sort"
	"strings"
)

type StarKey string

func defineStarKeys() map[StarKey]bool {
	stars := make(map[StarKey]bool)
	for _, stellarClass := range []string{"O", "B", "A", "F", "G", "K", "M", "L", "T", "Y", "D"} {
		for _, numericalSubclass := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"} {
			for _, luminocityClass := range []string{"Ia", "Ib", "II", "III", "IV", "V", "VI", "BD"} {
				key := fmt.Sprintf("%s%s %s", stellarClass, numericalSubclass, luminocityClass)
				switch stellarClass {
				case "L", "T", "Y", "D":
					key = fmt.Sprintf("%s%s", stellarClass, numericalSubclass)
					stars[StarKey(key)] = true
					continue
				case "O", "B", "A", "F":
					if luminocityClass == "VI" {
						key = strings.ReplaceAll(key, "VI", "V")
					}
				}
				if luminocityClass == "BD" {
					key = luminocityClass
					stars[StarKey(key)] = true
					continue
				}
				stars[StarKey(key)] = true
			}
		}
	}
	return stars
}

func Parse(s string) []StarKey {
	validKeys := defineStarKeys()

	keys := make([]string, 0, len(validKeys))
	for k := range validKeys {
		keys = append(keys, string(k))
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	result := make([]StarKey, 0)
	n := len(s)
	i := 0
	for i < n {
		found := false
		for _, key := range keys {
			l := len(key)
			if i+l <= n && s[i:i+l] == key {
				result = append(result, StarKey(key))
				i += l
				found = true
				break
			}
		}
		if !found {
			i++
		}
	}
	return result
}

// parseKey разбирает StarKey на составляющие:
// stellarClass (O,B,A,F,G,K,M,L,T,Y,D или пустая строка для BD),
// numericalSubclass (0-9 или пустая строка),
// luminocityClass (Ia,Ib,II,III,IV,V,VI,D,BD или пустая строка).
func parseKey(sk StarKey) (string, string, string) {
	s := string(sk)

	// Специальный случай: "BD" (коричневый карлик)
	if s == "BD" {
		return "", "", "BD"
	}

	if len(s) < 2 {
		return "", "", ""
	}

	// Первый символ должен быть буквой класса
	stellarClass := s[0:1]
	if !isValidStellarClass(stellarClass) {
		return "", "", ""
	}

	// Второй символ - цифра подкласса
	numericalSubclass := s[1:2]
	if numericalSubclass < "0" || numericalSubclass > "9" {
		return "", "", ""
	}

	// Если длина ровно 2 -> формат "Xy" (L,T,Y,D)
	if len(s) == 2 {
		return stellarClass, numericalSubclass, ""
	}

	// Иначе ожидаем пробел и luminosityClass
	if len(s) >= 4 && s[2] == ' ' {
		luminocityClass := s[3:]
		if isValidLuminosityClass(luminocityClass) {
			return stellarClass, numericalSubclass, luminocityClass
		}
		// даже если невалиден, возвращаем как есть (на случай расширения)
		return stellarClass, numericalSubclass, luminocityClass
	}

	// Неизвестный формат
	return "", "", ""
}

// isValidStellarClass проверяет допустимые спектральные классы (включая D)
func isValidStellarClass(c string) bool {
	switch c {
	case "O", "B", "A", "F", "G", "K", "M", "L", "T", "Y", "D":
		return true
	default:
		return false
	}
}

// isValidLuminosityClass проверяет допустимые классы светимости (включая BD)
func isValidLuminosityClass(lc string) bool {
	switch lc {
	case "Ia", "Ib", "II", "III", "IV", "V", "VI", "D", "BD":
		return true
	default:
		return false
	}
}
