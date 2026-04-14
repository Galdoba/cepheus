package tables

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/engine/dice"
)

var (
	DefaultUpperBound = 1000
	DefaultLowerBound = -1000
)

type GameTable struct {
	Name       string            `json:"name"`
	Expression string            `json:"expression"`
	Data       map[string]string `json:"data"`
	D66        bool              `json:"d_66"`
}

func New(name, expression string, data map[string]string) GameTable {
	return GameTable{Name: name, Expression: expression, Data: data}
}

func (t GameTable) WithD66Roll(d66 bool) GameTable {
	t.D66 = d66
	t.Expression = ""
	return t
}

func (t GameTable) Validate() error {
	if len(t.Name) == 0 {
		return errors.New("table name cannot be empty")
	}
	if len(t.Data) < 2 {
		return fmt.Errorf("table %q must have at least 2 entries", t.Name)
	}
	switch t.D66 {
	case false:
		if err := dice.ValidateExpression(t.Expression); err != nil {
			return fmt.Errorf("table %q expression is not parseable: %w", t.Name, err)
		}
	case true:
		if t.Expression != "" {
			return fmt.Errorf("table %q is D66 table: expression is not empty", t.Name)
		}
	}
	indexes := make([]int, 0, len(t.Data))
	indexMet := make(map[int]int)
	for k := range t.Data {
		idx, err := stringToIndexes(k)
		if err != nil {
			return fmt.Errorf("table %q has invalid index %q: %w", t.Name, k, err)
		}
		for _, i := range idx {
			indexMet[i]++
			if indexMet[i] > 1 {
				return fmt.Errorf("table %q: index duplication: %d", t.Name, i)
			}
		}
		indexes = append(indexes, idx...)
	}
	sort.Ints(indexes)
	min, max := indexes[0], indexes[len(indexes)-1]
	expectedCount := max - min + 1
	if !t.D66 && len(indexes) != expectedCount {
		return fmt.Errorf("table %q has holes in index range [%d, %d]", t.Name, min, max)
	}
	for _, idx := range indexes {
		if idx < -1000 || idx > 1000 {
			return fmt.Errorf("table %q has index %d out of bounds [-1000, 1000]", t.Name, idx)
		}
		if idx == andAbove || idx == andBelow {
			return fmt.Errorf("table %q contains marker index %d", t.Name, idx)
		}
	}
	for _, v := range t.Data {
		if len(v) == 0 {
			return fmt.Errorf("table %q has empty value", t.Name)
		}
	}
	return nil
}

type TableCollection struct {
	tables map[string]GameTable
}

const (
	andAbove = 1001
	andBelow = -1001
)

func indexesToString(indexes ...int) (string, error) {
	if len(indexes) == 0 {
		return "", nil
	}

	hasAbove := false
	hasBelow := false
	for _, idx := range indexes {
		if idx == andAbove {
			hasAbove = true
		}
		if idx == andBelow {
			hasBelow = true
		}
	}

	if hasAbove || hasBelow {
		if len(indexes) != 2 {
			return "", errors.New("andAbove/andBelow requires exactly 2 arguments")
		}
		if hasAbove {
			return fmt.Sprintf("%d+", indexes[0]), nil
		}
		return fmt.Sprintf("%d-", indexes[0]), nil
	}

	sorted := make([]int, len(indexes))
	copy(sorted, indexes)
	sort.Ints(sorted)

	var groups []string
	start := sorted[0]
	end := sorted[0]

	for i := 1; i < len(sorted); i++ {
		if sorted[i] == end+1 {
			end = sorted[i]
		} else {
			if start == end {
				groups = append(groups, fmt.Sprintf("%d", start))
			} else {
				groups = append(groups, fmt.Sprintf("%d - %d", start, end))
			}
			start = sorted[i]
			end = sorted[i]
		}
	}

	if start == end {
		groups = append(groups, fmt.Sprintf("%d", start))
	} else {
		groups = append(groups, fmt.Sprintf("%d - %d", start, end))
	}

	return strings.Join(groups, ", "), nil
}

func stringToIndexes(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty string is not valid input")
	}

	parts := strings.Split(s, ",")
	result := make(map[int]bool)

	rangePattern := regexp.MustCompile(`^(-?\d+)\s*-\s*(-?\d+)$`)
	plusPattern := regexp.MustCompile(`^(-?\d+)\+$`)
	minusPattern := regexp.MustCompile(`^(-?\d+)\-$`)
	numberPattern := regexp.MustCompile(`^(-?\d+)$`)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if match := rangePattern.FindStringSubmatch(part); match != nil {
			min := mustParseInt(match[1])
			max := mustParseInt(match[2])
			if min > max {
				return nil, fmt.Errorf("invalid range: %d > %d", min, max)
			}
			for i := min; i <= max; i++ {
				result[i] = true
			}
			continue
		}

		if match := plusPattern.FindStringSubmatch(part); match != nil {
			n := mustParseInt(match[1])
			for i := n; i <= DefaultUpperBound; i++ {
				result[i] = true
			}
			continue
		}

		if match := minusPattern.FindStringSubmatch(part); match != nil {
			n := mustParseInt(match[1])
			for i := DefaultLowerBound; i <= n; i++ {
				result[i] = true
			}
			continue
		}

		if match := numberPattern.FindStringSubmatch(part); match != nil {
			n := mustParseInt(match[1])
			result[n] = true
			continue
		}

		return nil, fmt.Errorf("invalid token: %q", part)
	}

	if len(result) == 0 {
		return nil, errors.New("no valid indexes found")
	}

	indexes := make([]int, 0, len(result))
	for k := range result {
		indexes = append(indexes, k)
	}
	sort.Ints(indexes)

	return indexes, nil
}

func mustParseInt(s string) int {
	n := 0
	fmt.Sscanf(s, "%d", &n)
	return n
}
