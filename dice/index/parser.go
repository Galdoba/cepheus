package index

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func Parse(str string) ([]int, error) {
	output := []int{}
	parts := strings.Split(str, ",")
	for _, part := range parts {
		//pure tn
		val, err := strconv.Atoi(part)
		if err == nil {
			output = append(output, val)
			continue
		}
		//tn or more
		if strings.HasSuffix(part, "+") {
			partSegment := strings.TrimSuffix(part, "+")
			val, err := strconv.Atoi(partSegment)
			if err == nil {
				for i := val; i <= 30; i++ {
					output = append(output, i)
				}
				continue
			}
		}
		//tn or less
		if strings.HasSuffix(part, "-") {
			partSegment := strings.TrimSuffix(part, "-")
			val, err := strconv.Atoi(partSegment)
			if err == nil {
				for i := val; i >= -30; i-- {
					output = append(output, i)
				}
				continue
			}
		}
		//tn from x to y
		segments := strings.Split(part, "...")
		if len(segments) == 2 {
			x, errX := strconv.Atoi(segments[0])
			y, errY := strconv.Atoi(segments[1])
			if errX == nil && errY == nil {
				for i := x; i <= y; i++ {
					output = append(output, i)
				}
				continue
			}
		}
		return nil, fmt.Errorf("failed to parse goal '%v'\n", str)

	}

	return uniqueInts(output), nil
}

func uniqueInts(sl []int) []int {
	intMap := make(map[int]int)
	for _, i := range sl {
		intMap[i]++
	}
	sl = []int{}
	for k := range intMap {
		sl = append(sl, k)
	}
	slices.Sort(sl)
	return sl
}

func Format(nums []int) (string, error) {
	// Проверка валидности входных чисел
	for _, num := range nums {
		if num < -100 || num > 100 {
			return "", fmt.Errorf("number %d is out of range [-100, 100]", num)
		}
	}

	if len(nums) == 0 {
		return "", nil
	}

	// Удаление дубликатов и сортировка
	set := make(map[int]struct{})
	unique := make([]int, 0, len(nums))
	for _, num := range nums {
		if _, exists := set[num]; !exists {
			set[num] = struct{}{}
			unique = append(unique, num)
		}
	}
	sort.Ints(unique)

	// Группировка последовательных чисел
	var groups [][]int
	currentGroup := []int{unique[0]}

	for i := 1; i < len(unique); i++ {
		if unique[i] == currentGroup[len(currentGroup)-1]+1 {
			currentGroup = append(currentGroup, unique[i])
		} else {
			groups = append(groups, currentGroup)
			currentGroup = []int{unique[i]}
		}
	}
	groups = append(groups, currentGroup)

	// Форматирование групп
	var parts []string
	for _, group := range groups {
		n := len(group)
		start := group[0]
		end := group[n-1]

		switch {
		case n == 1:
			parts = append(parts, strconv.Itoa(start))

		case start == -100 && end < 0:
			absEnd := -end
			parts = append(parts, strconv.Itoa(absEnd)+"-")

		case end == 100 && start >= 0:
			parts = append(parts, strconv.Itoa(start)+"+")

		default:
			parts = append(parts, fmt.Sprintf("%d-%d", start, end))
		}
	}

	return strings.Join(parts, ", "), nil
}
