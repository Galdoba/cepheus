package dice

import (
	"fmt"
	"slices"
)

type Manager struct {
	roller    Roller
	rollState *RollState
}

func New(roller Roller) (*Manager, error) {
	return newManager(roller), nil
}

func newManager(r Roller) *Manager {
	if r == nil {
		r = defaultRoller
	}
	return &Manager{
		roller:    r,
		rollState: newRollState(&stdInterpreter{}),
	}
}

func newRollState(i Interpreter) *RollState {
	rs := RollState{}
	rs.Result = &Result{}
	rs.Result.Raw = make(map[*Die]int)
	rs.Interpreter = i
	return &rs
}

func (m *Manager) Result() *Result {
	return m.rollState.Result
}

func (m *Manager) RollSafe(expr string) error {
	expStruct, ok := exprCache.Get(expr)
	if !ok {
		var err error
		expStruct, err = newExpression(expr)
		if err != nil {
			return fmt.Errorf("failed to parse expresion %q: %w", expr, err)
		}
		exprCache.Set(expr, expStruct)
	}
	m.rollState.Expression = expStruct
	m.rollState.Result = basicRoll(m.roller, m.rollState.Expression.Dicepool)
	return nil
}

type RollState struct {
	Expression  *Expression
	Result      *Result
	Interpreter Interpreter
}

type Expression struct {
	Code     string
	Dicepool Dicepool
	Mods     []Mod
}

func newExpression(expr string) (*Expression, error) {
	dp, mods, err := parseExpression(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression %q: %w", expr, err)
	}
	e := Expression{
		Code:     expr,
		Dicepool: dp,
		Mods:     mods,
	}
	return &e, nil
}

type Result struct {
	Raw map[*Die]int
}

func (rs *RollState) Sum() (int, error) {
	switch rs.Interpreter.(type) {
	default:
		return 0, fmt.Errorf("wrong interpreter type")
	case *stdInterpreter:
	}
	in, err := rs.Interpreter.Interpret(rs)
	if err != nil {
		return 0, fmt.Errorf("interpretation failed: %w", err)
	}
	return in.Sum, nil
}

type Interpreter interface {
	Interpret(*RollState) (Interpretation, error)
}

type Interpretation struct {
	Sum   int
	Code  string
	Valid bool
}

type stdInterpreter struct{}

func (si stdInterpreter) Interpret(rs *RollState) (Interpretation, error) {
	raw := []int{}
	for _, value := range rs.Result.Raw {
		raw = append(raw, value)
	}

	mid := slices.Clone(raw)
	var err error
	for _, m := range rs.Expression.Mods {
		mid, err = m.Apply(mid)
		if err != nil {
			return Interpretation{}, fmt.Errorf("failed to apply mod: %w", err)
		}
	}
	s := 0
	for _, v := range mid {
		s += v
	}
	in := Interpretation{}
	in.Sum = s
	in.Code = fmt.Sprintf("%d", s)
	in.Valid = true
	return in, nil
}

type concInterpreter struct{}
