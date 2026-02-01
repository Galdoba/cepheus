package tttable

import "fmt"

type RollableTable interface {
	Roll(Roller, ...string) (string, string, error)
	GetName() string
	Find(string) (string, error)
	GetAll() map[string]string
}

func AsTable(rt RollableTable) (*Table, error) {
	switch rt := rt.(type) {
	case *Table:
		return rt, nil
	default:
		return nil, fmt.Errorf("not a type 'Table'")
	}
}
