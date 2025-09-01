// hex.go
package cube

// Hex - неизменяемая структура гекса
type Cube struct {
	Q, R, S int // кубические координаты
}

// NewHex создает гекс с указанными координатами
func NewCube(q, r, s int) Cube {
	h := Cube{
		Q: q,
		R: r,
		S: s,
	}
	return h
}
