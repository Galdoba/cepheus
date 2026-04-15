package dice

import (
	"fmt"
	"slices"
)

func newManager(r roller) *Manager {
	if r == nil {
		r = defaultRoller
	}
	return &Manager{
		roller:    r,
		rollState: newRollState(&stdInterpreter{}),
	}
}

func newRollState(i interpreter) *rollState {
	return &rollState{
		result:      result{},
		interpreter: i,
	}
}

// Result returns the last roll result. The caller must not modify the returned slices.
// This method is safe for concurrent use.
func (m *Manager) Result() Result {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Result{dice: m.rollState.result.dice, raw: m.rollState.result.raw}
}

// roll parses the expression (using cache) and performs the basic roll.
// It must be called with the mutex held.
func (m *Manager) roll(expr string) error {
	expStruct, ok := exprCache.get(expr)
	if !ok {
		var err error
		expStruct, err = newExpression(expr)
		if err != nil {
			return fmt.Errorf("failed to parse expression %q: %w", expr, err)
		}
		exprCache.set(expr, expStruct)
	}
	m.rollState.expression = expStruct
	m.rollState.result = basicRoll(m.roller, m.rollState.expression.dicepool)
	return nil
}

type rollState struct {
	expression  *expression
	result      result
	interpreter interpreter
}

type expression struct {
	code     string
	dicepool dicepool
	mods     []mod
}

func newExpression(expr string) (*expression, error) {
	dp, mods, err := parseExpression(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression %q: %w", expr, err)
	}
	return &expression{
		code:     expr,
		dicepool: dp,
		mods:     mods,
	}, nil
}

type result struct {
	dice []die
	raw  []int
}

type interpreter interface {
	interpret(*rollState) (interpretation, error)
}

type interpretation struct {
	sum   int
	code  string
	valid bool
}

type stdInterpreter struct{}

func (si stdInterpreter) interpret(rs *rollState) (interpretation, error) {
	mid := slices.Clone(rs.result.raw)
	var err error
	for _, m := range rs.expression.mods {
		mid, err = m.apply(mid)
		if err != nil {
			return interpretation{}, fmt.Errorf("failed to apply mod: %w", err)
		}
	}
	sum := 0
	for _, v := range mid {
		sum += v
	}
	return interpretation{sum: sum, code: fmt.Sprintf("%d", sum), valid: true}, nil
}
