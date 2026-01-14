package coordinates

import "fmt"

/////////////////////////////////////////////////////////////////////////////////////////

type Cube struct {
	q, r, s int
}

func NewCube(q, r, s int) (Cube, error) {
	if q+r+s != 0 {
		return Cube{}, fmt.Errorf("invalid coordinates provided")
	}
	return Cube{q, r, s}, nil
}

func MustCube(q, r, s int) Cube {
	c, err := NewCube(q, r, s)
	if err != nil {
		panic(fmt.Sprintf("MustCube failed:\ninput: [%v, %v, %v]\nerror: %v", q, r, s, err))
	}
	return c
}

func (c Cube) Coordinates() (int, int, int) {
	return c.q, c.r, c.s
}

func (c Cube) Q() int {
	return c.q
}

func (c Cube) R() int {
	return c.r
}

func (c Cube) S() int {
	return c.s
}

/////////////////////////////////////////////////////////////////////////////////////////

type Global struct {
	x, y int
}

func NewGlobal(x, y int) Global {
	return Global{x, y}
}

func (g Global) Coordinates() (int, int) {
	return g.x, g.y
}

func (g Global) X() int {
	return g.x
}

func (g Global) Y() int {
	return g.y
}

/////////////////////////////////////////////////////////////////////////////////////////

type Local struct {
	sx, sy int
	lx, ly int
}

func NewLocal(sx, sy, lx, ly int) (Local, error) {
	if ly < 1 {
		return Local{}, fmt.Errorf("internal y (%v) is less than 1", ly)
	}
	if ly > sectorHeight {
		return Local{}, fmt.Errorf("internal y (%v) is less than sectorHeight (%v)", ly, sectorHeight)
	}
	if lx < 1 {
		return Local{}, fmt.Errorf("internal x (%v) is less than 1", lx)
	}
	if lx > sectorWidth {
		return Local{}, fmt.Errorf("internal x (%v) is less than sectorWidth (%v)", lx, sectorWidth)
	}
	return Local{sx, sy, lx, ly}, nil
}

func MustLocal(sx, sy, lx, ly int) Local {
	l, err := NewLocal(sx, sy, lx, ly)
	if err != nil {
		panic(fmt.Sprintf("MustLocal failed:\ninput: [%v, %v, %v, %v]\nerror: %v", sx, sy, lx, ly, err))
	}
	return l
}

func (l Local) Coordinates() (int, int, int, int) {
	return l.sx, l.sy, l.lx, l.ly
}

func (l Local) SectorX() int {
	return l.sx
}

func (l Local) SectorY() int {
	return l.sy
}

func (l Local) InternalX() int {
	return l.lx
}

func (l Local) InternalY() int {
	return l.ly
}
