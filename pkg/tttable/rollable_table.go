package tttable

import "io"

type RollableTable interface {
	Roll(Roller, ...string) (string, error)
	SetWriter(io.Writer)
	GetName() string
	FindBykey(string) (string, error)
}
